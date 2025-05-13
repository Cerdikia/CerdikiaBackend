package repositories

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"fmt"
)

func GetSiswaBeingVerified() (*[]users.VerifiedUserResponse, string) {
	var verifiedUsers []users.VerifiedUserResponse

	db := config.DB
	query := `
		SELECT 
			uv.email, 
			uv.verified_status, 
			s.nama, 
			s.id_kelas,
			k.kelas
		FROM 
			user_verified uv
		LEFT JOIN 
			siswa s ON uv.email = s.email
		LEFT JOIN 
			kelas k ON s.id_kelas = k.id_kelas;
	`

	tmpResult := db.Raw(query).Scan(&verifiedUsers)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, tmpResult.Error.Error()
	} else {
		return &verifiedUsers, "Data retrieved successfully"
	}
}

func GetSiswaBeingVerifiedByEmail(email string) (*users.VerifiedUserResponse, string) {
	var verifiedUser users.VerifiedUserResponse

	db := config.DB
	query := `
		SELECT 
			uv.email, 
			uv.verified_status, 
			s.nama, 
			s.id_kelas,
			k.kelas
		FROM 
			user_verified uv
		LEFT JOIN 
			siswa s ON uv.email = s.email
		LEFT JOIN 
			kelas k ON s.id_kelas = k.id_kelas
		WHERE 
			uv.email = ?;
	`

	tmpResult := db.Raw(query, email).Scan(&verifiedUser)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, tmpResult.Error.Error()
	} else {
		return &verifiedUser, "Data retrieved successfully"
	}
}
