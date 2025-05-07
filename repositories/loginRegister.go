package repositories

import (
	"coba1BE/config"
	"coba1BE/models"
	"coba1BE/models/points"
	"coba1BE/models/users"
	"fmt"
)

func GetAllSiswa() models.BaseResponseModel {
	var users []users.Siswa
	var result models.BaseResponseModel

	db := config.DB
	query := `SELECT email, nama, id_kelas, date_created FROM siswa`

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
		query = `SELECT email, nama, id_kelas, date_created FROM siswa`
	case "guru":
		query = `SELECT id, email, nama, jabatan, date_created FROM guru`
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

func GetDataActorByRoleAndEmail(role, email string) models.BaseResponseModel {
	var UserProfile []users.UserProfile
	var result models.BaseResponseModel
	var query string

	db := config.DB

	switch role {
	case "siswa":
		query = `SELECT email, nama, id_kelas, date_created FROM siswa WHERE email = ?`
	case "guru":
		query = `SELECT id, email, nama, jabatan, date_created FROM guru WHERE email = ?`
	case "admin":
		query = `SELECT email, nama, keterangan, date_created FROM admin WHERE email = ?`
	default:
		result = models.BaseResponseModel{
			Message: "undifine role",
			Data:    nil,
		}
		return result
	}
	// query := `SELECT email, nama, kelas, date_created FROM siswa`

	// tmpResult := db.Raw(query).Scan(&users)
	tmpResult := db.Raw(query, email).Scan(&UserProfile)

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

func GetUserByEmail(email, role string) (*users.UserProfile, string) {
	// var users users.Siswa
	var user users.UserProfile
	var query string
	fmt.Println("loginRequest role : " + role)

	db := config.DB
	// query := `SELECT email, nama, kelas, date_created FROM siswa WHERE email = ?`
	switch role {
	case "siswa":
		query = `SELECT email, nama, id_kelas, NULL AS jabatan, NULL AS keterangan, date_created, image_profile FROM siswa WHERE email = ?`
	case "guru":
		query = `SELECT id, email, nama, NULL AS id_kelas, jabatan, NULL AS keterangan, date_created, image_profile FROM guru WHERE email = ?`
	case "admin":
		query = `SELECT email, nama, NULL AS id_kelas, NULL AS jabatan, keterangan, date_created, image_profile FROM admin WHERE email = ?`
	default:
		return nil, "error bad request"
	}

	// tmpResult := db.Raw(query, role, email).Scan(&user)
	tmpResult := db.Raw(query, email).Scan(&user)

	if tmpResult.Error != nil {
		fmt.Println(tmpResult.Error)
		return nil, fmt.Sprintf("error query data : %e", tmpResult.Error)

	} else if tmpResult.RowsAffected == 0 {
		return nil, "no data found"
	} else {
		// Tambahkan asal
		user.Role = role
		// fmt.Println("email : " + user[0].Email)
		return &user, "Data retrieved successfully"
	}
}

func UpdateDataSiswa(actor users.Siswa) (*users.Siswa, string) {
	var query string
	// var result users.UserProfile
	db := config.DB
	query = "UPDATE siswa SET nama = ? , id_kelas = ? WHERE email = ?"

	tempResult := db.Exec(query, actor.Nama, actor.IdKelas, actor.Email)

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
	query = "UPDATE guru SET nama = ?, jabatan = ? WHERE email = ?"

	tempResult := db.Exec(query, actor.Nama, actor.Jabatan, actor.Email)

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

func CreatePointFirst(email string) error {
	userPoint := points.UserPoint{
		Email: email,
	}
	db := config.DB
	err := db.Create(&userPoint).Error
	return err
}

// func UpdatPointGeneral() {

// }

func CreateAcountVerifiedFirst(email string) error {
	userVerivied := users.UserVerified{
		Email: email,
	}
	db := config.DB
	err := db.Create(&userVerivied).Error
	return err
}
