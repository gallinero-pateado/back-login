package api

import (
	"login/internal/auth"
	"login/internal/upload"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "https://practicas.tssw.info", "https://descuentos.tssw.info", "https://roomies.tssw.info", "https://ulink.tssw.info"}, // Cambia el puerto si es necesario
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	router.POST("/register/user", auth.RegisterHandler)
	router.POST("/login/user", auth.UserLoginHandler)
	router.POST("/register_empresa", auth.RegisterHandler_empresa)
	router.POST("/login/company", auth.CompanyLoginHandler)
	router.GET("/verify-email", auth.VerifyEmailHandler)
	router.POST("/password-reset", auth.SendPasswordResetEmailHandler)
	router.POST("/resend-verification", auth.ResendVerificationEmailHandler)

	// Rutas protegidas
	protected := router.Group("/").Use(auth.AuthMiddleware) // Agrupar las rutas protegidas con el middleware
	{
		protected.POST("/complete-profile", auth.CompleteProfileHandler)                // Ruta para completar perfil
		protected.POST("/complete-profile/empresa", auth.CompleteProfileEmpresaHandler) // Ruta para completar perfil
		protected.POST("/upload-image", upload.UploadImageHandler)                      // Ruta para subir imágenes
		protected.GET("/profile-status", auth.GetProfileStatusHandler)                  // Ruta para ver si el perfil esta verificado
		protected.GET("/profile-status-empresa", auth.GetProfileStatusEmpresaHandler)   // Ruta para ver si el perfil esta verificado

	}

	return router
}
