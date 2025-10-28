curl localhost:8080/movimientos

curl localhost:8080/movimientos/2

curl -X DELETE localhost:8080/movimientos/2

curl -X POST http://localhost:8080/movimientos -H "Content-Type: application/json" -d '{
        "id_usuario": 2,
        "monto": 200.00,
        "tipo": "I",
        "descripcion": null,
        "fecha_movimiento": "2024-01-17T00:00:00Z"
  }'

curl -X POST http://localhost:8080/movimientos -H "Content-Type: application/json" -d '{
        "id_usuario": 2,
        "monto": 200.00,
        "tipo": "I",
        "descripcion": {"String":"PEPE","Valid":true},
        "fecha_movimiento": "2024-01-17T00:00:00Z"
    }'

curl -X PUT http://localhost:8080/movimientos/4 -H "Content-Type: application/json" -d '{ 
        "monto": 50000.00,
        "tipo": "I",
        "descripcion": {"String": "Bono BUENISIMO", "Valid": true},
        "fecha_movimiento": "2024-01-20T00:00:00Z"
    }' 