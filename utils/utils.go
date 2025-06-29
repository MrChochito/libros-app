package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"libros-app/database"
)

// Split divides a string by a separator, returns nil if empty
func Split(s, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

// GetUserIDFromSession obtiene el id del usuario a partir de la cookie de sesión (correo)
func GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		return 0, err
	}
	correo := cookie.Value
	var id int
	err = database.DB.QueryRow("SELECT id FROM usuarios WHERE correo = @p1", correo).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ObtenerUsuarioIDDesdeSession obtiene el id del usuario a partir del valor de la cookie de sesión (correo)
func ObtenerUsuarioIDDesdeSession(sessionValue string) (int, error) {
	correo := sessionValue
	var id int
	err := database.DB.QueryRow("SELECT id FROM usuarios WHERE correo = @p1", correo).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GenerateSlug genera un slug URL amigable a partir de un string
func GenerateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "á", "a")
	s = strings.ReplaceAll(s, "é", "e")
	s = strings.ReplaceAll(s, "í", "i")
	s = strings.ReplaceAll(s, "ó", "o")
	s = strings.ReplaceAll(s, "ú", "u")
	s = strings.ReplaceAll(s, "ñ", "n")
	// Quita caracteres no alfanuméricos ni guiones
	var out strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			out.WriteRune(r)
		}
	}
	return out.String()
}

// GenerateUniqueSlug genera un slug único para un libro, agregando un sufijo incremental si ya existe
func GenerateUniqueSlug(titulo string) string {
	baseSlug := GenerateSlug(titulo)
	slug := baseSlug
	var count int
	for i := 2; ; i++ {
		err := database.DB.QueryRow("SELECT COUNT(*) FROM libros WHERE slug = @p1", slug).Scan(&count)
		if err != nil {
			// Si hay error, por seguridad, retorna el slug base
			return baseSlug
		}
		if count == 0 {
			break
		}
		slug = baseSlug + "-" + strconv.Itoa(i)
	}
	return slug
}

// ParseFechaPrestamo convierte un string a time.Time
func ParseFechaPrestamo(fechaStr string) (time.Time, error) {
	// Intenta con formato completo (datetime)
	t, err := time.Parse("2006-01-02 15:04:05", fechaStr)
	if err == nil {
		return t, nil
	}
	// Intenta solo fecha (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", fechaStr)
	if err == nil {
		return t, nil
	}
	// Intenta formato ISO8601 (2025-07-18T00:00:00Z)
	t, err = time.Parse("2006-01-02T15:04:05Z", fechaStr)
	if err == nil {
		return t, nil
	}
	return time.Time{}, err
}

// PrestamoVigente verifica si el préstamo sigue vigente
func PrestamoVigente(fechaPrestamo time.Time, duracionDias int) bool {
	if duracionDias <= 0 {
		return false
	}
	return time.Now().Before(fechaPrestamo.AddDate(0, 0, duracionDias))
}

// LoanPDFPath devuelve la ruta temporal del PDF prestado
func LoanPDFPath(usuarioID, libroID int) string {
	return filepath.Join("static", "loans", strconv.Itoa(usuarioID), strconv.Itoa(libroID)+".pdf")
}

// CopiarPDFTemporal copia el PDF original a la carpeta temporal del préstamo
func CopiarPDFTemporal(origen, destino string) error {
	fmt.Printf("CopiarPDFTemporal: origen=%s, destino=%s\n", origen, destino)
	// Convertir ruta web a ruta de sistema de archivos
	path := origen
	if strings.HasPrefix(path, "/static/") {
		path = "." + path // relativo al root del proyecto
	}
	os.MkdirAll(filepath.Dir(destino), os.ModePerm)
	in, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error abriendo origen: %v\n", err)
		return err
	}
	defer in.Close()
	out, err := os.Create(destino)
	if err != nil {
		fmt.Printf("Error creando destino: %v\n", err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Printf("Error copiando archivo: %v\n", err)
	}
	return err
}

// Now devuelve la hora actual (wrapper para test/mocks)
func Now() time.Time {
	return time.Now()
}

// FechaHoyString devuelve la fecha actual en formato YYYY-MM-DD
func FechaHoyString() string {
	return time.Now().Format("2006-01-02")
}

// SumarDiasAFecha suma n días a una fecha (YYYY-MM-DD) y devuelve el resultado en el mismo formato
func SumarDiasAFecha(fecha string, dias int) string {
	parsed, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return fecha // fallback: regresa la original si hay error
	}
	return parsed.AddDate(0, 0, dias).Format("2006-01-02")
}
