package controllers

import (
	"coba1BE/models"
	"coba1BE/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChangeUserRoleRequest defines the request structure for changing a user's role
type ChangeUserRoleRequest struct {
	Email   string `json:"email" binding:"required"`
	OldRole string `json:"old_role" binding:"required"`
	NewRole string `json:"new_role" binding:"required"`
}

// ChangeUserRole handles the request to change a user's role
func ChangeUserRole(c *gin.Context) {
	var request ChangeUserRoleRequest

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid request format",
			Data:    nil,
		})
		return
	}

	// Validate roles
	validRoles := map[string]bool{
		"siswa":        true,
		"guru":         true,
		"admin":        true,
		"kepalaSekolah": true,
	}

	if !validRoles[request.OldRole] || !validRoles[request.NewRole] {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "Invalid role specified",
			Data:    nil,
		})
		return
	}

	// Don't allow changing to the same role
	if request.OldRole == request.NewRole {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: "New role must be different from old role",
			Data:    nil,
		})
		return
	}

	// Call the repository function to change the role
	userProfile, message := repositories.ChangeUserRole(request.Email, request.OldRole, request.NewRole)

	if userProfile == nil {
		c.JSON(http.StatusBadRequest, models.BaseResponseModel{
			Message: message,
			Data:    nil,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, models.BaseResponseModel{
		Message: message,
		Data:    userProfile,
	})
}
