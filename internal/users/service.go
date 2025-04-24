package users

import (
	"context"

	grpcclient "github.com/literaen/simple_project/tasks/internal/grpc/client"
)

type UserService struct {
	grpc *grpcclient.UserGRPCClient
}

func NewUserService(grpc *grpcclient.UserGRPCClient) *UserService {
	svc := &UserService{
		grpc: grpc,
	}

	return svc
}
func (s *UserService) GetUser(ctx context.Context, id uint64) error {
	if err := s.grpc.GetUser(ctx, id); err != nil {
		return err
	}

	return nil
}
