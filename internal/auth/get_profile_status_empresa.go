package auth

import (
	"login/internal/database"
	"login/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProfileStatusResponse representa la respuesta que incluye el estado de PerfilCompletado
type ProfileStatusResponses struct {
	PerfilCompletado bool `json:"perfil_completado"`
}

// GetProfileStatusHandler devuelve el valor de la variable PerfilCompletado de empresa
// @Summary Obtener estado del perfil
// @Description Retorna si el perfil ha sido completado o no
// @Tags profile
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} ProfileStatusResponse "Estado del perfil"
// @Failure 400 {object} string "Datos inválidos"
// @Failure 401 {object} string "Usuario no autenticado"
// @Failure 500 {object} string "Error interno del servidor"
// @Router /profile-status-empresa [get]
func GetProfileStatusEmpresaHandler(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Buscar el usuario por el uid de Firebase
	var empresa models.Usuario_empresa
	result := database.DB.Where("firebase_usuario_empresa = ?", uid).First(&empresa)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario en la base de datos"})
		return
	}

	// Responder con el estado de PerfilCompletado
	c.JSON(http.StatusOK, ProfileStatusResponses{PerfilCompletado: empresa.Perfil_Completado})
}
