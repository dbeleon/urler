package urler

import (
	"context"

	api "github.com/dbeleon/urler/urler/pkg/urler/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) MakeUrl(ctx context.Context, req *api.MakeUrlRequest) (*api.MakeUrlResponse, error) {
	short, err := s.businessLogic.MakeUrl(ctx, req.GetUser(), req.GetUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.MakeUrlResponse{
		Url: short,
	}, nil
}
