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

// GetRoutesOwners contiene las rutas de los "owner"
func GetRoutesOwners(r *mux.Router) {
	r.HandleFunc("/owners", createOwner).Methods(http.MethodPost)
	r.HandleFunc("/owners", getAllOwners).Methods(http.MethodGet)
	r.HandleFunc("/owners/{id:[0-9]+}", getOwner).Methods(http.MethodGet)
	r.HandleFunc("/owners/{id:[0-9]+}", updateOwner).Methods(http.MethodPut)
	r.HandleFunc("/owners/{id:[0-9]+}", deleteOwner).Methods(http.MethodDelete)
}
// updateOwner - Para modificar un owner. El id y la fecha de creación no se pueden modificar.
func updateOwner(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			owner := models.GetOwner(id,db)
			if owner != nil {
				jsonBytes, err := ioutil.ReadAll(r.Body)
				if err == nil {
					editOwner := new(models.Owner)
					err := json.Unmarshal(jsonBytes, editOwner)
					//Comprobamos que no haya error en el Unmarshall y la validación de los datos con govalidator
					if err == nil && editOwner.Valid(){
						//Solo administradores
						errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
						if errAuth == nil {
							owner.Name = editOwner.Name
							owner.Description = editOwner.Description
							owner.TypeOwner = editOwner.TypeOwner
							models.EditOwner(owner, db)
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
// deleteOwner - Para borrar una asociacion owner. Solo Admins
func deleteOwner(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			owner := models.GetOwner(id,db)
			if owner != nil {
				//Solo administradores
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
				if errAuth == nil {
					models.DeleteOwner(owner, db)
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


//Ver una asociacion owner. Los users registrados pueden verlo
func getOwner(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { //  *models.Owner
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			owner := models.GetOwner(id, db)
			if owner != nil {
				jsonTask, err := json.Marshal(owner)
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

// getAllOwners - Listar todos los owners. 
func getAllOwners(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.Owner
		db, _ := data.ConnectDB()
		defer db.Close()
		jsonTasks, err := json.Marshal(models.GetOwners(db))
		if err == nil {
			// Tiene permisos el owner?
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

// Crear owners. Solo admins
func createOwner(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			owner := new(models.Owner)
			err := json.Unmarshal(jsonBytes, owner)
			if err == nil && owner.Valid(){
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
				if errAuth == nil {
					db, _ := data.ConnectDB()
					defer db.Close()
					if err == nil {
						models.CreateOwner(owner,db)
						//Devolveremos por json la asociacion owner creada
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