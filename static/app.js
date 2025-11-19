const API_URL = "/movimientos"; // URL API (endpoint del servidor)

document.addEventListener('DOMContentLoaded', function() { // DOMContentLoaded -> Espera que el HTML este completamente cargado
    console.log("Página cargada - Iniciando aplicación");
    cargarMovimientos(); // Carga todos los movimientos existentes

    document.getElementById('form-movimiento').addEventListener('submit', agregarMovimiento); //Prepara el formulario para cuando el usuario mande datos

    document.addEventListener('click', function(e) { // Escucha todos los clicks que se hagan en la pagina
        if (e.target.classList.contains('btn-borrar')) { // Solo se ejecuta si el elemento que se clickeo tiene la clase "btn-borrar"
            const id = e.target.getAttribute('data-id'); // Lee el valor del atributo data-id del boton
            console.log("ID del movimiento a borrar:", id);
            borrarMovimiento(id);
        }
    });
});

async function cargarMovimientos() { // async y await -> Esperan la respuesta del servidor sin bloquear la pag.
    try {
        console.log("Cargando movimientos desde el servidor...");
        
        const respuesta = await fetch(API_URL); // Peticion de GET al servidor
        
        if (!respuesta.ok) {
            throw new Error('Error al cargar movimientos');
        }
        
        const movimientos = await respuesta.json();
        console.log("MOVIMIENTOS RECIBIDOS:", movimientos);
        
        mostrarMovimientos(movimientos);
        
    } catch (error) { // Manejo de errores & mensaje a usuario
        console.error("Error:", error);
        document.getElementById('lista_movimientos').innerHTML = `
            <div class="error">Error al cargar los movimientos: ${error.message}</div>
        `;
    }
}

function mostrarMovimientos(movimientos) { // Muestra los movimientos 
    const lista = document.getElementById('lista_movimientos');

    if (!movimientos || movimientos.length === 0) {
        lista.innerHTML = '<div class="no-movimientos">No hay movimientos registrados</div>';
        return;
    }
    
    const html = movimientos.map(mov => crearHTMLMovimiento(mov)).join(''); 
    lista.innerHTML = html;
}

function crearHTMLMovimiento(mov) { //Crear el HTML de cada movimiento
        const esIngreso = mov.tipo === 'I';
        const simbolo = esIngreso ? '+' : '-';
        const clase = esIngreso ? 'ingreso' : 'gasto';
        const texto = esIngreso ? 'INGRESO' : 'GASTO';

        let descripcion = mov.descripcion?.String || 'Sin descripción';

        let fechaFormateada = 'Fecha no disponible';
        if (mov.fecha_movimiento) {
            try {
                const fechaISO = mov.fecha_movimiento; // "2025-10-29T00:00:00Z"
                const [fechaPart] = fechaISO.split('T'); // "2025-10-29"
                const [anio, mes, dia] = fechaPart.split('-');
                fechaFormateada = `${dia}/${mes}/${anio}`; // "29/10/2025"
            } catch (e) {
                console.error('Error formateando fecha:', e);
            }
        }    
        return `
            <div class="movimiento-${clase}-flex">
                <div class="movimiento-detalles">
                    <span class="fecha"> Fecha: ${fechaFormateada}</span>
                        <br>
                    <span class="descripcion"> Descripcion: ${descripcion}</span>
                        <br>
                    <span class="monto ${clase}"> Monto: ${simbolo}$${mov.monto}</span>
                        <br>
                    <span class="gasto"> Tipo de Gasto: ${texto}</span>
                </div>
                <br>
                    <button class="btn-borrar" data-id="${mov.id_movimiento}"> Eliminar </button>
                <br><br>
            </div>
        `;
}

async function borrarMovimiento(idMovimiento) {// Funcion para borrar un movimiento de la lista de movimientos
    if (!confirm('¿Estás seguro de que queres eliminar este movimiento?')) {
        return;
    }
    try {
        const id = parseInt(idMovimiento);
        console.log(`Borrando movimiento ID: ${id}`);
        
        const respuesta = await fetch(`${API_URL}/${id}`, { // Elimina el movimiento del servidor
            method: 'DELETE'
        });
        
        if (respuesta.ok) {
            alert('Movimiento eliminado correctamente');
            cargarMovimientos(); // Recargamos la lista de movimientos
        } else {
            const errorTexto = await respuesta.text();
            alert(`Error al eliminar movimiento: ${errorTexto}`);
        }
        
    } catch (error) {
        console.error('Error de conexión:', error);
        alert('Error de conexión con el servidor');
    }
}

async function agregarMovimiento(evento) { // Funcion para agregar un movimiento
    evento.preventDefault(); // Evitar que se recargue la página
    
    const monto = document.getElementById('monto_mov').value;
    const tipo = document.getElementById('tipo_mov').value;
    const descripcionInput = document.getElementById('descripcion_mov').value;
    const fechaMovimiento = document.getElementById('fechaMovimiento').value;
    const idUsuario = 1; // Por ahora, hasta poder loguear y verificar cada usuario

    if (!monto || !tipo || !fechaMovimiento) { // Validamos los datos que son obligatorios completar
        alert('Por favor completa todos los campos');
        return;
    }
    
    const nuevoMovimiento = { //Creamos el nuevo movimiento con los datos ingresados
        id_usuario: parseFloat(idUsuario),
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
        const respuesta = await fetch(API_URL, { // Envia los datos al servidor
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(nuevoMovimiento)
        });
        
        if (respuesta.ok) {
            alert('Movimiento creado correctamente');
            document.getElementById('form-movimiento').reset(); // Reset formulario
            cargarMovimientos(); // Recargamos la lista de movimientos
        } else {
            const errorTexto = await respuesta.text();
            alert(`Error al crear movimiento: ${errorTexto}`);
        }
    } catch (error) {
        console.error('Error de conexión:', error);
        alert('Error de conexión con el servidor');
    }
}