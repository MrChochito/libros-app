package controllers

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"libros-app/database"
	"libros-app/models"
	"libros-app/utils"
)

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	// Si ya hay sesión, redirige a home
	if cookie, err := r.Cookie("session"); err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Mostrar mensajes de error si existen
	tmpl := template.New("login.html").Funcs(template.FuncMap{"split": utils.Split})
	tmpl = template.Must(tmpl.ParseFiles("views/templates/login.html"))

	// Obtener mensajes de error de la query
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")
	data := map[string]interface{}{
		"Error":   errorMsg,
		"Success": successMsg,
	}
	tmpl.Execute(w, data)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	correo := r.FormValue("correo")
	password := r.FormValue("password")

	// Obtener hash de BD
	var hash string
	err := database.DB.QueryRow(
		"SELECT password FROM usuarios WHERE correo = @p1",
		correo,
	).Scan(&hash)
	if err != nil {
		// usuario no existe o error DB
		http.Redirect(w, r, "/login?error=Correo+o+contrase%C3%B1a+incorrectos", http.StatusSeeOther)
		return
	}

	// Comparar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=Correo+o+contrase%C3%B1a+incorrectos", http.StatusSeeOther)
		return
	}

	// Crear sesión
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: correo,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegistrarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	nombre := r.FormValue("nombre")
	correo := r.FormValue("correo")
	password := r.FormValue("password")

	// Validar duplicado
	var count int
	err := database.DB.QueryRow(
		"SELECT COUNT(*) FROM usuarios WHERE correo = @p1",
		correo,
	).Scan(&count)
	if err != nil {
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Redirect(w, r, "/login?error=Correo+ya+registrado", http.StatusSeeOther)
		return
	}

	// Hashear password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	// Insertar usuario con avatar por defecto
	avatarDefault := "/static/img/avatar.png"
	_, err = database.DB.Exec(
		"INSERT INTO usuarios (nombre, correo, password, avatar) VALUES (@p1,@p2,@p3,@p4)",
		nombre, correo, string(hash), avatarDefault,
	)
	if err != nil {
		http.Error(w, "Error al registrar", http.StatusInternalServerError)
		return
	}

	// Iniciar sesión automático
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: correo,
		Path:  "/",
	})
	http.Redirect(w, r, "/?success=Cuenta+creada+correctamente", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	// Validar sesión
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	correo := cookie.Value

	// Obtener datos de usuario
	var user models.Usuario
	err = database.DB.QueryRow(
		"SELECT id, nombre, correo, avatar FROM usuarios WHERE correo = @p1",
		correo,
	).Scan(&user.ID, &user.Nombre, &user.Correo, &user.Avatar)
	if err != nil {
		http.Error(w, "Error al cargar perfil", http.StatusInternalServerError)
		return
	}

	// Marcar automáticamente como devueltos los préstamos vencidos
	_ = models.MarcarPrestamosVencidosComoDevueltos()

	// Obtener libros subidos por el usuario
	librosSubidos, err := models.GetLibrosSubidos(user.ID)
	if err != nil {
		librosSubidos = []models.Libro{}
	}

	// Obtener libros prestados por el usuario
	librosPrestados, err := models.GetLibrosPrestadosPorUsuario(user.ID)
	if err != nil {
		librosPrestados = []models.LibroPrestado{}
	}

	// Renderizar plantilla
	tmpl := template.Must(template.ParseFiles(
		"views/templates/perfil.html",
		"views/templates/header.html",
	))
	tmpl.Execute(w, map[string]interface{}{
		"User":            user,
		"LibrosSubidos":   librosSubidos,
		"LibrosPrestados": librosPrestados,
		"Active":          "perfil",
		"ErrorExpirado":   r.URL.Query().Get("error") == "expirado",
	})
}
