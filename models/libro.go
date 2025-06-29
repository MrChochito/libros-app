package models

import (
	"fmt"
	"libros-app/database"
	"libros-app/utils"
	"time"
)

type Libro struct {
	ID                   int
	Titulo               string
	Autor                string
	Imagen               string
	Resumen              string
	Etiquetas            string
	Disponible           bool
	Vistas               int
	VecesPrestado        int
	UsuarioID            int
	Slug                 string // Nuevo campo para el slug
	DuracionPrestamoDias int    // <-- Asegura que este campo esté presente
	PDF                  string // Ruta al PDF
}

// LibroPrestado representa un libro prestado con info de préstamo
type LibroPrestado struct {
	ID            int
	Titulo        string
	Autor         string
	Imagen        string
	Slug          string
	Disponible    bool
	PrestamoID    int
	DiasRestantes int
	Estado        string // Nuevo campo para el estado
}

func GetLibrosPrestados(usuarioID int) ([]Libro, error) {
	query := `SELECT id, titulo, autor, imagen FROM libros WHERE id_usuario = @p1`
	rows, err := database.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen)
		if err != nil {
			return nil, err
		}
		libros = append(libros, l)
	}

	return libros, nil
}

func ObtenerLibros() ([]Libro, error) {
	query := `SELECT id, titulo, autor, imagen FROM libros ORDER BY id DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen)
		if err != nil {
			return nil, err
		}
		libros = append(libros, l)
	}
	return libros, nil
}

func GetLibroByID(id int) (*Libro, error) {
	query := `SELECT id, titulo, autor, imagen, resumen, etiquetas, disponible, vistas, veces_prestado, id_usuario FROM libros WHERE id = @p1`
	row := database.DB.QueryRow(query, id)
	var l Libro
	err := row.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen, &l.Resumen, &l.Etiquetas, &l.Disponible, &l.Vistas, &l.VecesPrestado, &l.UsuarioID)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// GetLibroBySlug busca un libro por el slug del título
func GetLibroBySlug(slug string) (*Libro, error) {
	query := `SELECT id, titulo, autor, imagen, resumen, etiquetas, disponible, vistas, veces_prestado, id_usuario, slug, duracion_prestamo_dias, pdf FROM libros WHERE slug = @p1`
	row := database.DB.QueryRow(query, slug)
	var l Libro
	err := row.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen, &l.Resumen, &l.Etiquetas, &l.Disponible, &l.Vistas, &l.VecesPrestado, &l.UsuarioID, &l.Slug, &l.DuracionPrestamoDias, &l.PDF)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func ObtenerLibrosRecientes() ([]Libro, error) {
	query := `SELECT TOP 8 id, titulo, autor, imagen, slug FROM libros ORDER BY id DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen, &l.Slug)
		if err != nil {
			return nil, err
		}
		libros = append(libros, l)
	}
	return libros, nil
}

func ObtenerLibrosMasVistos() ([]Libro, error) {
	query := `SELECT TOP 8 id, titulo, autor, imagen, slug FROM libros ORDER BY vistas DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen, &l.Slug)
		if err != nil {
			return nil, err
		}
		libros = append(libros, l)
	}
	return libros, nil
}

// GetLibrosSubidos devuelve los libros subidos por un usuario
func GetLibrosSubidos(usuarioID int) ([]Libro, error) {
	query := `SELECT id, titulo, autor, imagen, slug FROM libros WHERE id_usuario = @p1 ORDER BY id DESC`
	rows, err := database.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen, &l.Slug)
		if err != nil {
			return nil, err
		}
		libros = append(libros, l)
	}
	return libros, nil
}

// PedirPrestamo registra un préstamo de libro para un usuario
func PedirPrestamo(libroID int, usuarioID int) error {
	// Verificar si ya existe un préstamo activo para este usuario y libro
	var existe int
	queryCheck := `SELECT COUNT(*) FROM prestamos WHERE usuario_id = @p1 AND libro_id = @p2 AND devuelto = 0`
	err := database.DB.QueryRow(queryCheck, usuarioID, libroID).Scan(&existe)
	if err != nil {
		return err
	}
	if existe > 0 {
		return fmt.Errorf("Ya tienes este libro prestado y no lo has devuelto")
	}

	// Obtener la duración predeterminada del préstamo del libro
	var duracionDias *int
	queryDuracion := `SELECT duracion_prestamo_dias FROM libros WHERE id = @p1`
	err = database.DB.QueryRow(queryDuracion, libroID).Scan(&duracionDias)
	if err != nil {
		return err
	}
	// Si el autor no definió duración, usar 21 días por defecto
	finalDuracion := 21
	if duracionDias != nil && *duracionDias > 0 {
		finalDuracion = *duracionDias
	}

	// Calcular fecha de devolución (hoy + duración)
	fechaPrestamo := utils.FechaHoyString() // Debes tener una función que devuelva la fecha actual en formato YYYY-MM-DD
	fechaDevolucion := utils.SumarDiasAFecha(fechaPrestamo, finalDuracion)

	// Insertar el préstamo con fecha de devolución
	queryInsert := `INSERT INTO prestamos (usuario_id, libro_id, duracion_dias, fecha_prestamo, fecha_devolucion) VALUES (@p1, @p2, @p3, @p4, @p5)`
	_, err = database.DB.Exec(queryInsert, usuarioID, libroID, finalDuracion, fechaPrestamo, fechaDevolucion)
	if err != nil {
		return err
	}
	return nil
}

// GetLibrosPrestadosPorUsuario obtiene los libros que el usuario ha pedido prestados (no los que subió)
func GetLibrosPrestadosPorUsuario(usuarioID int) ([]LibroPrestado, error) {
	query := `SELECT l.id, l.titulo, l.autor, l.imagen, l.slug, p.devuelto, p.id, p.fecha_devolucion, p.estado
		FROM prestamos p
		JOIN libros l ON p.libro_id = l.id
		WHERE p.usuario_id = @p1`
	rows, err := database.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []LibroPrestado
	for rows.Next() {
		var l LibroPrestado
		var devuelto bool
		var prestamoID int
		var fechaDevolucion, estado string
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Imagen, &l.Slug, &devuelto, &prestamoID, &fechaDevolucion, &estado)
		if err != nil {
			return nil, err
		}
		l.Disponible = !devuelto
		l.PrestamoID = prestamoID
		l.DiasRestantes = calcularDiasRestantesPorDevolucion(fechaDevolucion)
		// Lógica de estado dinámico
		hoy := utils.Now()
		fechaDev, _ := utils.ParseFechaPrestamo(fechaDevolucion)
		if devuelto {
			l.Estado = "devuelto"
		} else if hoy.After(fechaDev.Add(24 * time.Hour)) {
			l.Estado = "vencido"
		} else {
			l.Estado = "activo"
		}
		libros = append(libros, l)
	}
	return libros, nil
}

// calcularDiasRestantesPorDevolucion calcula los días restantes usando la fecha de devolución
func calcularDiasRestantesPorDevolucion(fechaDevolucion string) int {
	t, err := utils.ParseFechaPrestamo(fechaDevolucion)
	if err != nil {
		fmt.Println("Error parseando fechaDevolucion:", fechaDevolucion, err)
		return 0
	}
	// Solo comparar fechas (sin horas)
	hoy := utils.Now()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	hoy = time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, hoy.Location())
	dias := int(t.Sub(hoy).Hours() / 24)
	fmt.Printf("fechaDevolucion: %v, hoy: %v, dias: %d\n", t.Format("2006-01-02"), hoy.Format("2006-01-02"), dias)
	if dias < 0 {
		return 0
	}
	return dias + 1 // Incluye el día de devolución
}

// Eliminar préstamo si la fecha de devolución ya pasó
func EliminarPrestamosVencidos() error {
	_, err := database.DB.Exec(`DELETE FROM prestamos WHERE fecha_devolucion < CONVERT(date, GETDATE()) AND devuelto = 0`)
	return err
}

// Eliminar préstamo si la fecha de devolución ya pasó
func MarcarPrestamosVencidosComoDevueltos() error {
	_, err := database.DB.Exec(`UPDATE prestamos SET devuelto = 1 WHERE fecha_devolucion < CONVERT(date, GETDATE()) AND devuelto = 0`)
	return err
}
