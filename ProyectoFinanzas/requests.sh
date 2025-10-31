echo "Todos los movimientos: "
curl localhost:8080/movimientos

echo "--------->"

read -p "Ingrese id_movimiento existente: " id

echo "Todos los movimientos con id_movimiento = $id : "
curl localhost:8080/movimientos/$id

echo "--------->"

read -p "Ingrese id_movimiento existente para eliminarlo : " id

echo "Elimina el movimiento con id_movimiento = $id  : "
curl -X DELETE localhost:8080/movimientos/$id

echo "Buscando el movimiento borrado "id_movimiento = $id"  : "
curl localhost:8080/movimientos/$id

echo "--------->"

echo "Creamos un movimiento con todos estos datos (Prueba con descripcion = null) : "
curl -X POST http://localhost:8080/movimientos -H "Content-Type: application/json" -d '{
        "id_usuario": 2,
        "monto": 200.00,
        "tipo": "I",
        "descripcion": null,
        "fecha_movimiento": "2024-01-17T00:00:00Z"
  }'

echo "--------->"

echo "Creamos un movimiento con todos estos datos (Prueba con descripcion != null): "
curl -X POST http://localhost:8080/movimientos -H "Content-Type: application/json" -d '{
        "id_usuario": 1,
        "monto": 23400.00,
        "tipo": "G",
        "descripcion": {"String":"Sanguche","Valid":true},
        "fecha_movimiento": "2024-01-17T00:00:00Z"
    }'

echo "--------->"

read -p "Ingrese id_movimiento existente para actualizar sus valores: " id

echo "Actualizamos el movimiento con id_movimiento = $id : "
curl -X PUT http://localhost:8080/movimientos/$id -H "Content-Type: application/json" -d '{ 
        "monto": 5000.00,
        "tipo": "I",
        "descripcion": {"String": "Chicle", "Valid": true},
        "fecha_movimiento": "2024-01-20T00:00:00Z"
    }' 
echo "--------->"