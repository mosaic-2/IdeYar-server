package postservice

import (
	"context"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	pb "github.com/mosaic-2/IdeYar-server/pkg/postservicepb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {

	db := dbutil.GormDB(ctx)

	userID, _ := ctx.Value(util.ProfileIDCtx{}).(int64)

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
