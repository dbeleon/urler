package urler

import (
	"context"

	api "github.com/dbeleon/urler/urler/pkg/urler/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetShorts(ctx context.Context, req *api.GetShortsRequest) (*api.GetShortsResponse, error) {
	shorts, err := s.businessLogic.GetShorts(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetShortsResponse{
		Shorts: shorts,
	}, nil
}
