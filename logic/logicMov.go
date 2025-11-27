package logic

import (
	db "ProyectoFinanzas/db/sqlc"
	"context"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type MovCapaLogica struct {
	repo MovimientoRepository
}

func NewMovimientoLogic(r MovimientoRepository) *MovCapaLogica {
	return &MovCapaLogica{
		repo: r,
	}
}

func (l *MovCapaLogica) ListMovimientoAllLogic(ctx context.Context) ([]db.Movimiento, error) {
	return l.repo.ListMovimientoAll(ctx)
}

func (l *MovCapaLogica) ListMovimientoLogic(ctx context.Context, id int32) ([]db.Movimiento, error) {
	return l.repo.ListMovimiento(ctx, id)
}

func (l *MovCapaLogica) CreateMovimientoLogic(ctx context.Context, arg db.CreateMovimientoParams) (db.Movimiento, error) {

	if !MontoValido(arg.Monto) {
		return db.Movimiento{}, errors.New("el monto del movmiento no puede ser menor o igual a 0")
	}

	if !EsFechaValida(arg.FechaMovimiento) {
		return db.Movimiento{}, errors.New("la fecha debe ser menor a la actual")
	}

	return l.repo.CreateMovimiento(ctx, arg)
}

func (l *MovCapaLogica) GetMovimientoLogic(ctx context.Context, arg db.GetMovimientoParams) (db.Movimiento, error) {
	return l.repo.GetMovimiento(ctx, arg)
}

func (l *MovCapaLogica) UpdateMovimientoLogic(ctx context.Context, arg db.UpdateMovimientoParams) (db.Movimiento, error) {

	if !MontoValido(arg.Monto) {
		return db.Movimiento{}, errors.New("el monto del movmiento no puede ser menor o igual a 0")
	}

	if !EsFechaValida(arg.FechaMovimiento) {
		return db.Movimiento{}, errors.New("la fecha debe ser menor a la actual")
	}

	return l.repo.UpdateMovimiento(ctx, arg)
}

func (l *MovCapaLogica) DeleteMovimientoLogic(ctx context.Context, id int32) error {
	return l.repo.DeleteMovimiento(ctx, id)
}

func MontoValido(monto float64) bool {
	if monto <= 0.0 {
		return false
	}
	return true
}

func EsFechaValida(fechaIngresada time.Time) bool {
	horaActual := time.Now()

	if fechaIngresada.After(horaActual) {
		return false
	}
	return true
}
