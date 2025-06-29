package routes

import (
	"libros-app/controllers"
	"net/http"
)

// LoadRoutes configura las rutas principales de la aplicación
func LoadRoutes() {
	// Rutas públicas
	http.HandleFunc("/login", controllers.ShowLogin)
	http.HandleFunc("/auth", controllers.Authenticate)
	http.HandleFunc("/register", controllers.RegistrarUsuario)

	// Rutas privadas (requieren sesión)
	http.HandleFunc("/", controllers.HomeLogged) // Home protegido
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/upload", controllers.UploadHandler)
	http.HandleFunc("/perfil", controllers.Profile)
	http.HandleFunc("/libro/", controllers.LibroMultipropositoHandler) // Maneja detalle y edición
	http.HandleFunc("/descargar/", controllers.DescargarPDFPrestamoHandler)
	http.HandleFunc("/prestamo/", controllers.EliminarPrestamoHandler)

	// Rutas de la API para las gráficas
	http.HandleFunc("/api/graficas/mas-prestados", controllers.GraficaMasPrestados)
	http.HandleFunc("/api/graficas/categorias", controllers.GraficaCategorias)
	http.HandleFunc("/api/graficas/prestamos-semana", controllers.GraficaPrestamosSemana)
}

