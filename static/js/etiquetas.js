// Script para gestión visual de etiquetas en subir y editar libro
// Debe incluirse en upload.html y editar_libro.html

const etiquetasDisponibles = [
  "Ficción", "Historia", "Ciencia", "Infantil", "Aventura", "Romance", "Educativo", "Misterio", "Fantasía", "Tecnología"
]; // Puedes cargar esto dinámicamente si lo deseas

function renderEtiquetasSugeridas() {
  const sugeridas = document.getElementById('etiquetas-sugeridas');
  if (!sugeridas) return;
  sugeridas.innerHTML = '';
  etiquetasDisponibles.forEach(etiqueta => {
    if (!etiquetasSeleccionadas.includes(etiqueta)) {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'bg-[#e8f2ec] text-[#0e1a13] px-3 py-1 rounded flex items-center gap-1';
      btn.innerHTML = etiqueta + ' <span class="text-[#39e079] font-bold">+</span>';
      btn.onclick = () => agregarEtiqueta(etiqueta);
      sugeridas.appendChild(btn);
    }
  });
  sugeridas.style.display = 'flex';
}

let etiquetasSeleccionadas = [];

function agregarEtiqueta(etiqueta) {
  const normalizada = etiqueta.trim().toLowerCase();
  if (!normalizada) return;
  // Evita duplicados ignorando mayúsculas/minúsculas y guiones
  if (!etiquetasSeleccionadas.some(e => e.trim().toLowerCase() === normalizada)) {
    etiquetasSeleccionadas.push(etiqueta.trim());
    actualizarEtiquetasSeleccionadas();
    renderEtiquetasSugeridas();
  }
}

function quitarEtiqueta(etiqueta) {
  const normalizada = etiqueta.trim().toLowerCase();
  etiquetasSeleccionadas = etiquetasSeleccionadas.filter(e => e.trim().toLowerCase() !== normalizada);
  actualizarEtiquetasSeleccionadas();
  renderEtiquetasSugeridas();
}

function actualizarEtiquetasSeleccionadas() {
  const seleccionadas = document.getElementById('etiquetas-seleccionadas');
  const input = document.getElementById('etiquetas-input');
  if (!seleccionadas || !input) return;
  seleccionadas.innerHTML = '';
  etiquetasSeleccionadas.forEach(etiqueta => {
    const tag = document.createElement('span');
    tag.className = 'bg-[#39e079] text-white px-3 py-1 rounded flex items-center gap-1 mr-2 mb-2';
    tag.innerHTML = etiqueta + ' <button type="button" class="ml-1 text-white font-bold" onclick="quitarEtiqueta(\'' + etiqueta + '\')">-</button>';
    seleccionadas.appendChild(tag);
  });
  input.value = etiquetasSeleccionadas.join(',');
}

document.addEventListener('DOMContentLoaded', function () {
  renderEtiquetasSugeridas();
  // Inicializar etiquetas seleccionadas desde el valor inicial del input oculto (para edición)
  const input = document.getElementById('etiquetas-input');
  if (input && input.value) {
    etiquetasSeleccionadas = input.value.split(',').map(e => e.trim()).filter(e => e);
    actualizarEtiquetasSeleccionadas();
    renderEtiquetasSugeridas();
  } else {
    actualizarEtiquetasSeleccionadas();
  }
  const buscador = document.getElementById('etiqueta-buscador');
  if (buscador) {
    buscador.addEventListener('keydown', function (e) {
      if (e.key === 'Enter' && buscador.value.trim() !== '') {
        e.preventDefault();
        agregarEtiqueta(buscador.value.trim());
        buscador.value = '';
      }
    });
  }
});

