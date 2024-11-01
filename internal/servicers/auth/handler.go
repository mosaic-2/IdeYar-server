package auth

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	"github.com/mosaic-2/IdeYar-server/internal/sql/dbpkg"
	pb "github.com/mosaic-2/IdeYar-server/pkg/authpb"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// convert email to lower case
	in.Email = strings.ToLower(in.Email)

	if !util.ValidateEmail(in.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}
	if !util.ValidateUsername(in.Username) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid username")
	}
	if !util.ValidatePassword(in.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password")
	}

	EmailCnt, err := s.query.ExistsUserEmail(ctx, in.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if EmailCnt != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "email already exist")
	}

	UserNameCnt, err := s.query.ExistsUserUsername(ctx, in.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if UserNameCnt != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "username already exist")
	}

	signUpExpTime := time.Now().Add(5 * time.Minute)
	verificationCode := util.GenerateVerificationCode()
	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password")
	}

	signupID, err := s.query.InsertSignup(ctx, dbpkg.InsertSignupParams{
		Email:            in.Email,
		Username:         in.Username,
		Password:         string(bcryptPass),
		VerificationCode: verificationCode,
		Expire:           signUpExpTime,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// send signup email
	// go util.SendSignUpEmail(in.Email, verificationCode)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(signUpExpTime),
		Issuer:    "KhanWeb",
		Subject:   strconv.Itoa(int(signupID)),
		Audience:  jwt.ClaimStrings{"SignUp"},
	})

	tokenString, err := token.SignedString(s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token")
	}

	return &pb.SignUpResponse{Token: tokenString}, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

	// verify user
	// try to get user by email
	userEmail, err1 := s.query.GetUserByEmail(ctx, strings.ToLower(in.UserNameOrEmail))
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "Error retrieving user %s\n", in.UserNameOrEmail)
	}
	// try to get user by username
	userUsername, err2 := s.query.GetUserByUsername(ctx, in.UserNameOrEmail)
	if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "Error retrieving user %s\n", in.UserNameOrEmail)
	}
	// user doesn't exist
	if err1 != nil && err2 != nil {
		return nil, status.Errorf(codes.NotFound, "No such user %s\n", in.UserNameOrEmail)
	}

	var user dbpkg.Account

	if err1 == nil {
		user = userEmail
	} else {
		user = userUsername
	}
	// if not verified raise error
	// compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, status.Errorf(codes.InvalidArgument, "Incorrect password")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking password")
	}

	profileID, err := s.query.GetProfileID(ctx, user.ID)

	// generate Token
	tokenString, err := util.CreateLoginToken(strconv.FormatInt(profileID, 10), time.Hour*12, s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token")
	}
	refreshTokenString, err := util.CreateRefreshToken(strconv.FormatInt(profileID, 10), time.Hour*100, s.hmacSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating refresh token")
	}

	return &pb.LoginResponse{
		JwtToken:     tokenString,
		RefreshToken: refreshTokenString,
	}, nil

}
