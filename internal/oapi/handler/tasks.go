package oapihandler

import (
	"net/http"

	"github.com/literaen/simple_project/tasks/internal/tasks"

	"github.com/gin-gonic/gin"
	dto "github.com/literaen/simple_project/dto"
)

type TaskHandler struct {
	service *tasks.TaskService
}

func NewTaskHandler(service *tasks.TaskService) *TaskHandler {
	handler := &TaskHandler{service: service}

	return handler
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while getting tasks"})
		return
	}

	resp := make([]*dto.Task, len(tasks))
	for i := 0; i < len(tasks); i++ {
		resp[i] = &dto.Task{
			ID:          tasks[i].ID,
			UserID:      tasks[i].UserID,
			Description: tasks[i].Description,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TaskHandler) PostTasks(c *gin.Context) {
	var dtoTask dto.Task
	if err := c.ShouldBindJSON(&dtoTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if dtoTask.UserID == 0 || dtoTask.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields (user_id, description) are required"})
		return
	}

	task := &tasks.Task{
		UserID:      dtoTask.UserID,
		Description: dtoTask.Description,
	}

	if err := h.service.PostTask(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtoTask.ID = task.ID

	c.JSON(http.StatusCreated, dtoTask)
}

func (h *TaskHandler) DeleteTasksId(c *gin.Context, id uint64) {
	err := h.service.DeleteTaskByID(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) PatchTasksId(c *gin.Context, id uint64) {
	var dtoTask dto.Task
	if err := c.ShouldBindJSON(&dtoTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	task := &tasks.Task{}

	if dtoTask.UserID != 0 {
		task.UserID = dtoTask.UserID
	}

	if dtoTask.Description != "" {
		task.Description = dtoTask.Description
	}

	task, err := h.service.PatchTaskByID(id, task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Task{
		ID:          task.ID,
		UserID:      task.UserID,
		Description: task.Description,
	})
}
