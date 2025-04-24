package consumer

import (
	"context"
	"log"

	"github.com/literaen/simple_project/tasks/internal/tasks"

	"github.com/literaen/simple_project/pkg/kafka"
)

type UserEvent struct {
	// EventType string `json:"event_type"`
	UserID uint64 `json:"user_id"`
}

type UserEventConsumer struct {
	kfr         *kafka.KFR
	TaskService *tasks.TaskService
}

func NewUserEventConsumer(taskService *tasks.TaskService, brokers []string, topic string, groupID string) *UserEventConsumer {
	return &UserEventConsumer{
		kfr:         kafka.NewKafkaReader(brokers, topic, groupID),
		TaskService: taskService,
	}
}

func (c *UserEventConsumer) Start(ctx context.Context) {
	log.Println("UserEventConsumer started...")

	for {
		var event UserEvent
		eventType, err := c.kfr.ReadMessage(ctx, &event)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		log.Printf("Received event: %+v", event)

		switch eventType {
		case "user.deleted":
			err := c.TaskService.DeleteTasksByUserID(event.UserID)
			if err != nil {
				log.Printf("Failed to delete tasks for user %d: %v", event.UserID, err)
			} else {
				log.Printf("Successfully deleted tasks for user %d", event.UserID)
			}

			log.Printf("user.deleted %d", event.UserID)
		case "user.updated":
			// обработка обновления
		default:
			log.Printf("Unknown event type: %s", eventType)
		}
	}
}
