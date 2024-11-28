package authservice

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	pb "github.com/mosaic-2/IdeYar-server/pkg/authservicepb"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	// convert email address to lower case
	req.Email = strings.ToLower(req.GetEmail())

	db := dbutil.GormDB(ctx)

	err := checkSignUpPreconditions(req, db)
	if err != nil {
		return nil, err
	}

	signUpExpTime := time.Now().Add(5 * time.Minute)
	verificationCode := util.GenerateVerificationCode()
	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password")
	}

	signUp := model.SignUp{
		Email:            req.GetEmail(),
		Username:         req.GetUsername(),
		Password:         string(bcryptPass),
		VerificationCode: verificationCode,
		Expire:           signUpExpTime,
	}

	err = db.Create(&signUp).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	go util.SendSignUpEmail(req.GetEmail(), verificationCode)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(signUpExpTime),
		Issuer:    "KhanWeb",
		Subject:   strconv.Itoa(int(signUp.ID)),
		Audience:  jwt.ClaimStrings{"SignUp"},
	})

	tokenString, err := token.SignedString(s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token")
	}

	return &pb.SignUpResponse{Token: tokenString}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	db := dbutil.GormDB(ctx)

	user, err := getUser(req, db)
	if err != nil {
		return nil, err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, status.Errorf(codes.InvalidArgument, "Incorrect password")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking password")
	}

	// generate Token
	tokenString, err := util.CreateLoginToken(strconv.FormatInt(user.ID, 10), time.Hour*12, s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token")
	}
	refreshTokenString, err := util.CreateRefreshToken(strconv.FormatInt(user.ID, 10), time.Hour*100, s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating refresh token")
	}

	return &pb.LoginResponse{
		JwtToken:     tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *Server) CodeVerification(ctx context.Context, req *pb.CodeVerificationRequest) (*pb.CodeVerificationResponse, error) {
	db := dbutil.GormDB(ctx)

	token, err := jwt.Parse(req.SignUpToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method: %v", token.Header["alg"])
		}
		return s.hmacSecret, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	aud, _ := token.Claims.GetAudience()
	if len(aud) != 1 || aud[0] != "SignUp" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	signUpIDStr, _ := token.Claims.GetSubject()
	signUpID, err := strconv.Atoi(signUpIDStr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	signUpData := &model.SignUp{
		ID: int32(signUpID),
	}

	err = db.Take(signUpData).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if signUpData.VerificationCode != req.Code {
		return nil, status.Errorf(codes.InvalidArgument, "Wrong Code")
	}

	err = db.Delete(signUpData).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	tx := db.Begin(&sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	defer func(tx *gorm.DB) {
		_ = tx.Commit()
	}(tx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var duplicateUsername bool
	err = tx.Raw(`
		SELECT COUNT(*) > 0
		FROM user_t
		WHERE username = ?
	`, signUpData.Username).Scan(&duplicateUsername).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if duplicateUsername {
		return nil, status.Errorf(codes.InvalidArgument, "username already exists")
	}

	var duplicateEmail bool
	err = tx.Raw(`
		SELECT COUNT(*) > 0
		FROM user_t
		WHERE email = ?
	`, signUpData.Email).Scan(&duplicateEmail).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user := &model.User{
		Email:    signUpData.Email,
		Username: signUpData.Username,
		Password: signUpData.Password,
	}

	err = tx.Create(user).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	loginToken, err := util.CreateLoginToken(strconv.Itoa(int(user.ID)), time.Hour*12, s.hmacSecret)
	if err != nil {
		_ = tx.Rollback()
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	refreshTokenString, err := util.CreateRefreshToken(strconv.FormatInt(user.ID, 10), time.Hour*100, s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating refresh token")
	}

	return &pb.CodeVerificationResponse{
		JwtToken:     loginToken,
		RefreshToken: refreshTokenString,
	}, nil
}

func checkSignUpPreconditions(req *pb.SignUpRequest, db *gorm.DB) error {
	if !util.ValidateEmail(req.GetEmail()) {
		return status.Errorf(codes.InvalidArgument, "invalid email")
	}
	if !util.ValidateUsername(req.GetUsername()) {
		return status.Errorf(codes.InvalidArgument, "invalid username")
	}
	if !util.ValidatePassword(req.GetPassword()) {
		return status.Errorf(codes.InvalidArgument, "invalid password")
	}

	var duplicateEmail bool
	err := db.Raw(`
		SELECT count(*) > 0
		FROM user_t
		WHERE email = ?
	`, req.GetEmail()).Scan(&duplicateEmail).Error
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	if duplicateEmail {
		return status.Errorf(codes.AlreadyExists, "email already exist")
	}

	var duplicateUsername bool
	err = db.Raw(`
		SELECT COUNT(*) > 0
		FROM user_t
		WHERE username = ?
	`, req.GetUsername()).Scan(&duplicateUsername).Error
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	if duplicateUsername {
		return status.Errorf(codes.AlreadyExists, "username already exist")
	}

	return nil
}

func getUser(req *pb.LoginRequest, db *gorm.DB) (*model.User, error) {

	var user *model.User

	err := db.Where("email = ? OR username = ?", req.GetUserNameOrEmail(), req.GetUserNameOrEmail()).
		Take(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	}

	return user, nil
}
