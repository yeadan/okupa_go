package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// Owner - Struct con los datos de un owner.
type Owner struct {
	OwnerID   int `gorm:"primary_key;not null" json:"owner_id"`
	Name 	  string `json:"name" valid:"length(3|30)"`
	TypeOwner string `json:"type_owner" valid:"length(3|15)"` 
	Description *string `json:"description" valid:"length(3|200)"`
	Created     time.Time `json:"created"`
}

// Valid - Bool para saber si la validación es correcta
func (owner Owner) Valid() bool {
	ok, err := govalidator.ValidateStruct(owner)
	if err != nil {
		return false
	}
	return ok
}

//GetOwners - Lista todas las owners, las más recientes primero
func GetOwners(db *gorm.DB) []*Owner {
	 owners := []*Owner{}
	db.Order("created asc").Find(&owners)
	return owners
}


// GetOwner - Devuelve el owner con el ID
func GetOwner(id int, db *gorm.DB) *Owner {
	owner := new(Owner)
	db.Find(owner, id)
	if owner.OwnerID == id {
		return owner
	}
	return nil
}
func CreateOwner(owner *Owner, db *gorm.DB) {
	db.Create(owner)
}
func DeleteOwner(del *Owner, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Delete(del)
}
func EditOwner(editOwner *Owner, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Save(editOwner)
}
