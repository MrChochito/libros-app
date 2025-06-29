package controllers

import (
	"encoding/json"
	"libros-app/database"
	"net/http"
	"strings"
)

// Datos para gráfica de libros más prestados (top 7)
func GraficaMasPrestados(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT TOP 7 titulo, COUNT(*) as total FROM prestamos p JOIN libros l ON p.libro_id = l.id GROUP BY titulo ORDER BY total DESC`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	labels := []string{}
	values := []int{}
	for rows.Next() {
		var titulo string
		var total int
		if err := rows.Scan(&titulo, &total); err == nil {
			labels = append(labels, titulo)
			values = append(values, total)
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"labels": labels, "values": values})
}

// Datos para gráfica de categorías (circular)
func GraficaCategorias(w http.ResponseWriter, r *http.Request) {
	// Si no existe la columna 'categoria', usar etiquetas como fallback
	rows, err := database.DB.Query(`SELECT etiquetas FROM libros WHERE etiquetas IS NOT NULL AND etiquetas <> ''`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	categoriaCount := map[string]int{}
	for rows.Next() {
		var etiquetas string
		if err := rows.Scan(&etiquetas); err == nil {
			for _, cat := range strings.Split(etiquetas, ",") {
				cat = strings.TrimSpace(cat)
				if cat != "" {
					categoriaCount[cat]++
				}
			}
		}
	}
	labels := []string{}
	values := []int{}
	for cat, count := range categoriaCount {
		labels = append(labels, cat)
		values = append(values, count)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"labels": labels, "values": values})
}

// Datos para gráfica de préstamos por semana (últimas 4 semanas)
func GraficaPrestamosSemana(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT FORMAT(fecha_prestamo, 'yyyy-MM-dd') as semana, COUNT(*) as total FROM prestamos WHERE fecha_prestamo >= DATEADD(week, -4, GETDATE()) GROUP BY FORMAT(fecha_prestamo, 'yyyy-MM-dd') ORDER BY semana`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	labels := []string{}
	values := []int{}
	for rows.Next() {
		var semana string
		var total int
		if err := rows.Scan(&semana, &total); err == nil {
			labels = append(labels, semana)
			values = append(values, total)
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"labels": labels, "values": values})
}

