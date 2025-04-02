package main

import (
	"log"
	"net/http"

	"github.com/arthurskonrad/gotodolist/internal/db"
	"github.com/arthurskonrad/gotodolist/internal/handlers"
)

func main() {
	err := db.Load()
	if err != nil {
		log.Fatal("Erro ao carregar dados:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/add", handlers.AddTodo)
	mux.HandleFunc("/delete", handlers.DeleteTodo)
	mux.HandleFunc("/toggle", handlers.ToggleDone)

	log.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
