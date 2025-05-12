package chat

import (
	"time"
)

type Chat struct {
	IDMessage  int       `gorm:"column:id_message;primaryKey;autoIncrement"`
	Form       string    `gorm:"column:form"`
	Entity     string    `gorm:"column:entity;type:ENUM('personal', 'role')"`
	Dest       string    `gorm:"column:dest"`
	Subject    string    `gorm:"column:subject"`
	Message    string    `gorm:"column:message"`
	Status     string    `gorm:"column:status;type:ENUM('mengirim', 'terkirim', 'dibaca');default:'mengirim'"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for the Chat model
func (Chat) TableName() string {
	return "chat"
}

type ChatRequest struct {
	Notifications []ChatNotification `json:"notifications"`
}

type ChatNotification struct {
	Form    string `json:"form"`
	Entity  string `json:"entity"`
	Dest    string `json:"dest"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
