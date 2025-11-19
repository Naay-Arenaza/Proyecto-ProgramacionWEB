package handlers

import (
	"ProyectoFinanzas/logic"
	"ProyectoFinanzas/views"

	// "encoding/json"
	// "log"
	"net/http"
	// "strconv"
	// "strings"
)

type MovHandler struct {
	logic *logic.MovCapaLogica
}

func NewMovHandler(l *logic.MovCapaLogica) *MovHandler {
	return &MovHandler{logic: l}
}

// // //////// ->  /movimientos
// func (q *MovHandler) MovimientosHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		q.getMovimientos(w, r)
// 	case http.MethodPost:
// 		q.createMovimiento(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// // GET /movimientos - Listar todos los movimientos
// // Listamos todos los movimiento que haya en la BD, hasta saber como verificar el usuario
// func (h *MovHandler) getMovimientos(w http.ResponseWriter, r *http.Request) {
// 	var movimientos = []sqlc.Movimiento{}

// 	movimientos, err := h.logic.ListMovimientoAllLogic(r.Context())

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError) //Error 500
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(movimientos)
// }

// // POST /movimientos - Crear nuevo movimiento
// func (h *MovHandler) createMovimiento(w http.ResponseWriter, r *http.Request) {
// 	var newMovimiento sqlc.CreateMovimientoParams

// 	err := json.NewDecoder(r.Body).Decode(&newMovimiento) //Decodifica el producto en formato JSON del cuerpo de la peticion

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest) //Error 404
// 		return
// 	}

// 	movimiento, err1 := h.logic.CreateMovimientoLogic(r.Context(), newMovimiento)

// 	if err1 != nil {
// 		http.Error(w, err1.Error(), http.StatusBadRequest)
// 		//El código http.StatusBadRequest (400) es el estándar de la industria RESTful para fallos en la validación de la entrada del usuario
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated) //Codigo de estado http.StatusCreated(201)
// 	json.NewEncoder(w).Encode(movimiento)
// }

// // //////// ->  /movimientos/
// func (h *MovHandler) MovimientoHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("ID recibido: ")
// 	parts := strings.Split(r.URL.Path, "/")
// 	if len(parts) != 3 {
// 		http.Error(w, "Invalid URL", http.StatusBadRequest)
// 		return
// 	}
// 	id, err := strconv.Atoi(parts[2])
// 	log.Printf("ID recibido: %d", id)
// 	if err != nil {
// 		http.Error(w, "Id de movimiento invalido", http.StatusBadRequest)
// 		return
// 	}

// 	switch r.Method {
// 	case http.MethodGet:
// 		h.getMov(w, r, id)
// 	case http.MethodPut:
// 		h.updateMov(w, r, id)
// 	case http.MethodDelete:
// 		h.deleteMov(w, r, id)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// especificar header?
func (h *MovHandler) getMovimientos(w http.ResponseWriter, r *http.Request) {
	movimientos, err := h.logic.ListMovimientoAllLogic(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Error 500
		return
	}
	views.MovimientoList(movimientos).Render(r.Context(), w)
}

// // GET /movimientos/ - Listar movimiento
// func (h *MovHandler) getMov(w http.ResponseWriter, r *http.Request, id int) {
// 	movimiento, err := h.logic.GetMovimientoLogic(r.Context(), int32(id))

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(movimiento)
// }

// // PUT /movimientos/ - Actualizar movimiento
// func (h *MovHandler) updateMov(w http.ResponseWriter, r *http.Request, id int) {

// 	var aux sqlc.UpdateMovimientoParams
// 	var mov sqlc.Movimiento

// 	err := json.NewDecoder(r.Body).Decode(&aux)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	aux.IDMovimiento = int32(id)
// 	mov, err = h.logic.UpdateMovimientoLogic(r.Context(), aux)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(mov)
// }

// // POST /movimientos/ - Eliminar movimiento
// func (h *MovHandler) deleteMov(w http.ResponseWriter, r *http.Request, id int) {
// 	err := h.logic.DeleteMovimientoLogic(r.Context(), int32(id))

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNoContent)
// 		return
// 	}
// }
