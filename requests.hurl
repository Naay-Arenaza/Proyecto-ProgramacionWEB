# 1. CREAR UN MOVIMIENTO (Caso descripción NULL)
POST http://localhost:8080/movimientos
Content-Type: application/json
{
    "id_usuario": 2,
    "monto": 200.00,
    "tipo": "I",
    "descripcion": null,
    "fecha_movimiento": "2024-01-17T00:00:00Z"
}
HTTP 201
[Captures]
# CAMBIO AQUÍ: id -> id_movimiento
mov_id_1: jsonpath "$.id_movimiento" 
[Asserts]
status == 201
jsonpath "$.monto" == 200.00
jsonpath "$.tipo" == "I"
# CAMBIO AQUÍ: id -> id_movimiento
jsonpath "$.id_movimiento" isInteger 

# 2. CREAR SEGUNDO MOVIMIENTO (Caso descripción con valor)
POST http://localhost:8080/movimientos
Content-Type: application/json
{
    "id_usuario": 1,
    "monto": 23400.00,
    "tipo": "G",
    "descripcion": {"String":"Sanguche","Valid":true},
    "fecha_movimiento": "2024-01-17T00:00:00Z"
}
HTTP 201
[Captures]
# CAMBIO AQUÍ: id -> id_movimiento
mov_id_2: jsonpath "$.id_movimiento"
[Asserts]
status == 201
jsonpath "$.descripcion.String" == "Sanguche"

# 3. OBTENER TODOS LOS MOVIMIENTOS
GET http://localhost:8080/movimientos
HTTP 200
[Asserts]
jsonpath "$" count >= 2

# 4. OBTENER EL MOVIMIENTO 1 POR ID
GET http://localhost:8080/movimientos/{{mov_id_1}}
HTTP 200
[Asserts]
# CAMBIO AQUÍ: id -> id_movimiento
jsonpath "$.id_movimiento" == {{mov_id_1}}
jsonpath "$.monto" == 200.00

# 5. ACTUALIZAR EL MOVIMIENTO 1
PUT http://localhost:8080/movimientos/{{mov_id_1}}
Content-Type: application/json
{ 
    "monto": 5000.00,
    "tipo": "I",
    "descripcion": {"String": "Chicle", "Valid": true},
    "fecha_movimiento": "2024-01-20T00:00:00Z"
}
HTTP 200
[Asserts]
jsonpath "$.monto" == 5000.00
jsonpath "$.descripcion.String" == "Chicle"

# 6. ELIMINAR EL MOVIMIENTO 1
DELETE http://localhost:8080/movimientos/{{mov_id_1}}
HTTP 200

# 7. VERIFICAR QUE EL MOVIMIENTO 1 YA NO EXISTE
GET http://localhost:8080/movimientos/{{mov_id_1}}
HTTP 404