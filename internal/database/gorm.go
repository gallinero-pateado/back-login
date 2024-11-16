package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase inicializa la conexión con la base de datos
func InitDatabase() error {
	// Construir la URL de conexión a PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require disable_prepared_statements=true",
		os.Getenv("SUPABASE_HOST"),
		os.Getenv("SUPABASE_USER"),
		os.Getenv("SUPABASE_PASSWORD"),
		os.Getenv("SUPABASE_DB"),
		os.Getenv("SUPABASE_PORT"))

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
		return err
	}

	// Limpiar declaraciones preparadas previas
	db.Exec("DEALLOCATE ALL")

	DB = db
	log.Println("Conexión a la base de datos exitosa con GORM")
	return nil
}
