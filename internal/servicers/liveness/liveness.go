package livenessImp

import (
	"context"
	"github.com/mosaic-2/IdeYar-server/pkg/LivenessService"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	LivenessService.UnimplementedLivenessServer
}

func NewServer() (LivenessService.LivenessServer, error) {
	return server{}, nil
}

func (s server) CheckLiveness(ctx context.Context, _ *emptypb.Empty) (*LivenessService.CheckLivenessResponse, error) {
	return &LivenessService.CheckLivenessResponse{
		IsAlive:   true,
		Message:   "IdeYar is alive",
		Timestamp: timestamppb.Now(),
	}, nil
}
