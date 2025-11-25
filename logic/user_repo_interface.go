package logic

import (
	db "ProyectoFinanzas/db/sqlc"
	"context"
)

type UsuarioRepository interface {
	CreateUsuario(ctx context.Context, arg db.CreateUsuarioParams) (db.Usuario, error)
	GetUsuario(ctx context.Context, idUsuario int32) (db.Usuario, error)
	GetUsuarioMail(ctx context.Context, email string) (db.Usuario, error)

	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	GetSession(ctx context.Context, sessionToken string) (db.Session, error)
	DeleteSession(ctx context.Context, sessionToken string) error
}
