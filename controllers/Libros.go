package controllers

import (
	"html/template"
	"net/http"
	"sistema-libros/models"
	"strings"
)

var repo models.Repositorio

func SetRepositorio(r models.Repositorio) {
	repo = r
}

func FormAgregarLibro(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/agregar_libro.html"))
	tmpl.Execute(w, nil)
}

func AgregarLibro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/libros", http.StatusSeeOther)
		return
	}
	r.ParseForm()

	nextID := len(repo.ObtenerLibros()) + 1
	libro, err := models.NuevoLibro(
		nextID,
		r.FormValue("titulo"),
		r.FormValue("autor"),
		r.FormValue("genero"),
		r.FormValue("archivo"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repo.GuardarLibro(libro)
	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

func MostrarLibros(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/libros.html"))
	tmpl.Execute(w, repo.ObtenerLibros())
}

func BuscarLibro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	criterio := r.FormValue("criterio")
	resultados := []*models.Libro{}

	for _, libro := range repo.ObtenerLibros() {
		if strings.Contains(strings.ToLower(libro.Titulo), strings.ToLower(criterio)) ||
			strings.Contains(strings.ToLower(libro.Autor), strings.ToLower(criterio)) {
			resultados = append(resultados, libro)
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/resultados_busqueda.html"))
	tmpl.Execute(w, resultados)
	
}
func ExportarLibros(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formato := r.FormValue("formato")

	var exportador models.Exportador
	archivo := "libros_exportados"

	switch formato {
	case "json":
		exportador = models.ExportadorJSON{}
		archivo += ".json"
	case "csv":
		exportador = models.ExportadorCSV{}
		archivo += ".csv"
	default:
		http.Error(w, "Formato no válido", http.StatusBadRequest)
		return
	}

	err := exportador.Exportar(repo.ObtenerLibros(), archivo)
	if err != nil {
		http.Error(w, "Error exportando: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Exportación exitosa. Archivo generado: " + archivo))
}
