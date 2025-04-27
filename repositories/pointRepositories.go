package repositories

import (
	"coba1BE/config"
	"coba1BE/models/points"
	"fmt"
)

// func UpdatePointGeneral(email string, point points.DiamondOrExp) error {
// 	// var result users.UserProfile
// 	db := config.DB
// 	// GORM akan hanya update field yang tidak nil
// 	err := db.Model(&points.UserPoint{}).
// 		Where("email = ?", email).
// 		Updates().Error

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User point updated"})
// }

func GetUserPoint(email string) (*points.UserPointResponse, string) {
	var userPoint points.UserPointResponse

	query := `SELECT email, diamond, exp FROM user_points WHERE email = ?`
	db := config.DB

	// tmpResult := db.Raw(query, role, email).Scan(&user)
	tmpResult := db.Raw(query, email).Scan(&userPoint)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error query data : %e", tmpResult.Error)

	} else if tmpResult.RowsAffected == 0 {
		return nil, "no data found"
	} else {
		// fmt.Println("email : " + user[0].Email)
		return &userPoint, "Data retrieved successfully"
	}
}
