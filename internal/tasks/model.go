package tasks

import "gorm.io/gorm"

type Task struct {
	ID          uint64 `gorm:"primaryKey"`
	UserID      uint64
	Description string
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Task{})
}
