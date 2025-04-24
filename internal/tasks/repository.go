package tasks

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/literaen/simple_project/pkg/redis"
	"gorm.io/gorm"

	"github.com/literaen/simple_project/pkg/postgres"
)

type TaskRepository interface {
	// Получить все задания
	GetAllTasks() ([]Task, error)

	// Получить все задания пользователя
	GetUserAllTasks(userID uint64) ([]Task, error)

	// Получить задания по ID
	GetTaskByID(id uint64) (*Task, error)

	// Создать новое задание
	PostTask(user *Task) error

	// Изменить задание по ID
	PatchTaskByID(id uint64, user *Task) (*Task, error)

	// Удалить задание по ID
	DeleteTaskByID(id uint64) error

	// Удалить все задания по UserID
	DeleteTasksByUserID(userID uint64) error
}

type taskRepository struct {
	gdb   *postgres.GDB
	redis *redis.RDB
}

func NewTaskRepository(gdb *postgres.GDB, redis *redis.RDB) TaskRepository {
	return &taskRepository{gdb: gdb, redis: redis}
}

func taskIDtoUInt(s string) (uint64, error) {
	var id uint64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.gdb.DB.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetUserAllTasks(userID uint64) ([]Task, error) {
	var tasks []Task
	err := r.gdb.DB.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTaskByID(id uint64) (*Task, error) {
	var task Task

	taskKey := fmt.Sprintf("task:%d", id)
	taskData, err := r.redis.HGetAll(context.Background(), taskKey)

	if err == nil && len(taskData) > 0 {
		task.ID = id
		task.Description = taskData["description"]
		id, err := taskIDtoUInt(taskData["user_id"])
		if err != nil {
			return nil, err
		}
		task.UserID = id
	} else {
		if err := r.gdb.DB.First(&task, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("task with ID %d not found", id)
			}
			return nil, err
		}
	}

	return &task, nil
}

func (r *taskRepository) PostTask(task *Task) error {
	if err := r.gdb.DB.Create(task).Error; err != nil {
		return err
	}

	taskKey := fmt.Sprintf("task:%d", task.ID)
	taskData := map[string]interface{}{
		"description": task.Description,
		"user_id":     task.UserID,
	}

	err := r.redis.HSet(context.Background(), taskKey, taskData)
	if err != nil {
		log.Printf("error caching created task: %v", err)
	}

	return nil
}

func (r *taskRepository) PatchTaskByID(id uint64, task *Task) (*Task, error) {
	var resp *Task
	res := r.gdb.DB.
		Model(&Task{}).
		Where("id = ?", id).
		Updates(&task).
		Scan(&resp)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	return resp, nil
}

func (r *taskRepository) DeleteTaskByID(id uint64) error {
	result := r.gdb.DB.Where("id = ?", id).Delete(&Task{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return nil
}

func (r *taskRepository) DeleteTasksByUserID(userID uint64) error {
	result := r.gdb.DB.Where("user_id = ?", userID).Delete(&Task{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no tasks found for user ID %d", userID)
	}

	return nil
}
