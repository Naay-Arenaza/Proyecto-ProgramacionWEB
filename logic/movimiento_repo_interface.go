package logic

import (
	db "ProyectoFinanzas/db/sqlc"
	"context"
)

type MovimientoRepository interface {
	CreateMovimiento(ctx context.Context, arg db.CreateMovimientoParams) (db.Movimiento, error)
	GetMovimiento(ctx context.Context, idMovimiento int32) (db.Movimiento, error)
	UpdateMovimiento(ctx context.Context, arg db.UpdateMovimientoParams) (db.Movimiento, error)
	DeleteMovimiento(ctx context.Context, idMovimiento int32) error
	ListMovimiento(ctx context.Context, idUsuario int32) ([]db.Movimiento, error)
	ListMovimientoAll(ctx context.Context) ([]db.Movimiento, error)
}
