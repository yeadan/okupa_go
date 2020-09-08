package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yeadan/okupa/api/data"
	"github.com/yeadan/okupa/api/middlewares"
	"github.com/yeadan/okupa/api/models"
	"github.com/yeadan/okupa/lib"
)

// GetRoutesOkupas contiene las rutas de las asociaciones "okupas"
func GetRoutesOkupas(r *mux.Router) {
	r.HandleFunc("/okupas", createOkupa).Methods(http.MethodPost)
	r.HandleFunc("/okupas", getAllOkupas).Methods(http.MethodGet)
	r.HandleFunc("/okupas/{id:[0-9]+}", getOkupa).Methods(http.MethodGet)
	r.HandleFunc("/okupas/{id:[0-9]+}", updateOkupa).Methods(http.MethodPut)
	r.HandleFunc("/okupas/{id:[0-9]+}", deleteOkupa).Methods(http.MethodDelete)
	//Rutas de "UserOkupa" - La relación entre users y okupas
	r.HandleFunc("/okupas/{id:[0-9]+}/{usr:[0-9]+}", createUserOkupa).Methods(http.MethodPost)
	r.HandleFunc("/okupas/{id:[0-9]+}/{usr:[0-9]+}", deleteUserOkupa).Methods(http.MethodDelete)
	r.HandleFunc("/okupas/users/{id:[0-9]+}", getUserOkupa).Methods(http.MethodGet)
}
// getUserOkupa - Listar usuarios de una asociación okupa
func getUserOkupa(w http.ResponseWriter, r *http.Request) {
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
					okupaJSON := []*models.User{}
					getUser := models.GetUserOkupaOkupa(id,db)
					for _,proba := range getUser {
						single := models.GetUser(proba.UserID,db)
						okupaJSON = append (okupaJSON, single)
					}
					jsonFinal, _ := json.Marshal(okupaJSON)
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
// createUserOkupa - Añadir un usuario a una asociación okupa
func createUserOkupa(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if okpStr, ok := mux.Vars(r)["id"]; ok {
			if usrStr, ok := mux.Vars(r)["usr"]; ok {
				okp, _ := strconv.Atoi(okpStr)
				usr, _ := strconv.Atoi(usrStr)
				db, _ := data.ConnectDB()
				defer db.Close()
				okupa := models.GetOkupa(okp,db)
				usuario := models.GetUser(usr,db)
				if okupa != nil && usuario != nil {
					errAuth := lib.UserAllowed(userValid.(*models.User), &usr, lib.GetString("admin"), w)
					if errAuth == nil { 
						result := models.ExistUserOkupa(usr,okp,db)
						//miramos si ya existe esa relación entre usuario y asociación okupa
						if result == nil {
							newRel := new(models.UserOkupa)
							newRel.UserID = usuario.UserID
							newRel.OkupaID = okupa.OkupaID
							newRel.Created = time.Now()
							models.CreateUserOkupa(newRel, db)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusCreated)
							jsonTask, _ := json.Marshal(newRel)
							w.Write(jsonTask)
						} else {
							w.WriteHeader(http.StatusBadRequest)
							w.Write([]byte("Error. Ya existe esa relación."))
						}
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// deleteUserOkupa - Borrar la relación entre usuario y asociación okupa
func deleteUserOkupa(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if okpStr, ok := mux.Vars(r)["id"]; ok {
			if usrStr, ok := mux.Vars(r)["usr"]; ok {
				okp, _ := strconv.Atoi(okpStr)
				usr, _ := strconv.Atoi(usrStr)
				db, _ := data.ConnectDB()
				defer db.Close()
				okupa := models.GetOkupa(okp,db)
				usuario := models.GetUser(usr,db)
				if okupa != nil && usuario != nil {
					errAuth := lib.UserAllowed(userValid.(*models.User), &usr, lib.GetString("admin"), w)
					if errAuth == nil { 
						result := models.ExistUserOkupa(usr,okp,db)
						//miramos si ya existe esa relación entre usuario y asociación okupa
						if result != nil {
							models.DeleteUserOkupa(result,db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusNotFound)
							w.Write([]byte("Error. No existe esa relación."))
						}
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// updateOkupa - Para modificar un okupa. El id y la fecha de creación no se pueden modificar.
func updateOkupa(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			okupa := models.GetOkupa(id,db)
			if okupa != nil {
				jsonBytes, err := ioutil.ReadAll(r.Body)
				if err == nil {
					editOkupa := new(models.Okupa)
					err := json.Unmarshal(jsonBytes, editOkupa)
					//Comprobamos que no haya error en el Unmarshall y la validación de los datos con govalidator
					if err == nil && editOkupa.Valid(){
						//Solo administradores
						errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
						if errAuth == nil {
							okupa.Name = editOkupa.Name
							okupa.Description = editOkupa.Description
							models.EditOkupa(okupa, db)
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
// deleteOkupa - Para borrar una asociacion okupa. Solo Admins
func deleteOkupa(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			okupa := models.GetOkupa(id,db)
			if okupa != nil {
				//Solo administradores
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
				if errAuth == nil {
					models.DeleteOkupa(okupa, db)
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


//Ver una asociacion okupa. Los users registrados pueden verlo
func getOkupa(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { //  *models.Okupa
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			okupa := models.GetOkupa(id, db)
			if okupa != nil {
				jsonTask, err := json.Marshal(okupa)
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

// getAllOkupas - Listar todos los okupas. 
func getAllOkupas(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.Okupa
		db, _ := data.ConnectDB()
		defer db.Close()
		jsonTasks, err := json.Marshal(models.GetOkupas(db))
		if err == nil {
			// Tiene permisos el okupa?
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

// Crear asociaciones okupas. Solo admins
func createOkupa(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			okupa := new(models.Okupa)
			err := json.Unmarshal(jsonBytes, okupa)
			if err == nil && okupa.Valid(){
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
				if errAuth == nil {
					db, _ := data.ConnectDB()
					defer db.Close()
					if err == nil {
						models.CreateOkupa(okupa,db)
						//Devolveremos por json la asociacion okupa creada
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