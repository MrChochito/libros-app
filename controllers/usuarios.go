package controllers

import (
	"html/template"
	"net/http"
	"sistema-libros/models"

	"github.com/gorilla/sessions"
)

var (
	repoUsuarios models.RepositorioUsuario
	store        = sessions.NewCookieStore([]byte("clave-super-secreta"))
)

// SetRepositorioUsuarios asigna el repositorio para usuarios.
func SetRepositorioUsuarios(r models.RepositorioUsuario) {
	repoUsuarios = r
}

// MostrarFormularioRegistro muestra la página de registro.
func MostrarFormularioRegistro(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/registro.html"))
	tmpl.Execute(w, nil)
}

// RegistrarUsuario procesa la creación de un nuevo usuario.
func RegistrarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		usuario, err := models.NuevoUsuario(
			r.FormValue("nombre"),
			r.FormValue("email"),
			r.FormValue("password"),
			"usuario", // rol por defecto
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := repoUsuarios.Agregar(usuario); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// MostrarFormularioLogin muestra el formulario de login.
func MostrarFormularioLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

// LoginUsuario procesa la autenticación.
func LoginUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		usuario, err := repoUsuarios.ObtenerPorEmail(email)
		if err != nil || !usuario.VerificarPassword(password) {
			http.Error(w, "email o contraseña incorrectos", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "session")
		session.Values["usuario"] = usuario.Email()
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Logout cierra la sesión del usuario.
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "usuario")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
