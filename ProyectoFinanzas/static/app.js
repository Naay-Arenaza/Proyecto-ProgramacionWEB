// ¡AÑADE ESTA LÍNEA PRIMERO QUE NADA!
console.log("¡app.js se cargó correctamente!");
// Espera a que todo el HTML esté cargado
document.addEventListener('DOMContentLoaded', () => {

    // 1. Selecciona el formulario
    const form = document.getElementById('form-movimiento');

    // 2. Escucha el evento 'submit' del formulario
    form.addEventListener('submit', async (e) => {
        
        // 3. Evita que el formulario se envíe de la forma tradicional (recargando la página)
        e.preventDefault();

        // 4. Obtiene los valores de los inputs
        //    (Usamos .value para obtener el texto que el usuario escribió)
        const monto = document.getElementById('monto_mov').value;
        const tipo = document.getElementById('tipo_mov').value;
        const descripcion = document.getElementById('descripcion_mov').value;
        const fechaMovimiento = document.getElementById('fechaMovimiento').value;
        //    Por ahora, valor de ejemplo: 1
        const idUsuario = 1; 

        // 5. Validaciones básicas
        if (!fechaMovimiento) {
            alert('Tipo es obligatorio y Fecha son obligatorios');
            return;
        }

        // Convertir fecha a formato ISO con hora y timezone
        let fechaISO = '';
        if (fechaMovimiento) {
            // Agregar hora medianoche y timezone
            fechaISO = new Date(fechaMovimiento + 'T00:00:00-03:00').toISOString();
        } else {
            // Fecha actual si no se especifica
            fechaISO = new Date().toISOString();
        }

        // 6. Creacion del objeto JavaScript (¡CON CORRECCIONES!)
        const dataParaEnviar = {
            id_usuario: parseInt(idUsuario),
            monto: parseFloat(monto) || 0.0,
            tipo: tipo,
            descripcion: {
                string: descripcion.trim(),  // Campo "string" para el valor
                valid: descripcion.trim() !== ""  // Campo "valid" para indicar si es válido
            },
            fecha_movimiento: fechaISO
        };
        // DEBUG: Ver qué estamos enviando
        console.log('Datos a enviar:', dataParaEnviar);

        try {
            const response = await fetch('/movimientos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(dataParaEnviar)
            });

            const responseText = await response.text();
            console.log('Status:', response.status, 'Respuesta:', responseText);

            if (response.ok) {
                alert('¡Creación de Movimiento con Éxito!');
                form.reset();
            } else {
                alert(`Error ${response.status}: ${responseText}`);
            }

        } catch (error) {
            console.error('Error:', error);
            alert('Error de conexión');
        }
    });
});