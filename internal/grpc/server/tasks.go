package grpcserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/literaen/simple_project/tasks/internal/config"

	grpchandler "github.com/literaen/simple_project/tasks/internal/grpc/handler"

	grpcserver "github.com/literaen/simple_project/pkg/grpc/server"

	taskpb "github.com/literaen/simple_project/proto/gen"
)

type TaskGRPCServer struct {
	server *grpcserver.Server
}

func NewTaskGRPCServer(cfg *config.Config, taskService *grpchandler.TaskHandler) *TaskGRPCServer {
	srv := grpcserver.NewServer(5 * time.Second)

	go func() {
		taskpb.RegisterTaskServiceServer(srv.GetServer(), taskService)

		err := srv.Start(context.TODO(), fmt.Sprintf(":%s", cfg.GRPC_Port))
		if err != nil {
			log.Fatalf("error while starting grpc task server: %v", err)
		}
	}()

	return &TaskGRPCServer{server: srv}
}
