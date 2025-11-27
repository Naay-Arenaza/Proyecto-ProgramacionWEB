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

	movLogic := logic.NewMovimientoLogic(queries)

	userLogic := logic.NewUserLogic(queries)

	movWebHandler := handlers.NewMovimientoWebHandler(movLogic, userLogic)

	//Abrir el servidor
	staticDir := "./static"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/signin", movWebHandler.Signin)
	http.HandleFunc("/refresh", movWebHandler.RefreshHandler)
	http.HandleFunc("/logout", movWebHandler.AuthMiddleware(movWebHandler.LogoutHandler))

	http.HandleFunc("/", movWebHandler.AuthMiddleware(movWebHandler.ServeForm))
	http.HandleFunc("/movimientos/edit/", movWebHandler.AuthMiddleware(movWebHandler.EditMovimientoHandler))
	http.HandleFunc("/movimientos", movWebHandler.AuthMiddleware(movWebHandler.MovimientosHandler))
	http.HandleFunc("/movimientos/", movWebHandler.AuthMiddleware(movWebHandler.MovimientoHandler))

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
