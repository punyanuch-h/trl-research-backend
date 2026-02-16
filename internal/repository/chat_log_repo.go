package repository

import (
	"time"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type ChatLogRepository interface {
	CreateChatLog(chatLog *models.ChatLogs) error
}

type chatLogRepo struct {
	db *gorm.DB
}

func NewChatLogRepo(db *gorm.DB) ChatLogRepository {
	return &chatLogRepo{db: db}
}

func (r *chatLogRepo) CreateChatLog(chatLog *models.ChatLogs) error {
	id, err := utils.GenerateID(r.db, "chat_logs", "CL")
	if err != nil {
		return err
	}
	chatLog.ID = id
	chatLog.CreatedAt = time.Now()

	return r.db.Create(chatLog).Error
}
