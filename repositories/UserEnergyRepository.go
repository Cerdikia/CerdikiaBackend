package repositories

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"fmt"
	"time"
)

func CreateUserEnergyFirstTime(email string) string {
	db := config.DB

	now := time.Now()
	// now := time.Now().UTC()
	newUser := users.UserEnergy{
		Email:       email,
		Energy:      5,
		LastUpdated: now,
	}

	if err := db.Create(&newUser).Error; err != nil {
		fmt.Println(err)
		return "error : Failed to create user"
	}
	return "User energy initialized successfuly"

}
