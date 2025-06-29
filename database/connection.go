package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func Conectar() {
	server := "localhost"
	port := 1433
	database := "libros_app"

	connString := fmt.Sprintf(
		"server=%s;port=%d;database=%s;trusted_connection=yes",
		server, port, database,
	)

	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error al abrir conexión:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	log.Println("✅ Conexión a SQL Server exitosa.")

	// Crear tablas si no existen
	crearTablas()
}

func crearTablas() {
	// Crear tabla usuarios
	_, err := DB.Exec(`
		IF NOT EXISTS (
			SELECT * FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_NAME = 'usuarios'
		)
		BEGIN
			CREATE TABLE usuarios (
				id INT IDENTITY(1,1) PRIMARY KEY,
				nombre NVARCHAR(100),
				correo NVARCHAR(100) UNIQUE,
				password NVARCHAR(255),
				avatar NVARCHAR(255)
			)
		END
	`)
	if err != nil {
		log.Fatal("❌ Error creando tabla usuarios:", err)
	}

	// Crear tabla libros
	_, err = DB.Exec(`
		IF NOT EXISTS (
			SELECT * FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_NAME = 'libros'
		)
		BEGIN
			CREATE TABLE libros (
				id INT IDENTITY(1,1) PRIMARY KEY,
				titulo NVARCHAR(255),
				autor NVARCHAR(255),
				imagen NVARCHAR(255),
				resumen NVARCHAR(MAX),
				etiquetas NVARCHAR(255),
				disponible BIT DEFAULT 1,
				vistas INT DEFAULT 0,
				veces_prestado INT DEFAULT 0,
				usuario_id INT,
				FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
			)
		END
	`)
	if err != nil {
		log.Fatal("❌ Error creando tabla libros:", err)
	}

	// Agregar columna pdf si no existe
	_, err = DB.Exec(`
		IF COL_LENGTH('libros', 'pdf') IS NULL
		ALTER TABLE libros ADD pdf NVARCHAR(255)
	`)
	if err != nil {
		log.Fatal("❌ Error agregando columna pdf:", err)
	}

	// Agregar columna slug si no existe
	_, err = DB.Exec(`
		IF COL_LENGTH('libros', 'slug') IS NULL
		ALTER TABLE libros ADD slug NVARCHAR(255) UNIQUE
	`)
	if err != nil {
		log.Fatal("❌ Error agregando columna slug:", err)
	}

	// Crear tabla prestamos si no existe
	_, err = DB.Exec(`
		IF NOT EXISTS (
			SELECT * FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_NAME = 'prestamos'
		)
		BEGIN
			CREATE TABLE prestamos (
				id INT IDENTITY(1,1) PRIMARY KEY,
				usuario_id INT,
				libro_id INT,
				duracion_dias INT,
				fecha_prestamo DATETIME,
				fecha_devolucion DATETIME,
				devuelto BIT DEFAULT 0,
				FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
				FOREIGN KEY (libro_id) REFERENCES libros(id)
			)
		END
	`)
	if err != nil {
		log.Fatal("❌ Error creando tabla prestamos:", err)
	}
}

