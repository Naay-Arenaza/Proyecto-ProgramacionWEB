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

//Espera a que todo el HTML esté cargado y evita errores al intentar seleccionar elementos que aún no existen
document.addEventListener('DOMContentLoaded', () => {

    // Selecciona el formulario
    const form = document.getElementById('form-movimiento');

    // Escucha el evento 'submit' del formulario
    if (form) {
        form.addEventListener('submit', handleFormSubmitMov);
    }


});