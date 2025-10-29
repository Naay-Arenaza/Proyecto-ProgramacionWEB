const API_URL = "http://localhost:8080/movimientos";

// Cuando la página carga
document.addEventListener('DOMContentLoaded', function() {
    console.log("Página cargada - Iniciando aplicación");
    cargarMovimientos();
    
    // Escuchar cuando se envía el formulario
    document.getElementById('form-movimiento').addEventListener('submit', agregarMovimiento);
});
// Función para cargar y mostrar los movimientos
async function cargarMovimientos() {
    try {
        console.log("Cargando movimientos desde el servidor...");
        
        const respuesta = await fetch(API_URL);
        
        if (!respuesta.ok) {
            throw new Error('Error al cargar movimientos');
        }
        
        const movimientos = await respuesta.json();
        console.log("MOVIMIENTOS RECIBIDOS:", movimientos);
        
        mostrarMovimientos(movimientos);
        
    } catch (error) {
        console.error("Error:", error);
        document.getElementById('lista_movimientos').innerHTML = `
            <div class="error">Error al cargar los movimientos: ${error.message}</div>
        `;
    }
}

function crearHTMLMovimiento(mov) {
        const esIngreso = mov.tipo === 'I';
        const simbolo = esIngreso ? '+' : '-';
        const clase = esIngreso ? 'ingreso' : 'gasto';
        const texto = esIngreso ? 'INGRESO' : 'GASTO';
        
        // Obtener descripción - varias formas por si la estructura cambia
        let descripcion = 'Sin descripción';
        if (mov.descripcion) {
            if (typeof mov.descripcion === 'string') {
                descripcion = mov.descripcion;
            } else if (mov.descripcion.String) {
                descripcion = mov.descripcion.String;
            } else if (mov.descripcion.string) {
                descripcion = mov.descripcion.string;
            } else if (mov.descripcion.descripcion) {
                descripcion = mov.descripcion.descripcion;
            }
        }
        
        let fechaFormateada = 'Fecha no disponible';
        if (mov.fecha_movimiento) {
            try {
                const fechaISO = mov.fecha_movimiento; // "2025-10-29T00:00:00Z"
                const [fechaPart] = fechaISO.split('T'); // "2025-10-29"
                const [anio, mes, dia] = fechaPart.split('-');
                fechaFormateada = `${dia}/${mes}/${anio}`; // "29/10/2025"
                
                console.log(`Fecha original: ${fechaISO} -> Formateada: ${fechaFormateada}`); 
                
            } catch (e) {
                console.error('Error formateando fecha:', e);
            }
        }
        
        return `
            <div class="movimiento ${clase}">
                <div class="movimiento-header">
                    <span class="fecha">Fecha: ${fechaFormateada}</span>
                    <br>
                    <span class="descripcion">Descripcion: ${descripcion}</span>
                    <br>
                    <span class="monto ${clase}">Monto: ${simbolo}$${mov.monto}</span>
                </div>
                <div class="movimiento-detalles">
                    Tipo de Gasto: ${texto} | ID: ${mov.id_movimiento}
                </div>
                <br>
                </div>
                        <button id="Borrar">Load Data</button>
                </div>
                <br><br>
            </div>
        `;
}

// Función para mostrar los movimientos en pantalla
function mostrarMovimientos(movimientos) {
    const lista = document.getElementById('lista_movimientos');
    
    if (!movimientos || movimientos.length === 0) {
        lista.innerHTML = '<div class="no-movimientos">No hay movimientos registrados</div>';
        return;
    }
    
    const html = movimientos.map(mov => crearHTMLMovimiento(mov)).join('');
    lista.innerHTML = html;
}

// Función para agregar nuevo movimiento
async function agregarMovimiento(evento) {
    evento.preventDefault(); // Evitar que se recargue la página
    
    // Obtener valores del formulario
    const monto = document.getElementById('monto_mov').value;
    const tipo = document.getElementById('tipo').value;
    const descripcionInput = document.getElementById('descripcion_mov').value;
    const fechaMovimiento = document.getElementById('fechaMovimiento').value;
    
    // Validar campos obligatorios
    if (!monto || !tipo || !fechaMovimiento) {
        alert('Por favor completa todos los campos');
        return;
    }
    
    // Crear objeto con los datos
    const nuevoMovimiento = {
        id_usuario: 1,
        monto: parseFloat(monto),
        tipo: tipo,
        descripcion: {
            String: descripcionInput.trim(),
            Valid: descripcionInput.trim() !== ""
        },
        fecha_movimiento: new Date(fechaMovimiento + 'T00:00:00-03:00').toISOString()
    };
    
    console.log("Enviando movimiento:", nuevoMovimiento);
    
    try {
        // Enviar a la API
        const respuesta = await fetch(API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(nuevoMovimiento)
        });
        
        if (respuesta.ok) {
            alert('✅ Movimiento creado correctamente');
            // Limpiar formulario
            document.getElementById('form-movimiento').reset();
            // Recargar la lista de movimientos
            cargarMovimientos();
        } else {
            const errorTexto = await respuesta.text();
            alert(`❌ Error al crear movimiento: ${errorTexto}`);
        }
        
    } catch (error) {
        console.error('Error de conexión:', error);
        alert('❌ Error de conexión con el servidor');
    }
}