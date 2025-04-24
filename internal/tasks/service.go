package tasks

import (
	"context"
	"time"

	"github.com/literaen/simple_project/tasks/internal/users"
)

type TaskService struct {
	repo       TaskRepository
	userClient *users.UserService
}

func NewTaskService(repo TaskRepository, userClient *users.UserService) *TaskService {
	return &TaskService{repo: repo, userClient: userClient}
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetUserAllTasks(id uint64) ([]Task, error) {
	return s.repo.GetUserAllTasks(id)
}

func (s *TaskService) GetTaskByID(id uint64) (*Task, error) {
	return s.repo.GetTaskByID(id)
}

func (s *TaskService) PostTask(task *Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, существует ли пользователь
	if err := s.userClient.GetUser(ctx, task.UserID); err != nil {
		return err
	}

	return s.repo.PostTask(task)
}

func (s *TaskService) PatchTaskByID(id uint64, task *Task) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, существует ли пользователь
	if task.UserID != 0 {
		if err := s.userClient.GetUser(ctx, task.UserID); err != nil {
			return nil, err
		}
	}

	return s.repo.PatchTaskByID(id, task)
}

func (s *TaskService) DeleteTaskByID(id uint64) error {
	return s.repo.DeleteTaskByID(id)
}

func (s *TaskService) DeleteTasksByUserID(userID uint64) error {
	return s.repo.DeleteTasksByUserID(userID)
}
