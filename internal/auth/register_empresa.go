package auth

import (
	"context"
	"login/internal/database"
	"login/internal/models"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRequest_empresa estructura de los datos recibidos
type RegisterRequest_empresa struct {
	Email_empresa  string `json:"Email_empresa" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Nombre_empresa string `json:"Nombre_empresa" binding:"required"`
}

// RegisterResponse_empresa estructura de la respuesta de registro
type RegisterResponse_empresa struct {
	Message     string `json:"message"`
	FirebaseUID string `json:"firebase_uid"`
}

// RegisterHandler_empresa maneja el registro del usuario
// @Summary Registra un nuevo usuario
// @Description Crea un nuevo usuario en Firebase y lo guarda en la base de datos local
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest_empresa true "Datos del usuario a registrar"
// @Success 200 {object} RegisterResponse_empresa "Usuario registrado correctamente"
// @Failure 400 {object} RegisterResponse_empresa "Solicitud inválida"
// @Failure 500 {object} RegisterResponse_empresa "Error interno del servidor"
// @Router /register_empresa [post]
func RegisterHandler_empresa(c *gin.Context) {
	var req RegisterRequest_empresa
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el correo ya está registrado como usuario
	var usuario models.Usuario
	if result := database.DB.Where("correo = ?", req.Email_empresa).First(&usuario); result.Error != nil {
		if result.Error.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar si el correo está registrado como usuario"})
			return
		}
	}

	if usuario.Id > 0 { // Usar la variable usuario correctamente
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado como usuario"})
		return
	}

	// Verificar si el correo ya está registrado como empresa
	var empresa models.Usuario_empresa
	if result := database.DB.Where("correo_empresa = ?", req.Email_empresa).First(&empresa); result.Error != nil {
		if result.Error.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar si el correo está registrado como empresa"})
			return
		}
	}

	if empresa.Id_empresa > 0 { // Usar la variable empresa correctamente
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado como empresa"})
		return
	}

	// Crear el usuario en Firebase con email y password
	params := (&auth.UserToCreate{}).
		Email(req.Email_empresa).
		Password(req.Password)

	user, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario en Firebase: " + err.Error()})
		return
	}

	// Crear el usuario en la base de datos sin almacenar la contraseña
	usuario_empresa := models.Usuario_empresa{
		Correo_empresa:           req.Email_empresa,
		Nombre_empresa:           req.Nombre_empresa,
		Perfil_Completado:        false,
		Firebase_usuario_empresa: user.UID,
		Rol:                      "empresa", // Rol por defecto
	}

	if result := database.DB.Create(&usuario_empresa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el usuario en la base de datos"})
		return
	}

	// Generar token de verificación de correo
	token, err := GenerateVerificationToken(req.Email_empresa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token de verificación"})
		return
	}

	// Enviar correo de verificación
	err = SendVerificationEmail(req.Email_empresa, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo de verificación"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Usuario empresa creado correctamente. Verifica tu correo", "firebase_uid": user.UID})
}
