package auth

import (
	"login/internal/database"
	"login/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsuarios obtiene todos los usuarios
// @Summary Obtiene una lista de todos los usuarios
// @Description Recupera todos los registros de usuarios almacenados en la base de datos
// @Tags usuarios
// @Produce json
// @Success 200 {array} models.Usuario "Lista de usuarios obtenida con éxito"
// @Failure 500 {object} ErrorResponse "Error al obtener los usuarios"
// @Router /GetAllusuarios [get]
// GetAllUsuarios maneja la recuperación de todos los usuarios
func GetAllUsuarios(c *gin.Context) {
	var usuario []models.Usuario

	// Obtener todos los usuarios de la base de datos
	if err := database.DB.Find(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		return
	}

	// Respuesta exitosa con la lista de usuarios
	c.JSON(http.StatusOK, usuario)
}
