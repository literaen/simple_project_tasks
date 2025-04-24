//go:build wireinject
// +build wireinject

package app

import (
	"context"

	"github.com/literaen/simple_project/tasks/internal/config"
	"github.com/literaen/simple_project/tasks/internal/kafka/consumer"

	//github.com/literaen/simple_project/tasks/internal/database"
	grpcclient "github.com/literaen/simple_project/tasks/internal/grpc/client"
	grpchandler "github.com/literaen/simple_project/tasks/internal/grpc/handler"
	grpcserver "github.com/literaen/simple_project/tasks/internal/grpc/server"
	oapihandler "github.com/literaen/simple_project/tasks/internal/oapi/handler"

	"github.com/literaen/simple_project/tasks/internal/tasks"
	"github.com/literaen/simple_project/tasks/internal/users"

	"github.com/literaen/simple_project/pkg/postgres"
	"github.com/literaen/simple_project/pkg/redis"

	"github.com/google/wire"
)

type App struct {
	Config          *config.Config
	TaskOAPIHandler *oapihandler.TaskHandler
	TaskGRPCHandler *grpchandler.TaskHandler
	TaskGRPCServer  *grpcserver.TaskGRPCServer
}

func InitApp() (*App, error) {
	wire.Build(
		config.LoadEnv,

		config.ProvideDBCreds,
		postgres.NewGDB,

		config.ProvideRedisCreds,
		redis.NewRDB,

		grpcclient.NewUserGRPCClient,
		users.NewUserService,

		oapihandler.NewTaskHandler,

		tasks.NewTaskRepository,
		tasks.NewTaskService,

		grpcserver.NewTaskGRPCServer,
		grpchandler.NewTaskHandler,

		newApp,
	)
	return nil, nil
}

func newApp(
	config *config.Config,
	gdb *postgres.GDB,
	oapiTaskHandler *oapihandler.TaskHandler,
	taskService *tasks.TaskService,
	taskGRPCServer *grpcserver.TaskGRPCServer,
	grpcTaskHandler *grpchandler.TaskHandler,
) *App {
	tasks.Migrate(gdb.DB)

	consumer := consumer.NewUserEventConsumer(
		taskService,
		config.KAFKA_BROKERS,
		"users.events",
		"tasks.service",
	)

	go consumer.Start(context.TODO())

	return &App{
		Config:          config,
		TaskOAPIHandler: oapiTaskHandler,
		TaskGRPCHandler: grpcTaskHandler,
		TaskGRPCServer:  taskGRPCServer,
	}
}
