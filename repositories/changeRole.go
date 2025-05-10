package repositories

import (
	"coba1BE/config"
	"coba1BE/models/users"
	"fmt"
)

// ChangeUserRole changes a user's role by moving their data between tables
// It copies the user data from the source table to the target table and deletes from the source table
// For guru and kepalaSekolah roles (which use the same table), it updates the jabatan field instead
func ChangeUserRole(email string, oldRole, newRole string) (*users.UserProfile, string) {
	db := config.DB
	var userProfile users.UserProfile

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Special case: changing between guru and kepalaSekolah (same table, different jabatan)
	if (oldRole == "guru" && newRole == "kepalaSekolah") || (oldRole == "kepalaSekolah" && newRole == "guru") {
		// Get the current user data
		sourceQuery := `SELECT id, email, nama, NULL AS id_kelas, jabatan, NULL AS keterangan, date_created, image_profile FROM guru WHERE email = ?`
		result := tx.Raw(sourceQuery, email).Scan(&userProfile)
		if result.Error != nil {
			tx.Rollback()
			return nil, fmt.Sprintf("Error retrieving user data: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, "User not found in the source role"
		}

		// Update the jabatan field based on the new role
		var newJabatan string
		if newRole == "kepalaSekolah" {
			newJabatan = "kepala sekolah"
		} else {
			// If changing from kepalaSekolah to guru, set a default jabatan or empty string
			newJabatan = "guru" // You can adjust this default value as needed
		}

		// Update the jabatan in the guru table
		updateQuery := `UPDATE guru SET jabatan = ? WHERE email = ?`
		result = tx.Exec(updateQuery, newJabatan, email)
		if result.Error != nil {
			tx.Rollback()
			return nil, fmt.Sprintf("Error updating jabatan: %v", result.Error)
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Sprintf("Error committing transaction: %v", err)
		}

		// Update the role and jabatan in the response
		userProfile.Role = newRole
		jabatan := newJabatan
		userProfile.Jabatan = &jabatan

		return &userProfile, "Role changed successfully"
	}

	// For other role changes that involve different tables:
	// 1. Get user data from the source table based on oldRole
	var sourceQuery string
	switch oldRole {
	case "siswa":
		sourceQuery = `SELECT email, nama, id_kelas, NULL AS jabatan, NULL AS keterangan, date_created, image_profile FROM siswa WHERE email = ?`
	case "guru", "kepalaSekolah":
		sourceQuery = `SELECT id, email, nama, NULL AS id_kelas, jabatan, NULL AS keterangan, date_created, image_profile FROM guru WHERE email = ?`
	case "admin":
		sourceQuery = `SELECT email, nama, NULL AS id_kelas, NULL AS jabatan, keterangan, date_created, image_profile FROM admin WHERE email = ?`
	default:
		return nil, "Invalid source role"
	}

	// Execute the query to get user data
	result := tx.Raw(sourceQuery, email).Scan(&userProfile)
	if result.Error != nil {
		tx.Rollback()
		return nil, fmt.Sprintf("Error retrieving user data: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, "User not found in the source role"
	}

	// 2. Insert user data into the target table based on newRole
	var insertQuery string
	var insertArgs []interface{}

	switch newRole {
	case "siswa":
		insertQuery = `INSERT INTO siswa (email, nama, image_profile) VALUES (?, ?, ?)`
		insertArgs = []interface{}{userProfile.Email, userProfile.Nama, userProfile.ImageProfile}
	case "guru":
		insertQuery = `INSERT INTO guru (email, nama, image_profile) VALUES (?, ?, ?)`
		insertArgs = []interface{}{userProfile.Email, userProfile.Nama, userProfile.ImageProfile}
	case "kepalaSekolah":
		// For kepalaSekolah, we insert into guru table with jabatan="kepala sekolah"
		insertQuery = `INSERT INTO guru (email, nama, jabatan, image_profile) VALUES (?, ?, ?, ?)`
		insertArgs = []interface{}{userProfile.Email, userProfile.Nama, "kepala sekolah", userProfile.ImageProfile}
	case "admin":
		insertQuery = `INSERT INTO admin (email, nama, image_profile) VALUES (?, ?, ?)`
		insertArgs = []interface{}{userProfile.Email, userProfile.Nama, userProfile.ImageProfile}
	default:
		tx.Rollback()
		return nil, "Invalid target role"
	}

	// Execute the insert query
	result = tx.Exec(insertQuery, insertArgs...)
	if result.Error != nil {
		tx.Rollback()
		return nil, fmt.Sprintf("Error inserting user data to new role: %v", result.Error)
	}

	// 3. Delete user data from the source table
	var deleteQuery string
	switch oldRole {
	case "siswa":
		deleteQuery = `DELETE FROM siswa WHERE email = ?`
	case "guru", "kepalaSekolah":
		deleteQuery = `DELETE FROM guru WHERE email = ?`
	case "admin":
		deleteQuery = `DELETE FROM admin WHERE email = ?`
	default:
		tx.Rollback()
		return nil, "Invalid source role"
	}

	// Execute the delete query
	result = tx.Exec(deleteQuery, email)
	if result.Error != nil {
		tx.Rollback()
		return nil, fmt.Sprintf("Error deleting user data from old role: %v", result.Error)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Sprintf("Error committing transaction: %v", err)
	}

	// Update the role in the response
	userProfile.Role = newRole

	return &userProfile, "Role changed successfully"
}
