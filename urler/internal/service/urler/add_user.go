package urler

import (
	"context"

	api "github.com/dbeleon/urler/urler/pkg/urler/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddUser(ctx context.Context, req *api.AddUserRequest) (*api.AddUserResponse, error) {
	user, err := s.businessLogic.AddUser(ctx, req.GetName(), req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &api.AddUserResponse{User: user}
	return resp, nil
	//return &emptypb.Empty{}, nil
	//return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
