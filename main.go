package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/urfave/negroni"

	"github.com/yeadan/okupa/api/data"
	"github.com/yeadan/okupa/api/middlewares"
	"github.com/yeadan/okupa/api/routes"
)

func main() {

	// Nuevo mux.router y añadimos rutas
	router := mux.NewRouter().StrictSlash(true)
	routes.GetRoutesUsers(router)
	routes.GetRoutesOkupas(router)
	routes.GetRoutesOwners(router)
	routes.GetRoutesProperties(router)

	// Para habilitar CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:8080","https://apiv1.geoapi.es"})

	router.Use(middlewares.AuthUser)
	
	// Creamos Negroni y registamos middlewares
	middle := negroni.Classic()
	middle.UseHandler(router)

	//Iniciamos base de datos y .env
	data.InitDB()

	http.ListenAndServe(":4444", handlers.CORS(headers, methods, origins)(middle))
}
