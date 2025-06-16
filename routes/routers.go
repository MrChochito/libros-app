package routes

import (
	"sistema-libros/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	// Rutas libros
	r.HandleFunc("/libros", controllers.MostrarLibros).Methods("GET")
	r.HandleFunc("/libros/nuevo", controllers.FormAgregarLibro).Methods("GET")
	r.HandleFunc("/libros/nuevo", controllers.AgregarLibro).Methods("POST")

	// Rutas usuarios
	r.HandleFunc("/registro", controllers.MostrarFormularioRegistro).Methods("GET")
	r.HandleFunc("/registro", controllers.RegistrarUsuario).Methods("POST")
	r.HandleFunc("/login", controllers.MostrarFormularioLogin).Methods("GET")
	r.HandleFunc("/login", controllers.LoginUsuario).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")
	r.HandleFunc("/buscar", controllers.BuscarLibro).Methods("POST")
	r.HandleFunc("/exportar", controllers.ExportarLibros).Methods("GET")

}
