package postImpl

import (
	"context"
	"errors"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	pb "github.com/mosaic-2/IdeYar-server/pkg/postServicePb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {

	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	var id int64

	err := db.Transaction(func(tx *gorm.DB) error {
		post, postDetails, err := toPostCreatePayload(req, userID)
		if err != nil {
			return err
		}

		err = validatePost(post)
		if err != nil {
			return err
		}

		err = tx.Create(&post).Error
		if err != nil {
			return err
		}

		id = post.ID

		for _, detail := range postDetails {
			detail.PostID = id
		}

		err = tx.Create(&postDetails).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: id}, nil
}

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
		UserId:      post.UserID,
		Title:       post.Title,
		MinimumFund: post.MinimumFund.String(),
		FundRaised:  post.FundRaised.String(),
		PostDetails: postDetailsPb,
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
		return nil, status.Error(codes.Internal, "internal server error")
	}

	fund := model.Fund{
		UserID: userID,
		PostID: postID,
		Amount: amount,
	}

	err = db.Create(fund).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &emptypb.Empty{}, nil
}

func UserFunds(ctx context.Context, req *emptypb.Empty) (*pb.UserFundsResponse, error) {

	db := dbutil.GormDB(ctx)

	userID := util.GetUserIDFromCtx(ctx)

	userFunds := []*pb.PostOverview{}

	err := db.Table("fund AS f").
		Joins("JOIN post p ON f.post_id = p.id").
		Joins("JOIN post_detail pd ON p.id = pd.post_id").
		Where("f.user_id = ? AND pd.order = 0", userID).
		Select("p.id, p.title, pd.image").
		Scan(&userFunds).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserFundsResponse{PostOverview: userFunds}, nil
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

func toPostCreatePayload(req *pb.CreateRequest, userID int64) (*model.Post, []*model.PostDetail, error) {

	minimumFund, err := decimal.NewFromString(req.GetMinimumFund())
	if err != nil {
		return nil, nil, status.Errorf(codes.InvalidArgument, "invalid minimum fund format")
	}

	post := &model.Post{
		Title:       req.GetTitle(),
		MinimumFund: minimumFund,
		UserID:      userID,
	}

	postDetails := []*model.PostDetail{}

	for _, detail := range req.GetPostDetails() {
		postDetail := &model.PostDetail{
			Order:       detail.GetOrder(),
			Title:       detail.GetTitle(),
			Description: detail.GetDescription(),
		}
		postDetails = append(postDetails, postDetail)
	}

	return post, postDetails, nil
}

func validatePost(post *model.Post) error {
	if post.Title == "" {
		return status.Errorf(codes.InvalidArgument, "post title can not be empty")
	}

	return nil
}

func fetchUserIDProjects(userID int64, tx *gorm.DB) ([]*pb.PostOverview, error) {

	userProjects := []*pb.PostOverview{}

	err := tx.Table("post AS p").
		Joins("JOIN post_detail pd ON p.id = pd.post_id").
		Where("p.user_id = ? AND pd.order = 0", userID).
		Select("p.id, p.title, pd.image").
		Scan(&userProjects).Error
	if err != nil {
		return nil, err
	}

	return userProjects, nil
}
