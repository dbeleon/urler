package urler

import (
	"context"
	"errors"

	"github.com/dbeleon/urler/urler/internal/domain"
	api "github.com/dbeleon/urler/urler/pkg/urler/v1"

	//"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUrl(ctx context.Context, req *api.GetUrlRequest) (*api.GetUrlResponse, error) {
	long, err := s.businessLogic.GetUrl(ctx, req.GetUrl())
	if err != nil {
		if errors.Is(err, domain.ErrShortUrlNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	//header := metadata.Pairs("Location", long)
	//grpc.SendHeader(ctx, header)

	return &api.GetUrlResponse{
		Url: long,
	}, nil
}
