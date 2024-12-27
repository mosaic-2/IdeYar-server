package userProfileImpl

import (
	"context"
	"errors"
	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	pb "github.com/mosaic-2/IdeYar-server/pkg/UserProfileServicePb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

func (s *Server) ChangeEmail(ctx context.Context, in *pb.ChangeEmailRequest) (*pb.ChangeEmailResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	email := in.GetEmail()
	if !util.ValidateEmail(email) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format: %s", email)
	}

	// Directly update the user's email in the database
	if err := db.Model(&model.User{}).Where("id = ?", userID).Update("email", email).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update email: %v", err)
	}

	return &pb.ChangeEmailResponse{}, nil
}

func (s *Server) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	newPassword := in.GetNewPassword()
	if !util.ValidatePassword(newPassword) {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters long")
	}
	bcryptNewPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password")
	}

	if err := db.Model(&model.User{}).Where(
		"id = ?", userID,
	).Update("password", string(bcryptNewPass)).Error; err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update password: %v", err)
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (s *Server) GetProfileInfo(ctx context.Context, in *pb.GetProfileInfoRequest) (*pb.GetProfileInfoResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	var user model.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user profile: %v", err)
	}

	response := pb.GetProfileInfoResponse{
		Username:        user.Username,
		Phone:           user.Phone,
		Bio:             user.Bio,
		Birthday:        user.Birthday.Format("2006-01-02"),
		ProfileImageUrl: user.ProfileImageURL,
		Email:           user.Email,
	}

	return &response, nil
}

func (s *Server) UpdateProfileInfo(ctx context.Context, in *pb.UpdateProfileInfoRequest) (*pb.UpdateProfileInfoResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	updates := make(map[string]interface{})

	if in.Username != "" {
		updates["username"] = in.Username
	}
	if in.Phone != "" {
		updates["phone"] = in.Phone
	}
	if in.Bio != "" {
		updates["bio"] = in.Bio
	}
	if in.Birthday != "" {
		birthday, err := time.Parse("2006-01-02", in.Birthday)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date format for birthday: %v", err)
		}
		updates["birthday"] = birthday
	}

	if err := db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	var user model.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve updated profile: %v", err)
	}

	response := &pb.UpdateProfileInfoResponse{
		Username:        user.Username,
		Phone:           user.Phone,
		Bio:             user.Bio,
		Birthday:        user.Birthday.Format("2006-01-02"),
		ProfileImageUrl: user.ProfileImageURL,
	}

	return response, nil
}
