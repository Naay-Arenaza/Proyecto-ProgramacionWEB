package main

import (
	//"context"

	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	sqlc "ProyectoFinanzas/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	//Abrir base de datos
	connStr := "user=postgres password=12345 dbname=proyectos sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Cerrar conexion

	queries := sqlc.New(db)
	//ctx := context.Background()

	http.HandleFunc("/movimientos/", func(w http.ResponseWriter, r *http.Request) {
		movimientosHandler(w, r, queries)
	})
	//Abrir el servidor
	staticDir := "./static"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/", serveForm)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func movimientosHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Id de movimiento invalido", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getMov(w, r, id, queries)
	case http.MethodPut:
		updateMov(w, r, id, queries)
	case http.MethodDelete:
		deleteMov(w, r, id, queries)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getMov(w http.ResponseWriter, r *http.Request, id int, queries *sqlc.Queries) {
	movimiento, err := queries.GetMovimiento(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movimiento)
}

func updateMov(w http.ResponseWriter, r *http.Request, id int, queries *sqlc.Queries) {
	var aux sqlc.UpdateMovimientoParams
	var mov sqlc.Movimiento
	err := json.NewDecoder(r.Body).Decode(&aux)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mov, err = queries.UpdateMovimiento(r.Context(), aux)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mov)
}

func deleteMov(w http.ResponseWriter, r *http.Request, id int, queries *sqlc.Queries) {
	err := queries.DeleteMovimiento(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
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
