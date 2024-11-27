package auth

import (
	"context"
	"login/internal/database"
	"login/internal/models"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRequest estructura de los datos recibidos
type RegisterRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nombres    string `json:"nombres" binding:"required"`
	Apellidos  string `json:"apellidos" binding:"required"`
	Id_carrera uint   `json:"Id_carrera"`
}

// RegisterResponse estructura de la respuesta de registro
type RegisterResponse struct {
	Message     string `json:"message"`
	FirebaseUID string `json:"firebase_uid"`
}

// RegisterHandler maneja el registro del usuario
// @Summary Registra un nuevo usuario
// @Description Crea un nuevo usuario en Firebase y lo guarda en la base de datos local
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Datos del usuario a registrar"
// @Success 200 {object} RegisterResponse "Usuario registrado correctamente"
// @Failure 400 {object} RegisterResponse "Solicitud inválida"
// @Failure 500 {object} RegisterResponse "Error interno del servidor"
// @Router /register/user [post]
// RegisterHandler maneja el registro del usuario
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Verificar si el correo ya está registrado como empresa
	var empresa models.Usuario_empresa
	if result := database.DB.Where("correo_empresa = ?", req.Email).First(&empresa); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado como empresa"})
		return
	}

	// Verificar si el correo ya está registrado como usuario
	var usuarioExistente models.Usuario
	if result := database.DB.Where("correo = ?", req.Email).First(&usuarioExistente); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado como usuario"})
		return
	}

	// Crear el usuario en Firebase con email y password
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password)

	user, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario en Firebase: " + err.Error()})
		return
	}

	// Crear el usuario en la base de datos sin almacenar la contraseña
	usuario := models.Usuario{
		Correo:           req.Email,
		Nombres:          req.Nombres,
		Apellidos:        req.Apellidos,
		Firebase_usuario: user.UID,
		Id_carrera:       1,
		Rol:              "estudiante", // Rol por defecto
	}

	result := database.DB.Create(&usuario)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el usuario en la base de datos"})
		return
	}

	// Generar token de verificación de correo
	token, err := GenerateVerificationToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token de verificación"})
		return
	}

	// Enviar correo de verificación
	err = SendVerificationEmail(req.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo de verificación"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Usuario creado correctamente. Verifica tu correo", "firebase_uid": user.UID})
}
