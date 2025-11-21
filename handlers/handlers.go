package handlers

import (
	sqlc "ProyectoFinanzas/db/sqlc"
	"ProyectoFinanzas/logic"
	"ProyectoFinanzas/views"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

type MovimientoWebHandler struct {
	logic *logic.MovCapaLogica
}

func NewMovimientoWebHandler(l *logic.MovCapaLogica) *MovimientoWebHandler {
	return &MovimientoWebHandler{logic: l}
}

func (h *MovimientoWebHandler) EditMovimientoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/movimientos/edit/")
	id, _ := strconv.Atoi(idStr)

	mov, err := h.logic.GetMovimientoLogic(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r) // Si no se encuentra en la BD, devolvemos 404
			return
		}
		http.Error(w, "Error interno al cargar datos", http.StatusInternalServerError)
		return
	}

	html := views.MovimientoEditForm(mov)
	log.Printf("DEBUG: Cargando Movimiento ID %d, TIPO: '%s'", mov.IDMovimiento, mov.Tipo)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html.Render(r.Context(), w)
}

func (h *MovimientoWebHandler) ServeForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()

	movimientos, err := h.logic.ListMovimientoAllLogic(ctx)

	if err != nil {
		log.Printf("Error al cargar movimientos: %v", err)
		http.Error(w, "Error interno del servidor al cargar datos", http.StatusInternalServerError)
		return
	}
	// Lista de todos los movimientos + createFrom
	comp := views.Container(movimientos)
	templ.Handler(views.Layout("MovFinanzas", comp)).ServeHTTP(w, r)

}

// /////////////////////////////////////////////// --->  /MOVIMIENTOS
func (q *MovimientoWebHandler) MovimientosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		q.PostMovimiento(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ///////////////////////// -> CrearMovimiento
func (h *MovimientoWebHandler) PostMovimiento(w http.ResponseWriter, r *http.Request) {
	var newMovimiento sqlc.CreateMovimientoParams
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parseando formulario", http.StatusBadRequest)
		return
	}

	newMovimiento.IDUsuario = 1
	newMovimiento.Tipo = r.FormValue("tipo")
	monto, err1 := strconv.ParseFloat(r.FormValue("monto"), 64)

	if err1 != nil {
		http.Error(w, "Error de conversion del monto", http.StatusBadRequest)
		return
	}
	newMovimiento.Monto = monto

	newMovimiento.Descripcion = sql.NullString{
		String: r.FormValue("descripcion"),
		Valid:  r.FormValue("descripcion") != "",
	}

	fechaStr := r.FormValue("fechaMovimiento")
	fecha, _ := time.Parse("2006-01-02", fechaStr)
	if !logic.EsFechaValida(fecha) {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}
	newMovimiento.FechaMovimiento = fecha

	_, err2 := h.logic.CreateMovimientoLogic(r.Context(), newMovimiento)

	if err2 != nil {
		log.Printf("Error de lógica al crear: %v", err2)

		// Verificar si el error es por una regla de negocio específica (monto/fecha)
		if strings.Contains(err2.Error(), "el monto del movmiento no puede ser menor o igual a 0") ||
			strings.Contains(err2.Error(), "la fecha debe ser menor a la actual") {

			// Si es un error de validación de negocio, respondemos 400
			http.Error(w, "Error de validación: "+err2.Error(), http.StatusBadRequest)
			return
		}

		// Si es cualquier otro error (BD, conexión, etc.), devolvemos 500
		http.Error(w, "Error interno del servidor al guardar: "+err2.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// /////////////////////////////////////////////// --->  /MOVIMIENTO/
func (h *MovimientoWebHandler) MovimientoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ID recibido: ")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	log.Printf("ID recibido: %d", id)
	if err != nil {
		http.Error(w, "Id de movimiento invalido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// case http.MethodGet:
	// 	h.getMov(w, r, id)
	case http.MethodPost:
		metodoReal := r.FormValue("_method")

		switch metodoReal {
		case "DELETE":
			h.deleteMov(w, r, id)
		case "PUT":
			h.updateMovimiento(w, r, id) // Asumiendo que tienes esta función
		default:
			http.Error(w, "Accion no soportada", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ///////////////////////// -> UPDATE

func (h *MovimientoWebHandler) updateMovimiento(w http.ResponseWriter, r *http.Request, id int) {
	var newMovimiento sqlc.UpdateMovimientoParams
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}
	newMovimiento.IDMovimiento = int32(id)
	newMovimiento.Tipo = r.FormValue("tipo")
	montoStr := r.FormValue("monto")
	sanitizedStr := strings.ReplaceAll(montoStr, ",", ".")

	monto, err := strconv.ParseFloat(sanitizedStr, 64)

	if err != nil {
		fmt.Printf("Error de formato: %v. El valor original era: %s\n", err, monto)
		return
	}
	if !logic.MontoValido(monto) {
		http.Error(w, "Monto invalido", http.StatusBadRequest)
		return
	}
	newMovimiento.Monto = monto

	newMovimiento.Descripcion = sql.NullString{
		String: r.FormValue("descripcion"),
		Valid:  r.FormValue("descripcion") != "",
	}

	fechaStr := r.FormValue("fechaMovimiento")
	fecha, _ := time.Parse("2006-01-02", fechaStr)
	if !logic.EsFechaValida(fecha) {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}
	newMovimiento.FechaMovimiento = fecha

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	h.logic.UpdateMovimientoLogic(r.Context(), newMovimiento)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ///////////////////////// -> DELETE
func (h *MovimientoWebHandler) deleteMov(w http.ResponseWriter, r *http.Request, id int) {
	err := h.logic.DeleteMovimientoLogic(r.Context(), int32(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
