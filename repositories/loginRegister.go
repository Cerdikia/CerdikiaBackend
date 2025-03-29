package repositories

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/users"
	"fmt"
)

func GetAllUsers() models.BaseResponseModel {
	var users []users.User
	var result models.BaseResponseModel

	db := config.DB
	query := `SELECT email, nama, kelas, date_created, role FROM users`

	tmpResult := db.Raw(query).Scan(&users)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		result = models.BaseResponseModel{
			Message: tmpResult.Error.Error(),
			Data:    nil,
		}
	} else {
		result = models.BaseResponseModel{
			Message: "Data retrieved successfully",
			Data:    users,
		}
	}

	return result
}

func GetUserByEmail(email string) models.BaseResponseModel {
	var users users.User
	var result models.BaseResponseModel

	db := config.DB
	query := `SELECT email, nama, kelas, date_created, role FROM users WHERE email = ?`

	tmpResult := db.Raw(query, email).Scan(&users)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		result = models.BaseResponseModel{
			Message: tmpResult.Error.Error(),
			Data:    nil,
		}
	} else if tmpResult.RowsAffected == 0 {
		result = models.BaseResponseModel{
			Message: "no data found",
			Data:    nil,
		}
	} else {
		// fmt.Println("email : " + users[0].Email)
		result = models.BaseResponseModel{
			Message: "Data retrieved successfully",
			Data:    users,
		}
	}

	return result
}
