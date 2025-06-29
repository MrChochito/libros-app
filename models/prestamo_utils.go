package models

import (
	"libros-app/database"
	"time"
)

// Devuelve true si el usuario tiene un préstamo ACTIVO para el libro indicado
func UsuarioTienePrestamoActivo(libroID int, usuarioID int) (bool, error) {
	query := `SELECT TOP 1 fecha_devolucion, devuelto FROM prestamos WHERE libro_id = @p1 AND usuario_id = @p2 AND devuelto = 0 ORDER BY fecha_devolucion DESC`
	var fechaDevolucion string
	var devuelto bool
	err := database.DB.QueryRow(query, libroID, usuarioID).Scan(&fechaDevolucion, &devuelto)
	if err != nil {
		return false, nil // No hay préstamo activo
	}
	// Intentar parsear en varios formatos
	t, err := time.Parse("2006-01-02", fechaDevolucion)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", fechaDevolucion)
		if err != nil {
			t, err = time.Parse(time.RFC3339, fechaDevolucion)
			if err != nil {
				return false, nil
			}
		}
	}
	hoy := time.Now()
	hoy = time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, hoy.Location())
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	if hoy.After(t) {
		return false, nil // Ya venció si hoy es después de la fecha de devolución
	}
	return true, nil
}

