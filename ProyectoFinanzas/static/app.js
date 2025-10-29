const API_URL = "/movimientos";
async function handleFormSubmitMov(event){
    event.preventDefault(); // no se si va porque quite el action
    const form = event.currentTarget; // obtiene el form que disparo el evento

    try{
        const monto = document.getElementById('monto_mov').value;
        const tipo = document.getElementById('tipo_mov').value;
        const descripcion = document.getElementById('descripcion_mov').value;
        const fechaMovimiento = document.getElementById('fechaMovimiento').value;
        const idUsuario = 3; // Valor de ejemplo

        // Validaciones básicas
        let fechaISO = '';
        if (fechaMovimiento) {
            fechaISO = new Date(fechaMovimiento + 'T00:00:00-03:00').toISOString();
        } else {
            fechaISO = new Date().toISOString();
        }

        let montoParseado = parseFloat(monto) || 0.0;
        if (montoParseado <= 0) {
            // Si usamos una sola función, un 'alert' y 'return' es una forma
            // simple de detener la ejecución sin usar 'throw'.
            alert('El monto no puede ser menor a 0');
            return; // Detiene la función aquí
        }

        // Creacion del objeto JavaScript 
        const dataParaEnviar = {
            id_usuario: parseInt(idUsuario),
            monto: montoParseado,
            tipo: tipo,
            descripcion: {
                string: descripcion.trim(),
                valid: descripcion.trim() !== ""
            },
            fecha_movimiento: fechaISO
        };

        // guarda la respuesta en response
        const response = await fetch('/movimientos', { // Hace fetch a /movimientos 
            method: 'POST', // metodo del fetch
            headers: {
                'Content-Type': 'application/json' // cabecera, le dice lo que le va a enviar
            },
            body: JSON.stringify(dataParaEnviar) // convierte el js a json, que es lo que a enviar
        });

        // Manejo de la respuesta, guarda en responseText el texto de la respuesta del servidor
        const responseText = await response.text();
        console.log('Status:', response.status, 'Respuesta:', responseText);

        if (response.ok) {
            alert('¡Creación de Movimiento con Éxito!');
            form.reset(); // Limpia el formulario
        } else {
            alert(`Error ${response.status}: ${responseText}`); // error en el servidor (DENTRO)
        }

    } catch (error) {
        // Manejo de errores de conexión (si el fetch falla), osea si no llega ni siquiera el msj
        console.error('Error:', error);
        alert('Error de conexión');
    }
}

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
function mostrarMovimientos(movimientos) {
    const lista = document.getElementById('lista_movimientos');
    
    if (!movimientos || movimientos.length === 0) {
        lista.innerHTML = '<div class="no-movimientos">No hay movimientos registrados</div>';
        return;
    }
    
    const html = movimientos.map(mov => crearHTMLMovimiento(mov)).join('');
    lista.innerHTML = html; // innerHTML muestra el
}
function crearHTMLMovimiento(mov) {
        // va viendo si es ingreso o gasto y en base a eso guarda el valor para desp rellenar
        const esIngreso = mov.tipo === 'I';
        const simbolo = esIngreso ? '+' : '-';
        const clase = esIngreso ? 'ingreso' : 'gasto';
        const texto = esIngreso ? 'INGRESO' : 'GASTO';
        
        // Obtener descripción - varias formas por si la estructura cambia
        let descripcion = mov.descripcion?.String || 'Sin descripción';
        
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
            <div class="movimiento-${clase}-flex">
                 <div class="movimiento-detalles">
                    ID_Movimiento: ${mov.id_movimiento} <br> 
                    Tipo de Gasto: ${texto}
                </div>
                <div class="movimiento-header">
                    <span class="fecha">Fecha: ${fechaFormateada}</span>
                    <br>
                    <span class="descripcion">Descripcion: ${descripcion}</span>
                    <br>
                    <span class="monto ${clase}">Monto: ${simbolo}$${mov.monto}</span>
                </div>
                <button class="borrar_mov">Borrar</button>
            </div>
        `;
}

//Espera a que todo el HTML esté cargado y evita errores al intentar seleccionar elementos que aún no existen
document.addEventListener('DOMContentLoaded', () => {

    // Selecciona el formulario
    const form = document.getElementById('form-movimiento');
    cargarMovimientos();
    // si es null --> no existe form, js evalua false con null
    if (form) {
        form.addEventListener('submit', handleFormSubmitMov); // Escucha el evento 'submit' del formulario

    }
    

});