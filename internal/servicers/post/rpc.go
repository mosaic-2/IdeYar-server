package postImpl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	pb "github.com/mosaic-2/IdeYar-server/pkg/postServicePb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	db := dbutil.GormDB(ctx)

	postID := req.GetId()

	post := model.Post{}
	post.ID = postID

	err := db.Model(&post).Take(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "could not find post")
		}
		return nil, status.Error(codes.Internal, "error retreiving post")
	}

	user := model.User{ID: post.UserID}
	err = db.Model(&user).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "could not find user")
		}
		return nil, status.Error(codes.Internal, "error retreiving user")
	}

	postDetails := []model.PostDetail{}

	err = db.Model(model.PostDetail{}).
		Where("post_id = ?", postID).
		Scan(&postDetails).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Internal, "error retreiving post")
		}
	}

	postDetailsPb := []*pb.PostDetail{}

	for _, postD := range postDetails {
		postDetailsPb = append(postDetailsPb, &pb.PostDetail{
			Title:       postD.Title,
			Description: postD.Description,
			Order:       postD.Order,
			Image:       &postD.Image,
		})
	}

	return &pb.GetPostResponse{
		Id:               post.ID,
		UserId:           post.UserID,
		Username:         user.Username,
		UserProfileImage: user.ProfileImageURL,
		Title:            post.Title,
		MinimumFund:      post.MinimumFund.String(),
		FundRaised:       post.FundRaised.String(),
		DeadlineDate:     post.DeadlineDate.Format(time.DateOnly),
		Image:            post.Image,
		CreatedAt:        timestamppb.New(post.CreatedAt),
		PostDetails:      postDetailsPb,
	}, nil

}

func (s *Server) SearchPost(ctx context.Context, req *pb.SearchPostRequest) (*pb.SearchPostResponse, error) {

	db := dbutil.GormDB(ctx)

	title := req.GetTitle()
	offset := req.GetPage() * 20

	result := []*pb.PostOverview{}

	err := db.Raw(`
		SELECT p.id, p.title, pd.image
		FROM post p LEFT JOIN (
			SELECT pd.post_id, pd.image 
			FROM post_detail pd
			WHERE pd.order_c = 0
		) AS pd ON p.id = pd.post_id
		ORDER BY SIMILARITY(p.title, ?) DESC
		LIMIT 20 
		OFFSET ?
	`, title, offset).
		Scan(&result).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "error while retreiving posts")
	}

	return &pb.SearchPostResponse{
		PostOverview: result,
	}, nil
}

func (s *Server) LandingPosts(ctx context.Context, in *emptypb.Empty) (*pb.LandingPostsResponse, error) {
	db := dbutil.GormDB(ctx)

	result := []*pb.LandingPost{}

	err := db.Raw(`
		SELECT p.id, p.title, pd.image, p.fund_raised, p.minimum_fund
		FROM post p LEFT JOIN post_detail pd ON p.id = pd.post_id
		WHERE order_c = 0
		ORDER BY RANDOM()
	`).Scan(&result).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "error while retreiving posts")
	}

	return &pb.LandingPostsResponse{
		LandingPosts: result,
	}, nil
}

func (s *Server) FundPost(ctx context.Context, req *pb.FundPostRequest) (*emptypb.Empty, error) {

	db := dbutil.GormDB(ctx)

	postID := req.GetPostId()

	userID := util.GetUserIDFromCtx(ctx)

	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount format")
	}

	fund := model.Fund{
		UserID: userID,
		PostID: postID,
		Amount: amount,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&fund).Error
		if err != nil {
			return status.Error(codes.Internal, "internal server error")
		}

		err = db.Exec(`
			UPDATE post
			SET fund_raised = fund_raised + ?
			WHERE id = ?;
		`, amount, postID).Error
		if err != nil {
			return status.Error(codes.Internal, "internal server error")
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UserFunds(ctx context.Context, req *emptypb.Empty) (*pb.UserFundsResponse, error) {

	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	userFunds := []*pb.FundOverview{}

	err := db.Table("fund AS f").
		Joins("JOIN post p ON f.post_id = p.id").
		Joins("JOIN post_detail pd ON p.id = pd.post_id").
		Where("f.user_id = ? AND pd.order_c = 0", userID).
		Select("p.id, p.title, pd.image, f.amount").
		Scan(&userFunds).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserFundsResponse{FundOverview: userFunds}, nil
}

func (s *Server) UserProjects(ctx context.Context, req *emptypb.Empty) (*pb.UserProjectsResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	userProjects, err := fetchUserIDProjects(userID, db)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserProjectsResponse{PostOverview: userProjects}, nil
}

func (s *Server) UserIDProjects(ctx context.Context, req *pb.UserIDProjectsRequest) (*pb.UserProjectsResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := req.GetId()

	userProjects, err := fetchUserIDProjects(userID, db)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserProjectsResponse{PostOverview: userProjects}, nil
}

func validatePost(post model.Post) error {
	if post.Title == "" {
		return status.Errorf(codes.InvalidArgument, "post title can not be empty")
	}

	return nil
}

func validatePostDetail(postDetail model.PostDetail) error {
	if postDetail.Order < 0 || postDetail.Order > 9 {
		return status.Errorf(codes.InvalidArgument, "each post can have at most 10 parts")
	}

	if postDetail.PostID == 0 {
		return status.Errorf(codes.InvalidArgument, "invalid post id")
	}

	return nil
}

func hasCreateAccessPostDetail(tx *gorm.DB, postDetail model.PostDetail, userID int64) (bool, error) {
	var hasAccess bool

	err := tx.Table("post").
		Where("id = ? AND user_id = ?", postDetail.PostID, userID).
		Select("count(*) > 0").
		Scan(&hasAccess).Error
	if err != nil {
		return false, err
	}

	return hasAccess, err
}

func fetchUserIDProjects(userID int64, tx *gorm.DB) ([]*pb.PostOverview, error) {

	userProjects := []*pb.PostOverview{}

	err := tx.Table("post AS p").
		Joins("JOIN user_t AS u ON p.user_id = u.id").
		Where("p.user_id = ?", userID).
		Select("p.id, p.title, p.description, p.image, u.id AS user_id, u.username, u.profile_image_url AS user_profile_image").
		Scan(&userProjects).Error
	if err != nil {
		return nil, err
	}

	return userProjects, nil
}
