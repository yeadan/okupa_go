package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// Okupa - Struct con los datos de una asociacion de okupas.
type Okupa struct {
	OkupaID   int `gorm:"primary_key;not null" json:"okupa_id"`
	Name 	  string `json:"name" valid:"length(3|30)"`
	Description *string `json:"description" valid:"length(3|200)"`
	Created     time.Time `json:"created"`
}

// Valid - Bool para saber si la validación es correcta
func (okupa Okupa) Valid() bool {
	ok, err := govalidator.ValidateStruct(okupa)
	if err != nil {
		return false
	}
	return ok
}

//GetOkupas - Lista todas las okupas, las más recientes primero
func GetOkupas(db *gorm.DB) []*Okupa {
	 okupas := []*Okupa{}
	db.Order("created asc").Find(&okupas)
	return okupas
}


// GetOkupa - Devuelve la okupa con el ID
func GetOkupa(id int, db *gorm.DB) *Okupa {
	okupa := new(Okupa)
	db.Find(okupa, id)
	if okupa.OkupaID == id {
		return okupa
	}
	return nil
}
func CreateOkupa(okupa *Okupa, db *gorm.DB) {
	db.Create(okupa)
}
func DeleteOkupa(del *Okupa, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Delete(del)
}
func EditOkupa(editOkupa *Okupa, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Save(editOkupa)
}
