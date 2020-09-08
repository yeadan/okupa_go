package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//UserOkupa struct de los UserOkupas
type UserOkupa struct {
	UserOkupaID int     `gorm:"primary_key" json:"userokupa_id"`
	UserID    	int     `gorm:"not null;type:int REFERENCES users(user_id) ON DELETE CASCADE" json:"user_id"`
	OkupaID 	int     `gorm:"not null;type:int REFERENCES okupas(okupa_id) ON DELETE CASCADE" json:"okupa_id"`
	Created   time.Time `json:"created"`
}

// GetUserOkupaUser - Busca okupas por id de user
func GetUserOkupaUser(user int, db *gorm.DB) []*UserOkupa {
	okupas := []*UserOkupa{}
	db.Where("user_id = ?", user).Find(&okupas)
	return okupas
}

// GetUserOkupaOkupa - Busca usuarios por id de asociación okupa
func GetUserOkupaOkupa(okupa int, db *gorm.DB) []*UserOkupa {
	okupas := []*UserOkupa{}
	db.Where("okupa_id = ?", okupa).Find(&okupas)
	return okupas
}

// GetUserOkupa - Busca por id
func GetUserOkupa(id int, db *gorm.DB) *UserOkupa {
	userOkupa := new(UserOkupa)
	db.Find(userOkupa, id)
	if userOkupa.UserOkupaID == id {
		return userOkupa
	}
	return nil
}
// Mira si ya existe esa relación, y la devuelve si existe.
func ExistUserOkupa(usr int, okp int, db *gorm.DB) *UserOkupa {
	userOkupa := new(UserOkupa)
	result := db.Where("user_id = ? AND okupa_id = ?", usr,okp).First(&userOkupa)
	if result.RecordNotFound() {
		return nil
	}
	return userOkupa
}

func CreateUserOkupa(userOkupa *UserOkupa, db *gorm.DB) {
	db.Create(userOkupa)
}
func DeleteUserOkupa(del *UserOkupa, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Delete(del)
}

//EditUserOkupa no utilizado
func EditUserOkupa(editUserOkupa *UserOkupa, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Save(editUserOkupa)
}