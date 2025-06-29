package main

import (
	"encoding/json"
	"fmt"
	"log"
	"libros-app/models"
)

func main() {
	// Crear un libro de ejemplo
	libro := models.Libro{
		ID: 1,
		Titulo: "El Principito",
		Autor: "Antoine de Saint-Exupéry",
		Imagen: "el-principito.jpg",
		Resumen: "Un cuento para niños y adultos.",
		Etiquetas: "clásico,aventura",
		Disponible: true,
		Vistas: 100,
		VecesPrestado: 10,
		UsuarioID: 2,
		Slug: "el-principito",
		DuracionPrestamoDias: 7,
		PDF: "el-principito.pdf",
	}

	// Serialización a JSON
	jsonBytes, err := json.Marshal(libro)
	if err != nil {
		log.Fatal("Error serializando a JSON:", err)
	}
	fmt.Println("Libro serializado a JSON:", string(jsonBytes))

	// Deserialización desde JSON
	var libro2 models.Libro
	err = json.Unmarshal(jsonBytes, &libro2)
	if err != nil {
		log.Fatal("Error deserializando desde JSON:", err)
	}
	fmt.Printf("Libro deserializado: %+v\n", libro2)
}
