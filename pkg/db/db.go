package db

import (
	"database/sql"
	"os"
	_ "modernc.org/sqlite" 
)

// Использем глобальную переменную
var db *sql.DB

// Определяем строковую константу с SQL командами
const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX scheduler_date ON scheduler (date);`

// Проверяем существование файла переданного в dbFile
// Открываем базу данных и при необходимости создаём таблицу
func Init(dbFile string) error {
	// Получаем значение переменной окружения TODO_DBFILE
	envPath := os.Getenv("TODO_DBFILE")

	// Используем это значение если строчка не пустая
	if envPath != "" {
		dbFile = envPath
	}
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	dbConnection, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	db = dbConnection

	// Создание таблицы при необходимости
	if install {
		_, err = db.Exec(schema)
	if err != nil {
		return err
	}
	}

	return nil
}
