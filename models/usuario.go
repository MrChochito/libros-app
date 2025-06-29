// models/usuario.go
package models

type Usuario struct {
	ID       int
	Nombre   string
	Correo   string
	Password string
	Avatar   string
	LibrosSubidos []Libro
	LibrosPrestados []Libro
}
