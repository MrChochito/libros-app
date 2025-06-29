package controllers

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"libros-app/database"
	"libros-app/models"
	"libros-app/utils"
)

// Home muestra la página principal con los libros más recientes y populares.
func Home(w http.ResponseWriter, r *http.Request) {
	// Validar sesión del usuario
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtener libros recientes y populares
	recientes, err := models.ObtenerLibrosRecientes()
	if err != nil {
		http.Error(w, "Error obteniendo libros recientes", http.StatusInternalServerError)
		return
	}
	masVistos, err := models.ObtenerLibrosMasVistos()
	if err != nil {
		http.Error(w, "Error obteniendo libros populares", http.StatusInternalServerError)
		return
	}

	// Renderizar la plantilla principal
	tmpl := template.New("home.html").Funcs(template.FuncMap{"split": utils.Split})
	tmpl = template.Must(tmpl.ParseFiles(
		"views/templates/home.html",
		"views/templates/header.html",
	))
	err = tmpl.Execute(w, map[string]interface{}{
		"Recientes": recientes,
		"Populares": masVistos,
		"Active":    "home",
	})
	if err != nil {
		http.Error(w, "Error mostrando página", http.StatusInternalServerError)
	}
}

// HomeLogged muestra la página principal para usuarios logueados.
func HomeLogged(w http.ResponseWriter, r *http.Request) {
	// Igual que Home, pero para usuarios autenticados
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	recientes, err := models.ObtenerLibrosRecientes()
	if err != nil {
		http.Error(w, "Error obteniendo libros recientes", http.StatusInternalServerError)
		return
	}
	masVistos, err := models.ObtenerLibrosMasVistos()
	if err != nil {
		http.Error(w, "Error obteniendo libros populares", http.StatusInternalServerError)
		return
	}

	tmpl := template.New("home.html").Funcs(template.FuncMap{"split": utils.Split})
	tmpl = template.Must(tmpl.ParseFiles(
		"views/templates/home.html",
		"views/templates/header.html",
	))
	err = tmpl.Execute(w, map[string]interface{}{
		"Recientes": recientes,
		"Populares": masVistos,
		"Active":    "home",
	})
	if err != nil {
		http.Error(w, "Error mostrando página", http.StatusInternalServerError)
	}
}

// UploadHandler permite subir un nuevo libro con portada y PDF.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Validar sesión del usuario
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		// Mostrar formulario de subida
		tmpl := template.Must(template.ParseFiles(
			"views/templates/upload.html",
			"views/templates/header.html",
		))
		tmpl.Execute(w, map[string]interface{}{
			"Active": "upload",
		})
		return
	}

	// Procesar POST: obtener datos del formulario
	title := r.FormValue("titulo")
	author := r.FormValue("autor")
	resume := r.FormValue("resumen")
	tags := r.FormValue("etiquetas")

	// Generar slug único para el libro
	slug := utils.GenerateUniqueSlug(title)

	// Guardar imagen de portada con nombre basado en el slug
	file, handler, err := r.FormFile("imagen")
	if err != nil {
		http.Error(w, "Error imagen", http.StatusBadRequest)
		return
	}
	defer file.Close()
	uploadDir := "static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	imgExt := filepath.Ext(handler.Filename)
	imgFilename := slug + imgExt
	imgPath := filepath.Join(uploadDir, imgFilename)
	dst, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Error guardando imagen", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)
	imgWebPath := "/static/uploads/" + imgFilename

	// Guardar PDF con nombre basado en el slug
	pdfFile, pdfHandler, err := r.FormFile("pdf")
	var pdfPath string
	if err != nil || pdfHandler == nil {
		http.Error(w, "El PDF es obligatorio", http.StatusBadRequest)
		return
	}
	pdfDir := "static/uploads/pdf"
	os.MkdirAll(pdfDir, os.ModePerm)
	pdfExt := filepath.Ext(pdfHandler.Filename)
	pdfFilename := slug + pdfExt
	pdfPathFS := filepath.Join(pdfDir, pdfFilename)
	pdfDst, err := os.Create(pdfPathFS)
	if err == nil {
		defer pdfDst.Close()
		io.Copy(pdfDst, pdfFile)
		pdfPath = "/static/uploads/pdf/" + pdfFilename
	} else {
		pdfPath = ""
	}
	pdfFile.Close()

	// Obtener ID usuario desde la sesión
	userID, err := utils.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "No se pudo obtener el usuario de la sesión", http.StatusUnauthorized)
		return
	}

	// Insertar libro en la base de datos
	var libroID int64
	duracion := 21 // Duración del préstamo por defecto
	duracionStr := r.FormValue("duracion_prestamo_dias")
	if duracionStr != "" {
		duracionParsed, err := strconv.Atoi(duracionStr)
		if err == nil && duracionParsed > 0 {
			duracion = duracionParsed
		}
	}
	insertQuery := `INSERT INTO libros (
		titulo, autor, imagen, resumen, etiquetas, id_usuario, slug, pdf, duracion_prestamo_dias
	) OUTPUT INSERTED.id VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9)`
	err = database.DB.QueryRow(
		insertQuery,
		title, author, imgWebPath, resume, tags, userID, slug, pdfPath, duracion,
	).Scan(&libroID)
	if err != nil {
		http.Error(w, "Error guardando libro", http.StatusInternalServerError)
		return
	}

	// Mostrar mensaje de éxito
	tmpl := template.Must(template.ParseFiles(
		"views/templates/upload.html",
		"views/templates/header.html",
	))
	tmpl.Execute(w, map[string]interface{}{
		"Active":  "upload",
		"Success": true,
	})
}

// BookDetail muestra el detalle de un libro y controla el acceso al PDF.
func BookDetail(w http.ResponseWriter, r *http.Request) {
	// Extraer slug del libro desde la URL
	slug := r.URL.Path[len("/libro/"):]
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	// Buscar libro por slug
	libro, err := models.GetLibroBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	userID, err := utils.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "No se pudo obtener el usuario de la sesión", http.StatusUnauthorized)
		return
	}

	// Verificar si el usuario tiene préstamo activo para este libro
	tienePrestamoActivo := false
	if userID > 0 {
		tienePrestamoActivo, _ = models.UsuarioTienePrestamoActivo(libro.ID, userID)
	}

	// Permitir ver el PDF si es el dueño o tiene préstamo activo
	puedeVerPDF := tienePrestamoActivo || (userID == libro.UsuarioID)

	// Renderizar la plantilla de detalle
	tmpl := template.New("libro_detalle.html").Funcs(template.FuncMap{"split": utils.Split, "eq": func(a, b int) bool { return a == b }})
	tmpl = template.Must(tmpl.ParseFiles(
		"views/templates/libro_detalle.html",
		"views/templates/header.html",
	))
	tmpl.Execute(w, map[string]interface{}{
		"Book":           libro,
		"UserID":         userID,
		"Active":         "",
		"PrestamoActivo": tienePrestamoActivo,
		"PuedeVerPDF":    puedeVerPDF,
	})
}

// EditarLibroHandler permite al usuario dueño editar su libro.
func EditarLibroHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/libro/")
	if !strings.HasSuffix(slug, "/editar") {
		// No es la ruta de edición, delegar a BookDetail
		BookDetail(w, r)
		return
	}
	slug = strings.TrimSuffix(slug, "/editar")
	slug = strings.TrimSuffix(slug, "/")
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	libro, err := models.GetLibroBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	userID, err := utils.GetUserIDFromSession(r)
	if err != nil || libro.UsuarioID != userID {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles(
			"views/templates/editar_libro.html",
			"views/templates/header.html",
		))
		tmpl.Execute(w, map[string]interface{}{
			"Book":   libro,
			"Active": "",
		})
		return
	}
	// Procesar POST (solo campos editables)
	titulo := r.FormValue("titulo")
	autor := r.FormValue("autor")
	resumen := r.FormValue("resumen")
	etiquetas := r.FormValue("etiquetas")
	// Portada (opcional)
	file, handler, err := r.FormFile("imagen")
	imgWebPath := libro.Imagen
	if err == nil && handler != nil {
		uploadDir := "static/uploads"
		os.MkdirAll(uploadDir, os.ModePerm)
		path := filepath.Join(uploadDir, handler.Filename)
		dst, err := os.Create(path)
		if err == nil {
			defer dst.Close()
			io.Copy(dst, file)
			imgWebPath = "/static/uploads/" + strings.ReplaceAll(handler.Filename, " ", "-")
			imgWebPath = strings.ReplaceAll(imgWebPath, "\\", "/")
		}
		file.Close()
	}
	// PDF (opcional)
	pdfFile, pdfHandler, err := r.FormFile("pdf")
	pdfPath := libro.PDF // <-- Mantener el PDF anterior si no se sube uno nuevo
	if err == nil && pdfHandler != nil {
		pdfDir := "static/uploads/pdf"
		os.MkdirAll(pdfDir, os.ModePerm)
		pdfPathFS := filepath.Join(pdfDir, pdfHandler.Filename)
		pdfDst, err := os.Create(pdfPathFS)
		if err == nil {
			defer pdfDst.Close()
			io.Copy(pdfDst, pdfFile)
			pdfPath = "/static/uploads/pdf/" + strings.ReplaceAll(pdfHandler.Filename, " ", "-")
		}
		pdfFile.Close()
	}
	// Duración del préstamo (opcional, default 21 días)
	duracionStr := r.FormValue("duracion_prestamo_dias")
	duracion := 21
	if duracionStr != "" {
		duracionParsed, err := strconv.Atoi(duracionStr)
		if err == nil && duracionParsed > 0 {
			duracion = duracionParsed
		}
	}
	// Actualizar en BD
	_, err = database.DB.Exec(
		`UPDATE libros SET titulo=@p1, autor=@p2, resumen=@p3, etiquetas=@p4, imagen=@p5, pdf=@p6, duracion_prestamo_dias=@p8 WHERE slug=@p7`,
		titulo, autor, resumen, etiquetas, imgWebPath, pdfPath, slug, duracion,
	)
	if err != nil {
		tmpl := template.Must(template.ParseFiles(
			"views/templates/editar_libro.html",
			"views/templates/header.html",
		))
		tmpl.Execute(w, map[string]interface{}{
			"Book":   libro,
			"Error":  "Error actualizando libro",
			"Active": "",
		})
		return
	}
	// Recargar datos actualizados
	libro, _ = models.GetLibroBySlug(slug)
	tmpl := template.Must(template.ParseFiles(
		"views/templates/editar_libro.html",
		"views/templates/header.html",
	))
	tmpl.Execute(w, map[string]interface{}{
		"Book":    libro,
		"Success": true,
		"Active":  "",
	})
}

// LibroMultipropositoHandler decide si mostrar detalle o edición según la URL.
func LibroMultipropositoHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/editar") {
		EditarLibroHandler(w, r)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/prestar") && r.Method == http.MethodPost {
		// Extraer slug
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 3 {
			slug := parts[len(parts)-2]
			PedirPrestamoPorSlugHandler(w, r, slug)
			return
		}
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	BookDetail(w, r)
}

// PedirPrestamoPorSlugHandler maneja la solicitud de préstamo usando el slug.
func PedirPrestamoPorSlugHandler(w http.ResponseWriter, r *http.Request, slug string) {
	// Validar sesión
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	usuarioID, err := utils.ObtenerUsuarioIDDesdeSession(cookie.Value)
	if err != nil {
		http.Error(w, "No se pudo obtener el usuario", http.StatusUnauthorized)
		return
	}
	libro, err := models.GetLibroBySlug(slug)
	if err != nil {
		http.Error(w, "Libro no encontrado", http.StatusNotFound)
		return
	}
	err = models.PedirPrestamo(libro.ID, usuarioID)
	if err != nil {
		if err.Error() == "Ya tienes este libro prestado y no lo has devuelto" {
			http.Redirect(w, r, "/perfil#tab-borrowed", http.StatusSeeOther)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/perfil#tab-borrowed", http.StatusSeeOther)
}

// DescargarPDFPrestamoHandler permite descargar el PDF prestado si el préstamo está vigente.
func DescargarPDFPrestamoHandler(w http.ResponseWriter, r *http.Request) {
	// Validar sesión
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	usuarioID, err := utils.ObtenerUsuarioIDDesdeSession(cookie.Value)
	if err != nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}
	// Obtener prestamo_id de la URL: /descargar/{prestamo_id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	prestamoID := parts[len(parts)-1]

	// Buscar datos del préstamo (ahora incluye fecha_devolucion)
	var libroID, dbUsuarioID int
	var fechaDevolucion, pdfPath string
	query := `SELECT libro_id, usuario_id, fecha_devolucion, (SELECT pdf FROM libros WHERE id = libro_id) FROM prestamos WHERE id = @p1`
	err = database.DB.QueryRow(query, prestamoID).Scan(&libroID, &dbUsuarioID, &fechaDevolucion, &pdfPath)
	if err != nil {
		http.Error(w, "Préstamo no encontrado", http.StatusNotFound)
		return
	}
	if dbUsuarioID != usuarioID {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}
	// Verificar vigencia usando fecha_devolucion
	hoy := utils.Now()
	fechaDev, err := utils.ParseFechaPrestamo(fechaDevolucion)
	if err != nil || hoy.After(fechaDev.Add(24*time.Hour)) {
		// Borra el archivo temporal si existe
		loanPath := utils.LoanPDFPath(usuarioID, libroID)
		os.Remove(loanPath)
		// Redirige a perfil con mensaje de error
		http.Redirect(w, r, "/perfil?error=expirado#tab-borrowed", http.StatusSeeOther)
		return
	}
	// Copiar PDF a carpeta temporal si no existe
	loanPath := utils.LoanPDFPath(usuarioID, libroID)
	if _, err := os.Stat(loanPath); os.IsNotExist(err) {
		// Copiar PDF original
		err = utils.CopiarPDFTemporal(pdfPath, loanPath)
		if err != nil {
			http.Error(w, "No se pudo preparar el archivo", http.StatusInternalServerError)
			return
		}
	}
	// Servir el archivo
	w.Header().Set("Content-Disposition", "attachment; filename=libro.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	file, err := os.Open(loanPath)
	if err != nil {
		http.Error(w, "No se pudo abrir el archivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

// EliminarPrestamoHandler elimina un préstamo por su ID (solo si pertenece al usuario).
func EliminarPrestamoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	usuarioID, err := utils.ObtenerUsuarioIDDesdeSession(cookie.Value)
	if err != nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	prestamoID := parts[1]
	// Solo permite eliminar si el préstamo es del usuario
	_, err = database.DB.Exec("DELETE FROM prestamos WHERE id = @p1 AND usuario_id = @p2", prestamoID, usuarioID)
	if err != nil {
		http.Error(w, "Error eliminando préstamo", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/perfil#tab-borrowed", http.StatusSeeOther)
}
