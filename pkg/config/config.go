package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Cargar las variables de entorno desde el archivo .env
func LoadConfig(filePath string) error {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("Error cargando el archivo de configuración %s: %v", filePath, err)
	}
	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
