package models

import (
	"errors"
	"strings"
	"golang.org/x/crypto/bcrypt"
)

// Usuario representa un usuario del sistema.
type Usuario struct {
	nombre   string
	email    string
	password string // almacenaremos el hash de la contraseña
	rol      string // Ejemplo: "admin" o "usuario"
}

// NuevoUsuario crea un nuevo usuario con validación y encriptación de contraseña.
func NuevoUsuario(nombre, email, password, rol string) (*Usuario, error) {
	if strings.TrimSpace(nombre) == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("el email no puede estar vacío")
	}
	if len(password) < 6 {
		return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
	}
	if rol != "admin" && rol != "usuario" {
		return nil, errors.New("rol inválido")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Usuario{
		nombre:   nombre,
		email:    email,
		password: string(hashed),
		rol:      rol,
	}, nil
}

// Getters
func (u *Usuario) Nombre() string {
	return u.nombre
}

func (u *Usuario) Email() string {
	return u.email
}

func (u *Usuario) Rol() string {
	return u.rol
}

// VerificarPassword comprueba si la contraseña proporcionada es correcta.
func (u *Usuario) VerificarPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

// CambiarPassword permite actualizar la contraseña (encriptada).
func (u *Usuario) CambiarPassword(nuevaPassword string) error {
	if len(nuevaPassword) < 6 {
		return errors.New("la contraseña debe tener al menos 6 caracteres")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(nuevaPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = string(hashed)
	return nil
}

// RepositorioUsuario define la interfaz para manipular usuarios.
type RepositorioUsuario interface {
	Agregar(usuario *Usuario) error
	ObtenerPorEmail(email string) (*Usuario, error)
	ObtenerTodos() []*Usuario
}
