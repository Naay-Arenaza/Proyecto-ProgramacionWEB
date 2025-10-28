package logic

import (
	"context"
	"errors"
	"time"

	//"database/sql"
	//"encoding/json"
	//"fmt"
	//"log"
	//"net/http"
	//"strconv"
	//"strings"
	sqlc "ProyectoFinanzas/db/sqlc"

	_ "github.com/lib/pq"
)

type MovCapaLogica struct {
	querie *sqlc.Queries
}

func NewMovimientoLogic(q *sqlc.Queries) *MovCapaLogica {
	return &MovCapaLogica{
		querie: q,
	}
}

func (l *MovCapaLogica) ListMovimientoAllLogic(ctx context.Context) ([]sqlc.Movimiento, error) {
	return l.querie.ListMovimientoAll(ctx)
}

func (l *MovCapaLogica) CreateMovimientoLogic(ctx context.Context, arg sqlc.CreateMovimientoParams) (sqlc.Movimiento, error) {

	if arg.Monto <= 0 {
		return sqlc.Movimiento{}, errors.New("el monto del movmiento no puede ser menor o igual a 0")
	}

	if !esFechaValida(arg.FechaMovimiento) {
		return sqlc.Movimiento{}, errors.New("la fecha debe ser menor a la actual")
	}

	return l.querie.CreateMovimiento(ctx, arg)
}

func (l *MovCapaLogica) GetMovimientoLogic(ctx context.Context, id int32) (sqlc.Movimiento, error) {
	return l.querie.GetMovimiento(ctx, id)
}

func (l *MovCapaLogica) UpdateMovimientoLogic(ctx context.Context, arg sqlc.UpdateMovimientoParams) (sqlc.Movimiento, error) {

	if arg.Monto <= 0 {
		return sqlc.Movimiento{}, errors.New("el monto del movmiento no puede ser menor o igual a 0")
	}

	if !esFechaValida(arg.FechaMovimiento) {
		return sqlc.Movimiento{}, errors.New("la fecha debe ser menor a la actual")
	}

	return l.querie.UpdateMovimiento(ctx, arg)
}

func (l *MovCapaLogica) DeleteMovimientoLogic(ctx context.Context, id int32) error {
	return l.querie.DeleteMovimiento(ctx, id)
}

func esFechaValida(fechaIngresada time.Time) bool {
	horaActual := time.Now()

	if fechaIngresada.After(horaActual) {
		return false
	}
	return true
}
