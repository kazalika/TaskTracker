package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitPostgreSQLClient() *sql.DB {
	// Строка подключения к базе данных PostgreSQL
	connectionString := "host=postgresql port=5432 user=main_user password=very_strong_generated_password dbname=task_service_db sslmode=disable"

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return nil
	}

	// Проверяем подключение к базе данных
	err = db.Ping()
	if err != nil {
		fmt.Println("Ошибка при проверке подключения к базе данных:", err)
		return nil
	}

	fmt.Println("Подключение к базе данных PostgreSQL успешно установлено!")

	return db
}
