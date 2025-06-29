// Código para inicializar y mostrar las gráficas en el perfil

document.addEventListener('DOMContentLoaded', function () {
  const btnGraficas = document.getElementById('btn-graficas');
  const graficasSection = document.getElementById('graficas-section');
  if (btnGraficas && graficasSection) {
    btnGraficas.addEventListener('click', function () {
      graficasSection.classList.toggle('hidden');
      if (!graficasSection.classList.contains('hidden')) {
        cargarGraficas();
      }
    });
  }
});

function cargarGraficas() {
  fetch('/api/graficas/mas-prestados')
    .then(res => res.json())
    .then(data => {
      new Chart(document.getElementById('grafica-mas-prestados'), {
        type: 'bar',
        data: {
          labels: data.labels,
          datasets: [{
            label: 'Préstamos',
            data: data.values,
            backgroundColor: '#39e079',
          }]
        },
        options: {responsive: true}
      });
    });

  fetch('/api/graficas/categorias')
    .then(res => res.json())
    .then(data => {
      new Chart(document.getElementById('grafica-categorias'), {
        type: 'doughnut',
        data: {
          labels: data.labels,
          datasets: [{
            data: data.values,
            backgroundColor: ['#39e079', '#51946b', '#0e1a13', '#e8f2ec', '#ffd6d6', '#f8fbfa'],
          }]
        },
        options: {responsive: true}
      });
    });

  fetch('/api/graficas/prestamos-semana')
    .then(res => res.json())
    .then(data => {
      new Chart(document.getElementById('grafica-prestamos-semana'), {
        type: 'line',
        data: {
          labels: data.labels,
          datasets: [{
            label: 'Préstamos por semana',
            data: data.values,
            borderColor: '#0e1a13',
            backgroundColor: 'rgba(57,224,121,0.2)',
            fill: true
          }]
        },
        options: {responsive: true}
      });
    });
}
