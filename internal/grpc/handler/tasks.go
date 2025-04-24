package grpchandler

import (
	"context"

	taskpb "github.com/literaen/simple_project/proto/gen"

	"github.com/literaen/simple_project/tasks/internal/tasks"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	taskpb.UnimplementedTaskServiceServer
	service *tasks.TaskService
}

func NewTaskHandler(service *tasks.TaskService) *TaskHandler {
	handler := &TaskHandler{service: service}

	return handler
}

func (h *TaskHandler) GetAllTasks(ctx context.Context, req *taskpb.GetAllTasksRequest) (*taskpb.GetAllTasksResponse, error) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*taskpb.Task, len(tasks))
	for i := 0; i < len(tasks); i++ {
		resp[i] = &taskpb.Task{
			Id:          tasks[i].ID,
			UserId:      tasks[i].UserID,
			Description: tasks[i].Description,
		}
	}

	return &taskpb.GetAllTasksResponse{
		Tasks: resp,
	}, nil
}

func (h *TaskHandler) GetUserAllTasks(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetAllTasksResponse, error) {
	id := req.GetId()
	tasks, err := h.service.GetUserAllTasks(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*taskpb.Task, len(tasks))
	for i := 0; i < len(tasks); i++ {
		resp[i] = &taskpb.Task{
			Id:          tasks[i].ID,
			UserId:      tasks[i].UserID,
			Description: tasks[i].Description,
		}
	}

	return &taskpb.GetAllTasksResponse{
		Tasks: resp,
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	id := req.GetId()
	task, err := h.service.GetTaskByID(id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.GetTaskResponse{
		Task: &taskpb.Task{
			Id:          id,
			UserId:      task.UserID,
			Description: task.Description,
		},
	}, nil
}

func (h *TaskHandler) AddTask(ctx context.Context, req *taskpb.AddTaskRequest) (*taskpb.AddTaskResponse, error) {
	data := req.GetTask()

	task := &tasks.Task{
		UserID:      data.GetUserId(),
		Description: data.GetDescription(),
	}

	if err := h.service.PostTask(task); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.AddTaskResponse{
		Id: task.ID,
	}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	id := req.GetId()
	err := h.service.DeleteTaskByID(id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	id := req.GetId()
	data := req.GetTask()

	task, err := h.service.PatchTaskByID(id, &tasks.Task{
		UserID:      data.GetUserId(),
		Description: data.GetDescription(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:          task.ID,
			UserId:      task.UserID,
			Description: task.Description,
		},
	}, nil
}
