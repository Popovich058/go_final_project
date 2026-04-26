package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	
	"go_final_project/pkg/api"
)

func StartingServer() error {
	// Директория с фронтенд-файлами
	webDir := "./web"

	// Получаем значение переменной TODO_PORT
	portStr := os.Getenv("TODO_PORT")

	// Порт по умолчанию
	defaultPort := 7540
	var port int

	// Преобразуем значение переменной в число
	if portStr != "" {
		parsedPort, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("error: TODO_PORT value '%s'. Use a number between 1 and 65535", portStr)
	}

	// Проверяем допустимый диапазон
	if parsedPort < 1 || parsedPort > 65535 {
		return fmt.Errorf("error: port %d is not in the valid range (1–65535)", parsedPort)
	}

	// Если значения TODO_PORT не было, то порт по умолчанию
	port = parsedPort
		} else { 
	port = defaultPort
	}

	// Формируем строку для прослушивания порта
	listenAddress := fmt.Sprintf(":%d", port)
	
	//Вызываем функцию init
	api.Init()

	// Обработчик файл-сервера
	fileServer := http.FileServer(http.Dir(webDir))

	// Обработчик корневого пути
	http.Handle("/", http.StripPrefix("/", fileServer))
	
	//Получеаем значение переменной TODO_DBFILE
	envPath := os.Getenv("TODO_DBFILE")
	
	log.Printf("the server is running at http://localhost%s/", listenAddress)
	log.Printf("TODO_PORT variable: '%s' (port used: %d)", portStr, port)
	if envPath != "" {
		log.Printf("database path overridden by TODO_DBFILE: '%s'", envPath)
	} else {
		log.Printf("using default database path: 'scheduler.db'")
	}
		log.Println("to stop the server, press Ctrl+C")

	return http.ListenAndServe(listenAddress, nil)
}
