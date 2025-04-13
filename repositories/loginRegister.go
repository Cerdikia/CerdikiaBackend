package repositories

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/users"
	"fmt"
)

func GetAllSiswa() models.BaseResponseModel {
	var users []users.Siswa
	var result models.BaseResponseModel

	db := config.DB
	query := `SELECT email, nama, kelas, date_created FROM siswa`

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

func GetDataActor(role string) models.BaseResponseModel {
	var UserProfile []users.UserProfile
	var result models.BaseResponseModel
	var query string

	db := config.DB

	switch role {
	case "siswa":
		query = `SELECT email, nama, kelas, date_created FROM siswa`
	case "guru":
		query = `SELECT email, nama, id_mapel, date_created FROM guru`
	case "admin":
		query = `SELECT email, nama, keterangan, date_created FROM admin`
	default:
		result = models.BaseResponseModel{
			Message: "undifine role",
			Data:    nil,
		}
		return result
	}
	// query := `SELECT email, nama, kelas, date_created FROM siswa`

	// tmpResult := db.Raw(query).Scan(&users)
	tmpResult := db.Raw(query).Scan(&UserProfile)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		result = models.BaseResponseModel{
			Message: tmpResult.Error.Error(),
			Data:    nil,
		}
	} else {
		for i := range UserProfile {
			UserProfile[i].Role = role
		}
		result = models.BaseResponseModel{
			Message: "Data retrieved successfully",
			Data:    UserProfile,
		}
	}

	return result
}

// () models.BaseResponseModel {
// 	var users []users.Siswa
// 	var result models.BaseResponseModel

// 	db := config.DB
// 	query := `SELECT email, nama, kelas, date_created FROM siswa`

// 	tmpResult := db.Raw(query).Scan(&users)

// 	if tmpResult.Error != nil {
// 		fmt.Println(tmpResult.Error)
// 		result = models.BaseResponseModel{
// 			Message: tmpResult.Error.Error(),
// 			Data:    nil,
// 		}
// 	} else {
// 		result = models.BaseResponseModel{
// 			Message: "Data retrieved successfully",
// 			Data:    users,
// 		}
// 	}

// 	return result
// }

func GetAllUsers() models.BaseResponseModel {
	var users []users.UserProfile
	var result models.BaseResponseModel

	db := config.DB
	query :=
		`SELECT email, nama, 'siswa' AS role, date_created FROM siswa
UNION
SELECT email, nama, 'guru' AS role, date_created FROM guru
UNION
SELECT email, nama, 'admin' AS role, date_created FROM admin;`

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

func GetUserByEmail(email, role string) models.BaseResponseModel {
	// var users users.Siswa
	var users users.UserProfile
	var result models.BaseResponseModel
	var query string
	fmt.Println("loginRequest role : " + role)

	db := config.DB
	// query := `SELECT email, nama, kelas, date_created FROM siswa WHERE email = ?`
	switch role {
	case "siswa":
		query = `SELECT email, nama, kelas, NULL AS jabatan, NULL AS keterangan, date_created FROM siswa WHERE email = ?`
	case "guru":
		query = `SELECT email, nama, NULL AS kelas, jabatan, NULL AS keterangan, date_created FROM guru WHERE email = ?`
	case "admin":
		query = `SELECT email, nama, NULL AS kelas, NULL AS jabatan, keterangan, date_created FROM admin WHERE email = ?`
	default:
		result = models.BaseResponseModel{
			Message: "no data found",
			Data:    nil,
		}

		return result
	}

	// tmpResult := db.Raw(query, role, email).Scan(&users)
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
		// Tambahkan asal
		users.Role = role
		// fmt.Println("email : " + users[0].Email)
		result = models.BaseResponseModel{
			Message: "Data retrieved successfully",
			Data:    users,
		}
	}

	return result
}

func UpdateDataSiswa(actor users.Siswa) (*users.Siswa, string) {
	var query string
	// var result users.UserProfile
	db := config.DB
	query = "UPDATE siswa SET nama = ? , kelas = ? WHERE email = ?"

	tempResult := db.Exec(query, actor.Nama, actor.Kelas, actor.Email)

	if tempResult.Error != nil {
		return nil, fmt.Sprintf("error query : %e", tempResult.Error)
	} else {
		rowsAffected := tempResult.RowsAffected
		if rowsAffected == 0 {
			return nil, "Tidak ada data yang di temukan atau data tidak berubah"
		} else {
			return &actor, "Success"
		}
	}
}

func UpdateDataGuru(actor users.Guru) (*users.Guru, string) {
	var query string
	// var result users.UserProfile
	db := config.DB
	query = "UPDATE guru SET id_mapel = ?, nama = ?, jabatan = ? WHERE email = ?"

	tempResult := db.Exec(query, actor.IDMapel, actor.Nama, actor.Jabatan, actor.Email)

	if tempResult.Error != nil {
		return nil, fmt.Sprintf("error query : %e", tempResult.Error)
	} else {
		rowsAffected := tempResult.RowsAffected
		if rowsAffected == 0 {
			return nil, "Tidak ada data yang di temukan atau data tidak berubah"
		} else {
			return &actor, "Success"
		}
	}
}

func UpdateDataAdmin(actor users.Admin) (*users.Admin, string) {
	var query string
	// var result users.UserProfile
	db := config.DB
	query = "UPDATE admin SET nama = ?, keterangan = ? WHERE email = ?"

	tempResult := db.Exec(query, actor.Nama, actor.Keterangan, actor.Email)

	if tempResult.Error != nil {
		return nil, fmt.Sprintf("error query : %e", tempResult.Error)
	} else {
		rowsAffected := tempResult.RowsAffected
		if rowsAffected == 0 {
			return nil, "Tidak ada data yang di temukan atau data tidak berubah"
		} else {
			return &actor, "Success"
		}
	}
}
