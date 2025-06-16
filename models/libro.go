package models

import (
	"errors"
	"strings"
)

type Libro struct {
	ID      int
	Titulo  string
	Autor   string
	Genero  string
	Archivo string
}

func NuevoLibro(id int, titulo, autor, genero, archivo string) (*Libro, error) {
	if strings.TrimSpace(titulo) == "" {
		return nil, errors.New("el título no puede estar vacío")
	}
	if strings.TrimSpace(archivo) == "" {
		return nil, errors.New("la ruta de archivo no puede estar vacía")
	}
	return &Libro{
		ID:      id,
		Titulo:  titulo,
		Autor:   autor,
		Genero:  genero,
		Archivo: archivo,
	}, nil
}
