package repositories

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"fmt"
)

func GetSiswaBeingVerified() (*[]users.UserVerified, string) {
	var beingverifiedAcout []users.UserVerified

	db := config.DB
	query :=
		`SELECT email, verified FROM user_verified;`

	tmpResult := db.Raw(query).Scan(&beingverifiedAcout)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, tmpResult.Error.Error()
	} else {
		return &beingverifiedAcout, "Data retrieved successfully"
	}
}

func GetSiswaBeingVerifiedByEmail(email string) (*users.UserVerified, string) {
	var beingverifiedAcout users.UserVerified

	db := config.DB
	query :=
		`SELECT email, verified FROM user_verified WHERE email = ?;`

	tmpResult := db.Raw(query, email).Scan(&beingverifiedAcout)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, tmpResult.Error.Error()
	} else {
		return &beingverifiedAcout, "Data retrieved successfully"
	}
}
