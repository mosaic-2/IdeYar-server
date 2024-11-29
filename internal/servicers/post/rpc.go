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
