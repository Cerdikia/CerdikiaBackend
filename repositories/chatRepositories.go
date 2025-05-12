package repositories

import (
	"coba1BE/config"
	"coba1BE/models/chat"
	"fmt"
)

// CreateChatMessages creates multiple chat messages in the database
func CreateChatMessages(messages []chat.ChatNotification) ([]chat.Chat, error) {
	db := config.DB
	var createdMessages []chat.Chat

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, msg := range messages {
		newMessage := chat.Chat{
			Form:    msg.Form,
			Entity:  msg.Entity,
			Dest:    msg.Dest,
			Subject: msg.Subject,
			Message: msg.Message,
			Status:  msg.Status,
		}

		if err := tx.Create(&newMessage).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create message: %w", err)
		}

		createdMessages = append(createdMessages, newMessage)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdMessages, nil
}

// GetAllChatMessages gets all chat messages with optional filtering
func GetAllChatMessages(filters map[string]string) ([]chat.Chat, error) {
	db := config.DB
	var messages []chat.Chat

	query := db.Model(&chat.Chat{}).Order("created_at DESC")

	// Apply filters if they exist
	if form, exists := filters["form"]; exists && form != "" {
		query = query.Where("form = ?", form)
	}

	if dest, exists := filters["dest"]; exists && dest != "" {
		query = query.Where("dest = ?", dest)
	}

	if subject, exists := filters["subject"]; exists && subject != "" {
		query = query.Where("subject LIKE ?", "%"+subject+"%")
	}

	if entity, exists := filters["entity"]; exists && entity != "" {
		query = query.Where("entity = ?", entity)
	}

	if status, exists := filters["status"]; exists && status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

// GetChatMessagesByDest gets all chat messages for a specific recipient
func GetChatMessagesByDest(dest string) ([]chat.Chat, error) {
	db := config.DB
	var messages []chat.Chat

	if err := db.Where("dest = ?", dest).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

// GetChatMessagesByForm gets all chat messages from a specific sender
func GetChatMessagesByForm(form string) ([]chat.Chat, error) {
	db := config.DB
	var messages []chat.Chat

	if err := db.Where("form = ?", form).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

// GetChatMessagesBySubject gets all chat messages with a specific subject
func GetChatMessagesBySubject(subject string) ([]chat.Chat, error) {
	db := config.DB
	var messages []chat.Chat

	if err := db.Where("subject LIKE ?", "%"+subject+"%").Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

// GetChatMessageByID gets a chat message by its ID
func GetChatMessageByID(id int) (*chat.Chat, error) {
	db := config.DB
	var message chat.Chat

	if err := db.Where("id_message = ?", id).First(&message).Error; err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

// CountUnreadMessages counts the number of unread messages for a specific recipient
func CountUnreadMessages(dest string) (int64, error) {
	db := config.DB
	var count int64

	// Count messages where status is 'terkirim' (not read yet)
	query := db.Model(&chat.Chat{}).Where("dest = ? AND status = 'terkirim'", dest)
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	return count, nil
}

// CountAllUnreadMessages counts all unread messages in the system
func CountAllUnreadMessages() (int64, error) {
	db := config.DB
	var count int64

	// Count all messages where status is 'terkirim' (not read yet)
	query := db.Model(&chat.Chat{}).Where("status = 'terkirim'")
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count all unread messages: %w", err)
	}

	return count, nil
}

// GetUnreadMessagesByRecipient gets a summary of unread messages grouped by recipient
func GetUnreadMessagesByRecipient() ([]map[string]interface{}, error) {
	db := config.DB
	var results []map[string]interface{}

	// Execute raw SQL query to get counts grouped by recipient
	rows, err := db.Raw(`
		SELECT dest as email, COUNT(*) as count 
		FROM chat 
		WHERE status = 'terkirim' 
		GROUP BY dest 
		ORDER BY count DESC
	`).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to get unread messages by recipient: %w", err)
	}
	defer rows.Close()

	// Process the results
	for rows.Next() {
		var email string
		var count int64
		if err := rows.Scan(&email, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, map[string]interface{}{
			"email": email,
			"count": count,
		})
	}

	return results, nil
}

// UpdateChatMessageStatus updates the status of a chat message
func UpdateChatMessageStatus(id int, status string) error {
	db := config.DB

	if status != "mengirim" && status != "terkirim" && status != "dibaca" {
		return fmt.Errorf("invalid status: %s", status)
	}

	if err := db.Model(&chat.Chat{}).Where("id_message = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	return nil
}
