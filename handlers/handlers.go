package handlers

import (
	sqlc "ProyectoFinanzas/db/sqlc"
	"ProyectoFinanzas/logic" // <--- Importa tu capa lógica
	"ProyectoFinanzas/views" // <--- Importa tus vistas
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handle struct {
	Logic *logic.MovCapaLogica
}

// 2. Pide el motor en el constructor
func NewHandler(l *logic.MovCapaLogica) *Handle {
	return &Handle{
		Logic: l,
	}
}

func (h *MovHandler) ServeForm(w http.ResponseWriter, r *http.Request) {
	//Validar ruta
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	movimientos, err := h.logic.ListMovimientoAllLogic(r.Context())
	if err != nil {
		log.Println("Error obteniendo datos:", err)
		http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
		return
	}
	// 4. Renderizar el Layout completo con los datos obtenidos
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Layout(movimientos).Render(r.Context(), w)
}

func (h *MovHandler) MovimientoHandler(w http.ResponseWriter, r *http.Request) {
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
	// case http.MethodPut:
	// 	h.updateMov(w, r, id)
	case http.MethodPost:
		h.deleteMov(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MovHandler) deleteMov(w http.ResponseWriter, r *http.Request, id int) {
	err := h.logic.DeleteMovimientoLogic(r.Context(), int32(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (q *MovHandler) MovimientosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//case http.MethodGet:
	//q.getMovimientos(w, r)
	case http.MethodPost:
		q.PostMovimiento(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MovHandler) PostMovimiento(w http.ResponseWriter, r *http.Request) {
	var newMovimiento sqlc.CreateMovimientoParams
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parseando formulario", http.StatusBadRequest)
		return
	}
	// Asignar valores del formulario a newMovimiento acomodando tipos
	newMovimiento.IDUsuario = 1
	newMovimiento.Tipo = r.FormValue("tipo")
	newMovimiento.Monto, _ = strconv.ParseFloat(r.FormValue("monto"), 64)
	newMovimiento.Descripcion = sql.NullString{
		String: r.FormValue("descripcion"),
		Valid:  r.FormValue("descripcion") != "",
	}
	fechaStr := r.FormValue("fechaMovimiento")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}
	newMovimiento.FechaMovimiento = fecha
	_, err1 := h.logic.CreateMovimientoLogic(r.Context(), newMovimiento)

	if err1 != nil {
		http.Error(w, "Error guardando en base de datos", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
