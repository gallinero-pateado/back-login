package main

import (
	"log"
	"login/api"
	"login/internal/auth"
	"login/internal/database"
	"login/internal/models"
	"login/internal/storage"
	"login/pkg/config"

	_ "login/docs" // Importar la documentación generada

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Cargar las configuraciones desde el archivo .env
	err := config.LoadConfig("config/app.env")
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Inicializar la base de datos
	err = database.InitDatabase()
	if err != nil {
		log.Fatalf("Error inicializando la base de datos: %v", err)
	}

	// Realizar la migración de la tabla fuera de la transacción
	if !database.DB.Migrator().HasTable(&models.Usuario{}) {
		err = database.DB.AutoMigrate(&models.Usuario{})
		if err != nil {
			log.Fatalf("Error al migrar modelos: %v", err)
		}
	}

	// Realizar la migración de la tabla fuera de la transacción
	if !database.DB.Migrator().HasTable(&models.Usuario_empresa{}) {
		err = database.DB.AutoMigrate(&models.Usuario_empresa{})
		if err != nil {
			log.Fatalf("Error al migrar modelos: %v", err)
		}
	}

	// Inicializar Firebase
	err = auth.InitFirebase()
	if err != nil {
		log.Fatalf("Error inicializando Firebase: %v", err)
	}

	// Inicializar Firebase Storage
	_, err = storage.InitFirebaseStorage()
	if err != nil {
		log.Fatalf("Error inicializando Firebase Storage: %v", err)
	}

	//Registrar rutas
	router := api.SetupRoutes()

	// Agregar la ruta de Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Iniciar el servidor
	router.Run(":8080")
}
