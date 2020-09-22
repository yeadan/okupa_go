package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// Address - Struct de una direccion
type Address struct {
	Calle  			string `json:"calle" valid:"length(2|50)"` 
	Numero			int `json:"numero"`
	Piso			int `json:"piso"`
	Puerta			string `json:"puerta"` 
	CodigoPostal	string `json:"codigo_postal"`
	Nucleo			string `json:"nucleo" valid:"length(2|50)"` 
	Poblacion		string `json:"poblacion" valid:"length(2|50)"` 
	Municipio		string `json:"municipio" valid:"length(2|50)"` 
	Provincia		string `json:"provincia" valid:"length(2|50)"` 
	Comunidad		string `json:"comunidad" valid:"length(2|50)"`  
}
// Property - Struct con los datos de un property.
type Property struct {
	PropertyID  	int `gorm:"primary_key;not null" json:"property_id"`
	OwnerID   		int `gorm:"not null;type:int REFERENCES owners(owner_id) ON DELETE SET NULL" json:"owner_id"`
	OkupaID		    int `gorm:"not null;type:int REFERENCES okupas(okupa_id) ON DELETE SET NULL" json:"okupa_id"`
	UserID    		int `gorm:"not null;type:int REFERENCES users(user_id) ON DELETE SET NULL" json:"user_id"` 		
	Type	  		string `json:"type" valid:"length(3|15)"` 
	Description 	*string `json:"description" valid:"length(3|200)"`
	Created     	time.Time `json:"created"`
	Address
	//Añadidos para mejorar los listados
	User        	*User `json:"user"` //no son FK por los errores del gorm con autoincrementos
	Owner			*Owner `json:"owner"`
	Okupa			*Okupa `json:"okupa"`

}

// Valid - Bool para saber si la validación es correcta
func (property Property) Valid() bool {
	ok, err := govalidator.ValidateStruct(property)
	if err != nil {
		return false
	}
	return ok
}

//GetProperties - Lista todas las properties, las más recientes primero
func GetProperties(db *gorm.DB) []*Property {
	 properties := []*Property{}
	db.Order("created desc").Find(&properties)
	return properties
}


// GetProperty - Devuelve el property con el ID
func GetProperty(id int, db *gorm.DB) *Property {
	property := new(Property)
	db.Find(property, id)
	if property.PropertyID == id {
		return property
	}
	return nil
}
// GetPropertyOkupa - Busca todos los properties okupados por "okupa"
func GetPropertyOkupa(okupa int, db *gorm.DB) []*Property {
	properties := []*Property{}
	db.Where("Property.okupa_id = ?", okupa).Find(&properties)
	return properties
}
// GetPropertyUser - Busca todos los properties okupados por "user"
func GetPropertyUser(user int, db *gorm.DB) []*Property {
	properties := []*Property{}
	db.Where("Property.user_id = ?", user).Find(&properties)
	return properties
}
// GetPropertyOwner - Busca todos los properties okupados por "owner"
func GetPropertyOwner(owner int, db *gorm.DB) []*Property {
	properties := []*Property{}
	db.Where("Property.owner_id = ?", owner).Find(&properties)
	return properties
}
// GetPropertyType - Busca todos los properties okupados por "type"
func GetPropertyType(id int, db *gorm.DB) []*Property {
	properties := []*Property{}
	db.Where("Property.type = ?", id).Find(&properties)
	return properties
}


func CreateProperty(property *Property, db *gorm.DB) {
	db.Create(property)
}
func DeleteProperty(del *Property, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Delete(del)
}
func EditProperty(editProperty *Property, db *gorm.DB) {
	db.BlockGlobalUpdate(true)
	db.Save(editProperty)
}
