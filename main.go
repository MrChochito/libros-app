package main

import (
	"net/http"
	"sistema-libros/controllers"
	"sistema-libros/models"
	"sistema-libros/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	repo := models.NewRepositorioMemoria()
	controllers.SetRepositorio(repo)

	routes.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
