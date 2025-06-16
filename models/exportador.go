package models

import (
	"encoding/csv"
	"encoding/json"
	"os"
)

type Exportador interface {
	Exportar(libros []*Libro, archivo string) error
}

// ExportadorJSON exporta libros a JSON
type ExportadorJSON struct{}

func (e ExportadorJSON) Exportar(libros []*Libro, archivo string) error {
	data, err := json.MarshalIndent(libros, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(archivo, data, 0644)
}

// ExportadorCSV exporta libros a CSV
type ExportadorCSV struct{}

func (e ExportadorCSV) Exportar(libros []*Libro, archivo string) error {
	f, err := os.Create(archivo)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"ID", "Título", "Autor", "Género", "Archivo"})

	for _, libro := range libros {
		writer.Write([]string{
			toStr(libro.ID),
			libro.Titulo,
			libro.Autor,
			libro.Genero,
			libro.Archivo,
		})
	}
	return nil
}
