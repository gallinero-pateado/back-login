package auth

import (
	"login/internal/database"
	"login/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsuariouid trae la fila del usuario específico por UID
// @Summary Obtiene los datos de un usuario por su UID
// @Description Recupera la información de un usuario de la base de datos utilizando su UID de Firebase
// @Tags usuarios
// @Produce json
// @Param uid path string true "UID del usuario"
// @Success 200 {object} models.Usuario "Datos del usuario encontrados"
// @Failure 404 {object} ErrorResponse "Usuario no encontrado"
// @Router /Getusuario/{uid} [get]
// GetUsuariouid maneja la búsqueda del usuario por su UID
func GetUsuariouid(c *gin.Context) {
	var usuario models.Usuario
	uid := c.Param("uid") // Obtener el UID de los parámetros de la ruta

	// Buscar al usuario en la base de datos por su UID de Firebase
	if err := database.DB.First(&usuario, "firebase_usuario = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	// Retornar los datos del usuario en formato JSON
	c.JSON(http.StatusOK, usuario)
}
