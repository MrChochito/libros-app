package models

type Repositorio interface {
	GuardarLibro(libro *Libro) error
	ObtenerLibros() []*Libro
	BuscarPorID(id int) (*Libro, bool)
}

type RepositorioMemoria struct {
	libros []*Libro
}

func NewRepositorioMemoria() *RepositorioMemoria {
	return &RepositorioMemoria{libros: []*Libro{}}
}

func (r *RepositorioMemoria) GuardarLibro(libro *Libro) error {
	r.libros = append(r.libros, libro)
	return nil
}

func (r *RepositorioMemoria) ObtenerLibros() []*Libro {
	return r.libros
}

func (r *RepositorioMemoria) BuscarPorID(id int) (*Libro, bool) {
	for _, l := range r.libros {
		if l.ID == id {
			return l, true
		}
	}
	return nil, false
}
