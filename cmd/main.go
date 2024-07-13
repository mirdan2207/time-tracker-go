package main

import (
	"log"
	"net/http"
	"time-tracker-go/config"
	"time-tracker-go/migrations"
	"time-tracker-go/routes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	log.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established successfully.")

	// Выполнение миграций
	log.Println("Applying migrations...")
	migrations.Migrate(db)

	// Загрузка начальных данных
	log.Println("Seeding initial data...")
	migrations.Seed(db)

	// Настройка маршрутов
	log.Println("Setting up routes...")
	router := routes.SetupRoutes(db, cfg)

	// Запуск сервера
	serverAddr := ":8080"
	log.Printf("Server is running at http://localhost%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
