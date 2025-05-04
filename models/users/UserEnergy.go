package users

import "time"

type UserEnergy struct {
	Email       string    `gorm:"column:email;primaryKey"`
	Energy      int       `gorm:"column:energy"`
	LastUpdated time.Time `gorm:"column:last_updated"`
}

// TableName override untuk menentukan nama tabel secara manual
func (UserEnergy) TableName() string {
	return "user_energy"
}
