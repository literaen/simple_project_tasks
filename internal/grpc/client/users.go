package grpcclients

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/literaen/simple_project/tasks/internal/config"

	grpcclient "github.com/literaen/simple_project/pkg/grpc/client"

	userpb "github.com/literaen/simple_project/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type UserClientConstructor struct{}

func (c *UserClientConstructor) NewClient(conn *grpc.ClientConn) interface{} {
	return userpb.NewUserServiceClient(conn)
}

type UserGRPCClient struct {
	client *grpcclient.Client
}

func NewUserGRPCClient(cfg *config.Config) *UserGRPCClient {
	client := &UserGRPCClient{
		client: grpcclient.NewClient(5*time.Second, &UserClientConstructor{}),
	}

	client.Start(cfg)

	return client
}

func (s *UserGRPCClient) GetUser(ctx context.Context, id uint64) error {
	if !s.client.IsReady() {
		return fmt.Errorf("user service unavailable")
	}

	_, err := s.GetUserClient().GetUser(ctx, &userpb.GetUserRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return errors.New(st.Message())
		} else {
			return fmt.Errorf("unknown error: %v", err)
		}
	}

	return nil
}

// GetUserClient возвращает типизированный клиент
func (s *UserGRPCClient) GetUserClient() userpb.UserServiceClient {
	return s.client.GetClient().(userpb.UserServiceClient)
}

func (s *UserGRPCClient) Start(cfg *config.Config) {
	target := fmt.Sprintf("%s:%s", cfg.USER_SERVICE_HOST, cfg.USER_SERVICE_PORT)
	s.client.AutoReconnect(context.TODO(), target)
}

func (s *UserGRPCClient) Close() error {
	return s.client.Close()
}
