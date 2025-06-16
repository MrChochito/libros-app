# Sistema de Gestión de Libros Electrónicos

Este repositorio contiene el desarrollo de un sistema web basado en programación funcional con Go para gestionar una biblioteca digital de libros electrónicos. El sistema permite la autenticación de usuarios, administración de libros, búsqueda avanzada, y la descarga controlada de contenido digital.

## Objetivo

Diseñar y desarrollar un sistema modular, claro y escalable que facilite la gestión, consulta y descarga de libros electrónicos, alineado con principios de programación estructurada y buenas prácticas de desarrollo web.

---

## Tecnologías Utilizadas

- Lenguaje: Go
- Framework Web: Gorilla Mux
- Almacenamiento: Repositorio en memoria
- Frontend: HTML, CSS, Bootstrap (básico)
- Librerías externas:
  - github.com/gorilla/mux
  - github.com/gorilla/sessions
  - golang.org/x/crypto/bcrypt

---

## Módulos del Sistema

1. Módulo de Autenticación de Usuarios
   - Registro, inicio de sesión, recuperación de contraseña

2. Módulo de Gestión de Libros
   - Alta, baja, modificación y categorización de libros digitales

3. Módulo de Búsqueda y Filtros
   - Búsquedas por título, autor y filtros avanzados

4. Módulo Administrativo
   - Gestión de usuarios, monitoreo del sistema

---

## Funcionalidades Destacadas

- Registro seguro y control de sesiones
- Administración de libros con metadatos
- Búsquedas dinámicas por múltiples campos
- Panel de control para administradores
- Aplicación estructurada en paquetes, con uso de structs, interfaces, y manejo de errores

---

## Alcances y Límites

Alcance:
- Plataforma web funcional para gestión de libros electrónicos
- Soporte para múltiples usuarios y roles

Límites:
- No se incluye base de datos persistente (uso en memoria)
- No contempla sistema de pagos o lector de libros integrado
- No se implementa versión móvil dedicada
