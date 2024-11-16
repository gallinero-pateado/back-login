package auth

import (
	"net/http"

	"login/internal/database"
	"login/internal/models"

	"github.com/gin-gonic/gin"
)

// EmailRequest estructura para recibir el email
type EmailRequest struct {
	Email string `json:"email" binding:"required"`
}

// SuccessResponse estructura para respuestas exitosas
type SuccessResponse struct {
	Message string `json:"message"`
}

// ResendVerificationEmailHandler maneja el reenvío del correo de verificación
// @Summary Reenviar correo de verificación
// @Description Reenvía el correo de verificación a un usuario registrado
// @Tags verification
// @Accept json
// @Produce json
// @Param email body EmailRequest true "Correo del usuario"
// @Success 200 {object} SuccessResponse "Correo de verificación enviado nuevamente"
// @Failure 400 {object} ErrorResponse "Email requerido"
// @Failure 404 {object} ErrorResponse "Usuario no encontrado"
// @Failure 500 {object} ErrorResponse "Error interno del servidor"
// @Router /resend-verification [post]
func ResendVerificationEmailHandler(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email requerido"})
		return
	}

	// Buscar al usuario en la base de datos usando su email
	var usuario models.Usuario
	result := database.DB.Where("correo = ?", req.Email).First(&usuario)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Usuario no encontrado"})
		return
	}

	// Generar token de verificación
	token, err := GenerateVerificationToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error al generar el token de verificación"})
		return
	}

	// Enviar correo de verificación
	err = SendVerificationEmail(req.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error al enviar el correo de verificación"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Correo de verificación enviado nuevamente"})
}
