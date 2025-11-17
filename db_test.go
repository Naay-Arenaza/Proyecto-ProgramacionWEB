package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"

	sqlc "ProyectoFinanzas/db/sqlc"
)

func TestQueries_CRUD(t *testing.T) {

	//Abrir base de datos
	connStr := "user=postgres password=12345 dbname=proyectos sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Cerrar conexion

	queries := sqlc.New(db)
	ctx := context.Background()

	//////////////////////////////////////////////////////////////////////////////
	//USUARIOS

	//CreateUsuario
	createdUser, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:     "Jose",
		Apellido:   "Lopez",
		Email:      "joselopez@example.com",
		Contraseña: "12345"})

	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("Created user: %+v\n", createdUser)

	//GetUsuario
	user, err := queries.GetUsuario(ctx, createdUser.IDUsuario) // Read One
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("Retrieved user: %+v\n", user)

	//////////////////////////////////////////////////////////////////////////////
	//MOVIMIENTO

	//CreateMovimiento
	createdMov, err := queries.CreateMovimiento(ctx, sqlc.CreateMovimientoParams{
		IDUsuario: createdUser.IDUsuario,
		Monto:     1000.00,
		Tipo:      "I",
		Descripcion: sql.NullString{ // Sting que puede ser NULL
			String: "sube",
			Valid:  true, // indica que NO es NULL
		},
		FechaMovimiento: time.Date(2025, 9, 28, 0, 0, 0, 0, time.Local),
		// Año, mes, dia, hora, minuto, segundo, nanosegundo, zona
	})

	if err != nil {
		log.Fatalf("Failed to create movement: %v", err)
	}

	fmt.Printf("Created movement: %+v\n", createdMov)

	//GetMovimiento
	userMov, err := queries.GetMovimiento(ctx, createdMov.IDMovimiento) // Read One
	if err != nil {
		log.Fatalf("Failed to get movement: %v", err)
	}
	fmt.Printf("Retrieved movement: %+v\n", userMov)

	//ListMovimiento
	users, err := queries.ListMovimiento(ctx, createdUser.IDUsuario) // Read Many
	if err != nil {
		log.Fatalf("Failed to list movement: %v", err)
	}
	fmt.Printf("All movement: %+v\n", users)

	//UpdateMovimiento

	_, err = queries.UpdateMovimiento(ctx, sqlc.UpdateMovimientoParams{
		IDMovimiento: createdMov.IDMovimiento,
		Monto:        2000000.00,
		Tipo:         "I",
		Descripcion: sql.NullString{ // Sting que puede ser NULL
			String: "sube",
			Valid:  true, // indica que NO es NULL
		},
		FechaMovimiento: time.Date(2025, 9, 28, 0, 0, 0, 0, time.Local),
		// Año, mes, dia, hora, minuto, segundo, nanosegundo, zona
	})
	if err != nil {
		log.Fatalf("Failed to update movement: %v", err)
	}
	fmt.Println("Usmovementer updated successfully")

	updatedMov, err := queries.GetMovimiento(ctx, createdMov.IDMovimiento)
	if err != nil {
		log.Fatalf("failed to get updated movement: %v", err)
	}

	fmt.Printf("Updated movement: %+v\n", updatedMov)

	//DeleteMovimiento
	err = queries.DeleteMovimiento(ctx, createdMov.IDMovimiento) // Delete
	if err != nil {
		log.Fatalf("Failed to delete movement: %v", err)
	}

	fmt.Println("Movement deleted successfully")

	_, err = queries.GetMovimiento(ctx, createdMov.IDMovimiento)
	if err == sql.ErrNoRows {
		fmt.Println("Movement not found after deletion")
	} else if err != nil {
		log.Fatalf("Failed to get movement after deletion: %v", err)
	}

}
