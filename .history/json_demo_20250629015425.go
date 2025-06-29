package main

import (
	"encoding/json"
	"fmt"
	"libros-app/models"
	"log"
	"os"
)

func main() {
	// Crear un libro de ejemplo
	libro := models.Libro{
		ID:                   1,
		Titulo:               "El Principito",
		Autor:                "Antoine de Saint-Exupéry",
		Imagen:               "el-principito.jpg",
		Resumen:              "Un cuento para niños y adultos.",
		Etiquetas:            "clásico,aventura",
		Disponible:           true,
		Vistas:               100,
		VecesPrestado:        10,
		UsuarioID:            2,
		Slug:                 "el-principito",
		DuracionPrestamoDias: 7,
		PDF:                  "el-principito.pdf",
	}

	// Serialización a JSON y guardar en archivo
	jsonBytes, err := json.MarshalIndent(libro, "", "  ")
	if err != nil {
		log.Fatal("Error serializando a JSON:", err)
	}
	fmt.Println("Libro serializado a JSON:")
	fmt.Println(string(jsonBytes))

	file, err := os.Create("libro.json")
	if err != nil {
		log.Fatal("No se pudo crear el archivo:", err)
	}
	defer file.Close()
	file.Write(jsonBytes)

	// Leer desde archivo y deserializar
	fileRead, err := os.ReadFile("libro.json")
	if err != nil {
		log.Fatal("No se pudo leer el archivo:", err)
	}
	var libro2 models.Libro
	err = json.Unmarshal(fileRead, &libro2)
	if err != nil {
		log.Fatal("Error deserializando desde JSON:", err)
	}
	fmt.Printf("Libro deserializado: %+v\n", libro2)
}
