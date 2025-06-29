# Libros Electrónicos

**Autor:** Erick Argoti  
**Fecha:** 22 de junio de 2025

## Objetivo del Programa
Este proyecto es una aplicación web desarrollada en Go para la gestión de libros electrónicos. Permite a los usuarios:
- Visualizar libros disponibles
- Subir y descargar libros en formato PDF
- Realizar préstamos de libros
- Gestionar usuarios y sesiones
- Visualizar detalles y estadísticas de cada libro

## Principales Funcionalidades
- **Autenticación de usuarios:** Inicio de sesión y registro seguro.
- **Gestión de libros:** Agregar, editar, eliminar y visualizar libros electrónicos.
- **Préstamos:** Solicitar y devolver libros, con control de disponibilidad y fechas.
- **Carga y descarga de archivos:** Subida de portadas e integración de archivos PDF.
- **Panel de usuario:** Visualización de libros prestados y perfil personal.
- **Serialización JSON:** Serialización y deserialización de objetos Go a JSON en el archivo `json.go`. Este archivo es solo de ejemplo y no debe ejecutarse junto con la aplicación principal.

## Estructura del Proyecto
- `cmd/` - Código principal de arranque del servidor
- `controllers/` - Lógica de controladores para usuarios y libros
- `models/` - Definición de modelos de datos (Libro, Usuario, Préstamo)
- `routes/` - Definición de rutas HTTP
- `static/` - Archivos estáticos (CSS, JS, imágenes, PDFs)
- `views/` - Plantillas HTML
- `utils/` - Funciones utilitarias
- `json.go` - Ejemplo de serialización y deserialización JSON (no ejecutar junto con la app principal)

## Mejoras a Futuro
- Implementar un sistema DRM (Digital Rights Management) para proteger los archivos PDF descargados y evitar su distribución no autorizada.

---

Este repositorio está diseñado para facilitar la gestión y préstamo de libros electrónicos en un entorno académico o personal.

## Cómo ejecutar este proyecto

1. Clona el repositorio:
   ```sh
   git clone https://github.com/MrChochito/libros-app.git
   cd libros-app
   ```

2. Instala Go desde [https://go.dev/dl/](https://go.dev/dl/).

3. Configura la base de datos:
   - Ejecuta el archivo `base_de_datos.sql` en tu SQL Server local para crear y poblar la base de datos necesaria para el sistema.
   - Puedes generar este archivo siguiendo las instrucciones de la siguiente sección.

4. Ejecuta el proyecto principal:
   ```sh
   go run main.go
   ```

5. Abre tu navegador y accede a `http://localhost:8080`.

---

## Cómo crear y restaurar el archivo base_de_datos.sql en SQL Server

1. Abre SQL Server Management Studio (SSMS) y conéctate a tu servidor.
2. Haz clic derecho sobre la base de datos que quieres exportar.
3. Selecciona **Tareas > Generar scripts...**
4. En el asistente, selecciona **Todos los objetos** y en opciones avanzadas elige **Esquema y datos**.
5. Guarda el archivo como `base_de_datos.sql`.
6. Para restaurar la base de datos en otra PC, abre SSMS, crea una base de datos vacía, abre el archivo `base_de_datos.sql` y ejecútalo sobre esa base de datos.

---
