package main

import (
	sqlc "ProyectoFinanzas/db/sqlc"
	"ProyectoFinanzas/logic"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type QLogic struct {
	movCapaLogica *logic.MovCapaLogica
}

func main() {
	//Abrir base de datos
	connStr := "user=postgres password=12345 dbname=proyectos sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Cerrar conexion

	queries := sqlc.New(db)

	movLogic := logic.NewMovimientoLogic(queries)

	qLogic := &QLogic{
		movCapaLogica: movLogic,
	}

	//Abrir el servidor
	staticDir := "./static"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/", serveForm)

	http.HandleFunc("/movimientos", qLogic.movimientosHandler)
	http.HandleFunc("/movimientos/", qLogic.movimientoHandler)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func serveForm(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "static/index.html")
}

// //////// ->  /movimientos
func (q *QLogic) movimientosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q.getMovimientos(w, r)
	case http.MethodPost:
		q.createMovimiento(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /movimientos - Listar todos los movimientos
//
//	Listamos todos los movimiento que haya en la BD, hasta saber como verificar el usuario
func (q *QLogic) getMovimientos(w http.ResponseWriter, r *http.Request) {
	var movimientos = []sqlc.Movimiento{}

	movimientos, err := q.movCapaLogica.ListMovimientoAllLogic(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Error 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movimientos)
}

// POST /movimientos - Crear nuevo movimiento
func (q *QLogic) createMovimiento(w http.ResponseWriter, r *http.Request) {
	var newMovimiento sqlc.CreateMovimientoParams

	err := json.NewDecoder(r.Body).Decode(&newMovimiento) //Decodifica el producto en formato JSON del cuerpo de la peticion

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //Error 404
		return
	}

	movimiento, err1 := q.movCapaLogica.CreateMovimientoLogic(r.Context(), newMovimiento)

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		//El código http.StatusBadRequest (400) es el estándar de la industria RESTful para fallos en la validación de la entrada del usuario
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //Codigo de estado http.StatusCreated(201)
	json.NewEncoder(w).Encode(movimiento)
}

// //////// ->  /movimientos/
func (q *QLogic) movimientoHandler(w http.ResponseWriter, r *http.Request) {
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
	case http.MethodGet:
		q.getMov(w, r, id)
	case http.MethodPut:
		q.updateMov(w, r, id)
	case http.MethodDelete:
		q.deleteMov(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /movimientos/ - Listar movimiento
func (q *QLogic) getMov(w http.ResponseWriter, r *http.Request, id int) {
	movimiento, err := q.movCapaLogica.GetMovimientoLogic(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movimiento)
}

// PUT /movimientos/ - Actualizar movimiento
func (q *QLogic) updateMov(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("ID recibido: %d", id)
	var aux sqlc.UpdateMovimientoParams
	var mov sqlc.Movimiento
	err := json.NewDecoder(r.Body).Decode(&aux)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("ID recibido: %d", id)
	log.Printf("Datos recibidos: %+v", aux)
	aux.IDMovimiento = int32(id)
	mov, err = q.movCapaLogica.UpdateMovimientoLogic(r.Context(), aux)
	log.Printf("Datos para update: %+v", mov)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mov)
}

// POST /movimientos/ - Eliminar movimiento
func (q *QLogic) deleteMov(w http.ResponseWriter, r *http.Request, id int) {
	err := q.movCapaLogica.DeleteMovimientoLogic(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
}
