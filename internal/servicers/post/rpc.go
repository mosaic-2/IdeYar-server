package postImpl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

const (
	landingPostsCount = 5
)

type Post struct {
	ID              int64
	UserID          int64
	Username        string
	ProfileImageUrl string
	Title           string
	Description     string
	Image           string
	MinimumFund     decimal.Decimal
	FundRaised      decimal.Decimal
	DeadlineDate    time.Time
	CreatedAt       time.Time
}

func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	db := dbutil.GormDB(ctx)

	postID := req.GetId()

	post := Post{}
	post.ID = postID

	err := db.Table("post AS p").
		Joins("JOIN user_t AS u ON p.user_id = u.id").
		Where("p.id = ?", postID).
		Select("p.*, u.username, u.profile_image_url").
		Scan(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "could not find post")
		}
		return nil, status.Error(codes.Internal, "error retreiving post")
	}

	postPb := &pb.Post{
		Id:              postID,
		Title:           post.Title,
		Image:           post.Image,
		MinimumFund:     post.MinimumFund.String(),
		FundRaised:      post.FundRaised.String(),
		DeadlineDate:    post.DeadlineDate.Format(time.DateOnly),
		CreatedAt:       timestamppb.New(post.CreatedAt),
		UserId:          post.UserID,
		Username:        post.Username,
		ProfileImageUrl: post.ProfileImageUrl,
	}

	postDetails := []*model.PostDetail{}

	err = db.Model(model.PostDetail{}).
		Where("post_id = ?", postID).
		Scan(&postDetails).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Internal, "error retreiving post")
		}
	}

	postDetailsPb := []*pb.PostDetail{}

	for _, postDetail := range postDetails {
		postDetailsPb = append(postDetailsPb, &pb.PostDetail{
			Title:       postDetail.Title,
			Description: postDetail.Description,
			Image:       &postDetail.Image,
			Order:       postDetail.Order,
		})
	}

	return &pb.GetPostResponse{
		Post:        postPb,
		PostDetails: postDetailsPb,
	}, nil

}

func (s *Server) SearchPost(ctx context.Context, req *pb.SearchPostRequest) (*pb.SearchPostResponse, error) {
	db := dbutil.GormDB(ctx)
	title := req.GetTitle()
	offset := int(req.GetPage()) * 20

	var result []*pb.PostOverview
	filter := req.GetFilter()

	query := db.Table("post AS p").
		Select("p.id, p.title, pd.image").
		Joins("LEFT JOIN (SELECT pd.post_id, pd.image FROM post_detail pd WHERE pd.order_c = 0) AS pd ON p.id = pd.post_id").
		Where("SIMILARITY(p.title, ?) > 0", title)

	if filter != nil {
		categories := filter.GetCategories()
		if len(categories) > 0 {
			query = query.Where("p.category IN (?)", categories)
		}

		// FIXME: This code has possible sql injection
		switch filter.GetSortBy() {
		case pb.SearchPostRequest_Filters_CREATED_TIME:
			query = query.Order("p.created_at DESC")
		case pb.SearchPostRequest_Filters_DEADLINE:
			query = query.Order("p.deadline_date ASC")
		default:
			query = query.Order(fmt.Sprintf("SIMILARITY(p.title, '%s') DESC", title))
		}
	} else {
		query = query.Order(fmt.Sprintf("SIMILARITY(p.title, '%s') DESC", title))
	}

	query = query.Limit(20).Offset(offset)

	if err := query.Scan(&result).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving posts: %v", err)
	}

	return &pb.SearchPostResponse{PostOverview: result}, nil
}

func (s *Server) LandingPosts(ctx context.Context, in *emptypb.Empty) (*pb.LandingPostsResponse, error) {
	db := dbutil.GormDB(ctx)

	posts := []*Post{}

	err := db.Raw(`
		SELECT p.*, u.username, u.profile_image_url
		FROM post p
		JOIN user_t u ON p.user_id = u.id
		ORDER BY RANDOM()
		LIMIT ?
	`, landingPostsCount).Scan(&posts).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "error while retreiving posts")
	}

	postsPb := convertPostToPostPb(posts)

	return &pb.LandingPostsResponse{
		Posts: postsPb,
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

	userFunds := []struct {
		Post
		Amount decimal.Decimal
	}{}

	err := db.Table("fund AS f").
		Joins("JOIN post p ON f.post_id = p.id").
		Joins("JOIN user_t u ON p.user_id = u.id").
		Where("f.user_id = ?", userID).
		Select("p.*, f.amount, u.username, u.profile_image_url").
		Scan(&userFunds).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	userFundsPb := []*pb.FundOverview{}

	for _, fund := range userFunds {
		userFundsPb = append(userFundsPb, &pb.FundOverview{
			Post: &pb.Post{
				Id:              fund.ID,
				UserId:          fund.UserID,
				Username:        fund.Username,
				ProfileImageUrl: fund.ProfileImageUrl,
				Title:           fund.Title,
				Description:     fund.Description,
				MinimumFund:     fund.MinimumFund.String(),
				FundRaised:      fund.FundRaised.String(),
				DeadlineDate:    fund.DeadlineDate.Format(time.DateOnly),
				Image:           fund.Image,
				CreatedAt:       timestamppb.New(fund.DeadlineDate),
			},
			Amount: fund.Amount.String(),
		})
	}

	return &pb.UserFundsResponse{FundOverviews: userFundsPb}, nil
}

func (s *Server) UserProjects(ctx context.Context, req *emptypb.Empty) (*pb.UserProjectsResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	userProjects, err := fetchUserIDProjects(userID, db)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserProjectsResponse{Posts: userProjects}, nil
}

func (s *Server) UserIDProjects(ctx context.Context, req *pb.UserIDProjectsRequest) (*pb.UserProjectsResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := req.GetId()

	userProjects, err := fetchUserIDProjects(userID, db)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserProjectsResponse{Posts: userProjects}, nil
}

func (s *Server) BookmarkPost(ctx context.Context, req *pb.BookmarkPostRequest) (*emptypb.Empty, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)
	postID := req.GetPostId()

	bookmark := model.Bookmark{
		UserID: userID,
		PostID: postID,
	}

	var existingBookmark model.Bookmark
	err := db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingBookmark).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.Internal, "error accessing the database")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = db.Create(&bookmark).Error
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error while bookmarking")
		}
	} else {
		err = db.Delete(&existingBookmark).Error
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error while removing bookmark")
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UserBookmarks(ctx context.Context, _ *emptypb.Empty) (*pb.UserBookmarksResponse, error) {
	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	userBookmarks := []*Post{}

	err := db.Table("post AS p").
		Joins("JOIN bookmark AS b ON p.id = b.post_id").
		Joins("JOIN user_t AS u ON u.id = p.user_id").
		Where("b.user_id = ?", userID).
		Select("p.*, u.id AS user_id, u.username, u.profile_image_url").
		Scan(&userBookmarks).Error
	if err != nil {
		return nil, err
	}

	posts := convertPostToPostPb(userBookmarks)
	return &pb.UserBookmarksResponse{Posts: posts}, nil
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

func fetchUserIDProjects(userID int64, tx *gorm.DB) ([]*pb.Post, error) {

	userProjects := []*Post{}

	err := tx.Table("post AS p").
		Joins("JOIN user_t AS u ON p.user_id = u.id").
		Where("p.user_id = ?", userID).
		Select("p.*, u.id AS user_id, u.username, u.profile_image_url").
		Scan(&userProjects).Error
	if err != nil {
		return nil, err
	}

	return convertPostToPostPb(userProjects), nil
}

func convertPostToPostPb(posts []*Post) []*pb.Post {

	result := []*pb.Post{}

	for _, post := range posts {
		result = append(result, &pb.Post{
			Id:              post.ID,
			UserId:          post.ID,
			Username:        post.Username,
			ProfileImageUrl: post.ProfileImageUrl,
			Title:           post.Title,
			Description:     post.Description,
			MinimumFund:     post.MinimumFund.String(),
			FundRaised:      post.FundRaised.String(),
			DeadlineDate:    post.DeadlineDate.Local().Format(time.DateOnly),
			Image:           post.Image,
			CreatedAt:       timestamppb.New(post.CreatedAt),
		})
	}

	return result
}
