package main

import (
	"log"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	// Подключаемся к базе данных
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("error connection database: %v", err)
	}
	// Запускаем сервер
	err = server.StartingServer()
	if err != nil {
		log.Fatalf("error when starting the server: %v", err)
		}
	}

