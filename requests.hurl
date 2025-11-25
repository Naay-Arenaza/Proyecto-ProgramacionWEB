# 1. Iniciar sesión y obtener la cookie de sesión
POST http://localhost:8080/signin
[FormParams]
usuario: ana@gmail.com
password: 123
HTTP 200
###

# 2. Acceder a la ruta protegida (hURL maneja la cookie automáticamente)
GET http://localhost:8080/
HTTP 200
[Asserts]
# Validamos que el título de la página sea el correcto
body contains "MovFinanzas"
# Validamos que aparezca el título principal 
body contains "Gestion de Movimientos"
# Validamos que aparezca el botón de cerrar sesión
body contains "Cerrar Sesión"
###

# 1. CREAR UN MOVIMIENTO (POST)
POST http://localhost:8080/movimientos
Content-Type: application/x-www-form-urlencoded
[FormParams]
monto: 1200.00
tipo: I
descripcion: Movimiento Base para Test
fechaMovimiento: 2025-11-20
HTTP 200

# 2. INTENTAR CREAR CON ERROR DE MONTO
POST http://localhost:8080/movimientos
Content-Type: application/x-www-form-urlencoded
[FormParams]
monto: -50.00
tipo: G
descripcion: Monto Invalido
fechaMovimiento: 2025-11-20

HTTP 400
[Asserts]
status == 400

# 3. NAVEGAR A LA PÁGINA PRINCIPAL Y CAPTURAR EL ID (GET )
GET http://localhost:8080/
HTTP 200
[Captures]
# Captura el ID de la URL del formulario de Editar, que es la única fuente de ID en el HTML.
# El patrón busca la primera URL que contiene /movimientos/edit/ y extrae el número que le sigue.
mov_id_captured: regex "movimientos/edit/(\\d+)"
[Asserts]
body contains "Movimiento Base para Test"
header "Content-Type" contains "text/html"

# 4. CARGAR PÁGINA DE EDICIÓN (GET /movimientos/edit/{{mov_id_captured}})
GET http://localhost:8080/movimientos/edit/{{mov_id_captured}}
HTTP 200
[Asserts]
status == 200
# El aserto se ajusta al valor real de la BD: 1200.00
body contains "value=\"1200.00\""
body contains "Editar Movimiento"

# 5. ACTUALIZAR EL MOVIMIENTO (POST con _method=PUT)
POST http://localhost:8080/movimientos/{{mov_id_captured}}
Content-Type: application/x-www-form-urlencoded
[FormParams]
_method: PUT
monto: 5000.00
tipo: G
descripcion: Gasto actualizado Hurl (5000)
fechaMovimiento: 2025-11-19

HTTP 303

# 6. ELIMINAR EL MOVIMIENTO (POST con _method=DELETE)
DELETE http://localhost:8080/movimientos/{{mov_id_captured}}
Content-Type: application/x-www-form-urlencoded

HTTP 200

# 3. Refrescar la sesión (obtiene una nueva cookie)
POST http://localhost:8080/refresh
HTTP 200
###

# 4. Cerrar la sesión
GET http://localhost:8080/logout
HTTP 200