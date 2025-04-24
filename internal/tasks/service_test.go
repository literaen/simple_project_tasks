package tasks

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/literaen/simple_project/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*postgres.GDB, func()) {
	// Создаем in-memory SQLite базу
	db, err := gorm.Open(sqlite.Open("file::memory:?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to sqlite in-memory DB: %v", err)
	}

	// Миграции — необходимо для создания нужных таблиц
	if err := db.AutoMigrate(&Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Возвращаем закрытие базы данных в конце теста
	return &postgres.GDB{DB: db}, func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get sql.DB instance: %v", err)
		}
		sqlDB.Close()
	}
}

func TestDeleteTasksByUserID(t *testing.T) {
	// Настроим тестовую БД
	gdb, teardown := setupTestDB(t)
	defer teardown()

	// Создаем репозиторий с in-memory SQLite
	taskRepo := NewTaskRepository(gdb, nil)
	taskService := NewTaskService(taskRepo, nil)

	// Заполняем тестовыми данными
	task1 := Task{UserID: 1, Description: "Task 1"}
	task2 := Task{UserID: 1, Description: "Task 2"}
	task3 := Task{UserID: 2, Description: "Task 3"}

	// Сохраняем задачи в БД
	gdb.DB.Create(&task1)
	gdb.DB.Create(&task2)
	gdb.DB.Create(&task3)

	// Проверяем, что задачи существуют
	var count int64
	gdb.DB.Model(&Task{}).Where("user_id = ?", 1).Count(&count)
	assert.Equal(t, int64(2), count, "Expected 2 tasks for user with ID 1")

	// Удаляем задачи пользователя с ID 1
	err := taskService.DeleteTasksByUserID(1)
	assert.NoError(t, err, "Expected no error during task deletion")

	// Проверяем, что задачи пользователя с ID 1 были удалены
	gdb.DB.Model(&Task{}).Where("user_id = ?", 1).Count(&count)
	assert.Equal(t, int64(0), count, "Expected 0 tasks for user with ID 1")

	// Проверяем, что задачи другого пользователя остались
	gdb.DB.Model(&Task{}).Where("user_id = ?", 2).Count(&count)
	assert.Equal(t, int64(1), count, "Expected 1 task for user with ID 2")
}
