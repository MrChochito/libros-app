package main

import (
	"log"
	"net/http"
	"strings"

	"libros-app/database"
	"libros-app/routes"
)

func split(s, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

func main() {
	// Conexión a la base de datos
	database.Conectar()

	// Servir archivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Cargar rutas
	routes.LoadRoutes()

	// Mensaje de inicio
	log.Println("✅ Servidor iniciado en http://localhost:8080")

	// Iniciar servidor web
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("❌ Error iniciando servidor:", err)
	}
}
 