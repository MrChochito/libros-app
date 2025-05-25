# Sistema de Gestión de Libros Electrónicos

Este repositorio contiene el desarrollo de un sistema web basado en **programación funcional con Python** para gestionar una biblioteca digital de libros electrónicos. El sistema permite la autenticación de usuarios, administración de libros, búsqueda avanzada, y la descarga controlada de contenido digital.

## Objetivo

Diseñar y desarrollar un sistema modular, claro y escalable que facilite la gestión, consulta y descarga de libros electrónicos, alineado con principios de programación funcional y buenas prácticas de desarrollo web.

---

##  Tecnologías Utilizadas

- **Lenguaje**: Python (con enfoque funcional)
- **Framework Web**: Flask
- **Base de Datos**: SQLite / PostgreSQL
- **Frontend**: HTML, CSS, Bootstrap
- **Librerías externas**:
  - `Flask`
  - `Flask-Login`
  - `SQLAlchemy`
  - `Marshmallow`
  - `functools`, `itertools`, `pathlib`

---

##  Módulos del Sistema

1. **Módulo de Autenticación de Usuarios**
   - Registro, inicio de sesión, recuperación de contraseña

2. **Módulo de Gestión de Libros**
   - Alta, baja, modificación y categorización de libros digitales

3. **Módulo de Búsqueda y Filtros**
   - Búsquedas por título, autor, género y filtros avanzados

4. **Módulo de Descarga**
   - Descarga controlada de archivos con registro de actividad

5. **Módulo Administrativo**
   - Gestión de usuarios, monitoreo del sistema y estadísticas

---

##  Funcionalidades Destacadas

- Registro seguro y control de sesiones
- Administración de libros con metadatos y archivos
- Búsquedas dinámicas por múltiples campos
- Control de descargas y acceso de usuarios
- Panel de control para administradores
- Código enfocado en la inmutabilidad y uso de funciones puras

---

##  Alcances y Límites

**Alcance:**
- Plataforma web completa para gestión de libros electrónicos
- Soporte para múltiples usuarios y roles
- Acceso a estadísticas básicas del sistema

**Límites:**
- No se incluye lector de libros integrado
- No contempla sistema de pago o comercialización
- No se implementa versión móvil (solo responsive)

---
