package main

import (
	"ProyectoFinanzas/db"
	sqlc "ProyectoFinanzas/db/sqlc"
	"ProyectoFinanzas/handlers"
	"ProyectoFinanzas/logic"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	//Abrir base de datos
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Cerrar conexion

	queries := sqlc.New(db)
	movLogic := logic.NewMovimientoLogic(queries)  // -> Capa Logica
	movHandler := handlers.NewMovHandler(movLogic) // -> /movimientos y /movimientos/

	//Abrir el servidor
	staticDir := "./static"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/", movHandler.ServeForm)
	http.HandleFunc("/movimientos", movHandler.PostMovimiento)
	http.HandleFunc("/movimientos/", movHandler.MovimientoHandler)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
