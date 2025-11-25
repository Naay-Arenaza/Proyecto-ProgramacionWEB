-- name: CreateUsuario :one
INSERT INTO Usuario ( nombre, apellido, email, contraseña) 
            VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: GetUsuario :one
SELECT * FROM Usuario WHERE id_usuario = $1;

-- name: GetUsuarioMail :one
SELECT * FROM Usuario WHERE email = $1;

