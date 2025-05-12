package controllers

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/chat"
	"coba1BE/repositories"
	"coba1BE/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMessages handles the creation of multiple chat messages
func CreateMessages(c *gin.Context) {
	var request chat.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Error binding JSON:", err.Error())
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid request format: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Validate the request
	if len(request.Notifications) == 0 {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "No messages provided",
			Data:    nil,
		})
		return
	}

	// Set all statuses to 'terkirim' regardless of what was sent
	for i := range request.Notifications {
		// Override any provided status with 'terkirim'
		request.Notifications[i].Status = "terkirim"
	}

	// Create the messages
	createdMessages, err := repositories.CreateChatMessages(request.Notifications)
	if err != nil {
		fmt.Println("Error creating messages:", err.Error())
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to create messages: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.BaseResponseModel{
		Message: "Messages created successfully",
		Data:    createdMessages,
	})
}

// GetAllMessages gets all messages with optional filtering
func GetAllMessages(c *gin.Context) {
	// Extract filter parameters from query
	filters := make(map[string]string)
	
	// Check for form filter
	if form := c.Query("form"); form != "" {
		filters["form"] = form
	}
	
	// Check for dest filter
	if dest := c.Query("dest"); dest != "" {
		filters["dest"] = dest
	}
	
	// Check for subject filter
	if subject := c.Query("subject"); subject != "" {
		filters["subject"] = subject
	}
	
	// Check for entity filter
	if entity := c.Query("entity"); entity != "" {
		filters["entity"] = entity
	}
	
	// Check for status filter
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	
	fmt.Println("Filters applied:", filters)
	
	messages, err := repositories.GetAllChatMessages(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to get messages: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Messages retrieved successfully",
		Data:    messages,
	})
}

// GetMessagesByRecipient gets all messages for a specific recipient
func GetMessagesByRecipient(c *gin.Context) {
	dest := c.Param("dest")
	if dest == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Recipient email is required",
			Data:    nil,
		})
		return
	}

	messages, err := repositories.GetChatMessagesByDest(dest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to get messages: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Messages retrieved successfully",
		Data:    messages,
	})
}

// GetMessagesBySender gets all messages from a specific sender
func GetMessagesBySender(c *gin.Context) {
	form := c.Param("form")
	if form == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Sender email is required",
			Data:    nil,
		})
		return
	}

	messages, err := repositories.GetChatMessagesByForm(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to get messages: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Messages retrieved successfully",
		Data:    messages,
	})
}

// GetMessagesBySubject gets all messages with a specific subject
func GetMessagesBySubject(c *gin.Context) {
	subject := c.Param("subject")
	if subject == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Subject is required",
			Data:    nil,
		})
		return
	}

	messages, err := repositories.GetChatMessagesBySubject(subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to get messages: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Messages retrieved successfully",
		Data:    messages,
	})
}

// UpdateMessageStatus updates the status of a message
func UpdateMessageStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid message ID",
			Data:    nil,
		})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid request format",
			Data:    nil,
		})
		return
	}

	if statusUpdate.Status != "mengirim" && statusUpdate.Status != "terkirim" && statusUpdate.Status != "dibaca" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid status. Must be 'mengirim', 'terkirim', or 'dibaca'",
			Data:    nil,
		})
		return
	}

	err = repositories.UpdateChatMessageStatus(id, statusUpdate.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to update message status: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Message status updated successfully",
		Data:    nil,
	})
}

// MarkMessageAsRead marks a message as read (dibaca) when a user opens it
func MarkMessageAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid message ID",
			Data:    nil,
		})
		return
	}

	// Always set status to 'dibaca'
	err = repositories.UpdateChatMessageStatus(id, "dibaca")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to mark message as read: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Get the updated message to return in the response
	db := config.DB
	var message chat.Chat
	if err := db.First(&message, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Message marked as read but failed to retrieve updated message: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Message marked as read successfully",
		Data:    message,
	})
}

// GetChatMessageByIDAndMarkAsRead gets a chat message by ID and automatically marks it as read
func GetChatMessageByIDAndMarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid message ID",
			Data:    nil,
		})
		return
	}

	// First get the message
	message, err := repositories.GetChatMessageByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponseModel{
			Message: "Message not found: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Only update the status if it's currently 'terkirim'
	if message.Status == "terkirim" {
		// Update the status to 'dibaca'
		err = repositories.UpdateChatMessageStatus(id, "dibaca")
		if err != nil {
			// Even if updating fails, we still return the message
			fmt.Println("Failed to update message status:", err.Error())
		} else {
			// Update the status in our local copy to return the updated version
			message.Status = "dibaca"
		}
	}

	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Message retrieved successfully",
		Data:    message,
	})
}

// CountUnreadMessages counts the number of unread messages for a specific recipient
func CountUnreadMessages(c *gin.Context) {
	// Get the recipient email from the query parameter
	email := c.Query("email")
	if email == "" {
		// If no email provided, try to get it from the authorization header
		authHeader := c.GetHeader("Authorization")
		userEmail, err := services.ExtractEmailFromAuthHeader(authHeader)
		if err == nil && userEmail != "" {
			email = userEmail
			fmt.Println("Extracted email from token:", email)
		} else if err != nil {
			fmt.Println("Error extracting email from token:", err)
		}
	}

	// If still no email, return an error
	if email == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Email parameter is required",
			Data:    nil,
		})
		return
	}

	// Count unread messages for this recipient
	count, err := repositories.CountUnreadMessages(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to count unread messages: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Return the count in a structured format
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "Unread message count retrieved successfully",
		Data: map[string]interface{}{
			"email": email,
			"count": count,
		},
	})
}

// CountAllUnreadMessagesAdmin counts all unread messages in the system for admin
func CountAllUnreadMessagesAdmin(c *gin.Context) {
	// First, get total count of all unread messages
	totalCount, err := repositories.CountAllUnreadMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to count all unread messages: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Then, get breakdown by recipient
	recipientCounts, err := repositories.GetUnreadMessagesByRecipient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponseModel{
			Message: "Failed to get unread messages by recipient: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Return the counts in a structured format
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: "All unread message counts retrieved successfully",
		Data: map[string]interface{}{
			"total_count": totalCount,
			"by_recipient": recipientCounts,
		},
	})
}
