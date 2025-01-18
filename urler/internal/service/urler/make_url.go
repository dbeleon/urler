package urler

import (
	"context"
	"errors"

	"github.com/dbeleon/urler/urler/internal/domain"
	api "github.com/dbeleon/urler/urler/pkg/urler/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) MakeUrl(ctx context.Context, req *api.MakeUrlRequest) (*api.MakeUrlResponse, error) {
	short, err := s.businessLogic.MakeUrl(ctx, req.GetUser(), req.GetUrl())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, domain.ErrToManyHashCollisions) {
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.MakeUrlResponse{
		Url: short,
	}, nil
}
