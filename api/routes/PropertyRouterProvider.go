package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yeadan/okupa/api/data"
	"github.com/yeadan/okupa/api/middlewares"
	"github.com/yeadan/okupa/api/models"
	"github.com/yeadan/okupa/lib"
)

// GetRoutesProperties contiene las rutas de los "property"
func GetRoutesProperties(r *mux.Router) {
	r.HandleFunc("/properties", createProperty).Methods(http.MethodPost)
	r.HandleFunc("/properties", getAllPropertys).Methods(http.MethodGet)
	r.HandleFunc("/properties/{id:[0-9]+}", getProperty).Methods(http.MethodGet)
	r.HandleFunc("/properties/{id:[0-9]+}", updateProperty).Methods(http.MethodPut)
	r.HandleFunc("/properties/{id:[0-9]+}", deleteProperty).Methods(http.MethodDelete)
	r.HandleFunc("/properties/users/{id:[0-9]+}", getUserProperty).Methods(http.MethodGet)
	r.HandleFunc("/properties/owners/{id:[0-9]+}", getOwnerProperty).Methods(http.MethodGet)
	r.HandleFunc("/properties/okupas/{id:[0-9]+}", getOkupaProperty).Methods(http.MethodGet)
	r.HandleFunc("/properties/types/{id:[0-9]+}", getTypeProperty).Methods(http.MethodGet)
}
// getTypeProperty - Devuelve las propiedades por "type"
func getTypeProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			if id > 0 {
				//Comprobar permisos y roles
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					propertyJSON := models.GetPropertyType(id,db)
					jsonFinal, _ := json.Marshal(propertyJSON)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(jsonFinal)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// getOkupaProperty - Devuelve las propiedades por "okupa"
func getOkupaProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			okupa := models.GetOkupa(id, db)
			if okupa != nil {
				//Comprobar permisos y roles
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					propertyJSON := models.GetPropertyOkupa(id,db)
					jsonFinal, _ := json.Marshal(propertyJSON)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(jsonFinal)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// getOwnerProperty - Devuelve las propiedades por "owner"
func getOwnerProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			owner := models.GetOwner(id, db)
			if owner != nil {
				//Comprobar permisos y roles
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					propertyJSON := models.GetPropertyOwner(id,db)
					jsonFinal, _ := json.Marshal(propertyJSON)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(jsonFinal)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// getUserProperty - Devuelve las propiedades por "user"
func getUserProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			user := models.GetOwner(id, db)
			if user != nil {
				//Comprobar permisos y roles
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					propertyJSON := models.GetPropertyUser(id,db)
					jsonFinal, _ := json.Marshal(propertyJSON)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(jsonFinal)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// updateProperty - Para modificar un property. El id y la fecha de creación no se pueden modificar.
func updateProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			property := models.GetProperty(id,db)
			if property != nil {
				jsonBytes, err := ioutil.ReadAll(r.Body)
				if err == nil {
					editProperty := new(models.Property)
					err := json.Unmarshal(jsonBytes, editProperty)
					//Comprobamos que no haya error en el Unmarshall y la validación de los datos con govalidator
					if err == nil && editProperty.Valid(){
						//Solo administradores
						errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
						if errAuth == nil {
							
							property.Description = editProperty.Description
							property.Type = editProperty.Type
							models.EditProperty(property, db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusForbidden)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				} else {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Error leyendo el body"))
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
// deleteProperty - Para borrar una asociacion property. Solo Admins
func deleteProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			property := models.GetProperty(id,db)
			if property != nil {
				//Solo administradores
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
				if errAuth == nil {
					models.DeleteProperty(property, db)
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusForbidden)
				} 
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func getProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			property := models.GetProperty(id, db)
			if property != nil {
				jsonTask, err := json.Marshal(property)
				if err == nil {
					//Comprobar permisos y roles
					errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
					if errAuth == nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(jsonTask)
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// getAllProperties - Listar todos los properties. 
func getAllPropertys(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.Property
		db, _ := data.ConnectDB()
		defer db.Close()
		jsonTasks, err := json.Marshal(models.GetProperties(db))
		if err == nil {
			// Tiene permisos el property?
			errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
			if errAuth == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonTasks)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// Crear asociaciones properties. Solo admins
func createProperty(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			property := new(models.Property)
			err := json.Unmarshal(jsonBytes, property)
			if err == nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
				if errAuth == nil {
					db, _ := data.ConnectDB()
					defer db.Close()
					if err == nil {
						models.CreateProperty(property,db)
						//Devolveremos por json la asociacion property creada
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)

					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error en el body"))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error leyendo el body"))
		}
	}
}