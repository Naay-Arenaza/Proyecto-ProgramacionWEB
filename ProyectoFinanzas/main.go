package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	sqlc "ProyectoFinanzas/db/sqlc"

	_ "github.com/lib/pq"
)

type QueriesStruct struct {
	queries *sqlc.Queries
}

func main() {
	//Abrir base de datos
	connStr := "user=postgres password=12345 dbname=proyectos sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Cerrar conexion

	queries := &QueriesStruct{
		queries: sqlc.New(db),
	}

	//Abrir el servidor
	staticDir := "./static"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/", serveForm)

	http.HandleFunc("/movimientos", queries.movimientosHandler)
	//http.HandleFunc("/movimientos/", queries.movimientoHandler)

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

func (q *QueriesStruct) movimientosHandler(w http.ResponseWriter, r *http.Request) {
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
func (q *QueriesStruct) getMovimientos(w http.ResponseWriter, r *http.Request) {
	var movimientos = []sqlc.Movimiento{}

	movimientos, err := q.queries.ListMovimiento(r.Context(), 1)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Error 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movimientos)
}

// POST /movimientos - Crear nuevo movimiento
func (q *QueriesStruct) createMovimiento(w http.ResponseWriter, r *http.Request) {
	var newMovimiento sqlc.CreateMovimientoParams

	err := json.NewDecoder(r.Body).Decode(&newMovimiento) //Decodifica el producto en formato JSON del cuerpo de la peticion

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //Error 404
		return
	}

	movimiento, err1 := q.queries.CreateMovimiento(r.Context(), newMovimiento)

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		//El código http.StatusBadRequest (400) es el estándar de la industria RESTful para fallos en la validación de la entrada del usuario
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //Codigo de estado http.StatusCreated(201)
	json.NewEncoder(w).Encode(movimiento)
}
