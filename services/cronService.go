package services

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func StartEnergyCron() error {
	db := config.DB
	c := cron.New()

	// Tangkap error dari AddFunc
	err := c.AddFunc("@every 10m", func() {
		var users []users.UserEnergy
		if err := db.Find(&users).Error; err != nil {
			fmt.Println("Error fetching users:", err)
			return // hanya keluar dari fungsi cron, tidak dari startEnergyCron
		}

		for _, user := range users {
			if user.Energy < 5 {
				user.Energy += 1
				if user.Energy > 5 {
					user.Energy = 5
				}
				user.LastUpdated = time.Now()
				db.Save(&user)
			}
		}
		fmt.Println("Energy updated at", time.Now())
	})

	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	c.Start()
	return nil
}

func AddEnergy(email string) error {
	db := config.DB

	var user users.UserEnergy
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		fmt.Println(err)
		return fmt.Errorf("User tidak ditemukan")
		// c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan", "success": false})
	}

	if user.Energy < 5 {
		user.Energy += 1
		now := time.Now()
		user.LastUpdated = now
		db.Save(&user)
	} else {
		return fmt.Errorf("energy is at maximum")
	}

	return nil
}

func UseEnergy(email string) error {
	db := config.DB

	var user users.UserEnergy
	// if err := db.First(&user, "email = ?", email).Error; err != nil {
	if err := db.Take(&user, "email = ?", email).Error; err != nil {

		// c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan", "success": false})
		return fmt.Errorf("User tidak ditemukan")
	}

	if user.Energy > 0 {
		user.Energy -= 1
		now := time.Now()
		user.LastUpdated = now
		db.Save(&user)
	} else {
		return fmt.Errorf("Energy 0, Please wait 10 minuts to recharge")
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Energi berhasil dikurangi", "data": user, "success": true})
	return nil
}
