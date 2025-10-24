-- name: CreateMovimiento :one
INSERT INTO Movimiento ( id_usuario, monto, tipo, descripcion, fecha_movimiento )
            VALUES ( $1, $2, $3, $4 , $5 )
RETURNING *;

-- name: GetMovimiento :one
SELECT * FROM Movimiento WHERE id_movimiento = $1;

-- name: ListMovimiento :many
SELECT * FROM Movimiento WHERE id_usuario = $1 ORDER BY fecha_movimiento DESC;

-- name: UpdateMovimiento :one
UPDATE Movimiento SET monto = $2, tipo = $3, descripcion = $4, fecha_movimiento = $5 WHERE id_movimiento = $1
RETURNING *;

-- name: DeleteMovimiento :exec
DELETE FROM Movimiento WHERE id_movimiento = $1;
