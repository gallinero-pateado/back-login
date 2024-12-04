package auth

import (
	"log"
	"login/internal/database"
	"login/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProfileUpdateRequests define los campos para actualizar el perfil de la empresa
type ProfileUpdateRequests struct {
	Sector            string `json:"Sector"`
	Descripcion       string `json:"Descripcion"`
	Direccion         string `json:"Direccion"`
	Persona_contacto  string `json:"Persona_contacto"`
	Correo_contacto   string `json:"Correo_contacto"`
	Telefono_contacto int    `json:"Telefono_contacto"`
	Perfil_Completado bool   `json:"Perfil_Completado"`
}

// CompleteProfileEmpresaHandler permite a los usuarios completar o actualizar su perfil
// @Summary Completar o actualizar perfil de usuario empresa
// @Description Permite a los usuarios autenticados completar o actualizar su perfil de empresa
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param profile body ProfileUpdateRequests true "Datos para actualizar el perfil"
// @Success 200 {object} SuccessResponse "Perfil actualizado correctamente"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 401 {object} ErrorResponse "Usuario no autenticado"
// @Failure 500 {object} ErrorResponse "Error al actualizar el perfil"
// @Router /complete-profile/empresa [post]
func CompleteProfileEmpresaHandler(c *gin.Context) {
	// Verificar si el UID está presente en el contexto
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener los datos del perfil desde el cuerpo de la solicitud
	var req ProfileUpdateRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Verificar si la conexión con la base de datos es válida
	if database.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la conexión a la base de datos"})
		return
	}

	// Buscar la empresa por el uid del Firebase
	var empresa models.Usuario_empresa
	result := database.DB.Where("firebase_usuario = ?", uid).First(&empresa)
	if result.Error != nil {
		log.Println("Error al buscar empresa:", result.Error) // Log para depuración
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
		return
	}

	// Actualizar solo los campos no relacionados con la foto de perfil
	updateResult := database.DB.Model(&empresa).Where("firebase_usuario = ?", uid).Updates(models.Usuario_empresa{
		Sector:            req.Sector,
		Descripcion:       req.Descripcion,
		Direccion:         req.Direccion,
		Persona_contacto:  req.Persona_contacto,
		Correo_contacto:   req.Correo_contacto,
		Telefono_contacto: req.Telefono_contacto,
		Perfil_Completado: req.Perfil_Completado,
	})

	// Manejo de errores en la actualización de los datos
	if updateResult.Error != nil {
		log.Println("Error al actualizar el perfil:", updateResult.Error) // Log para depuración
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el perfil"})
		return
	}

	// Responder con éxito si la actualización es correcta
	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente"})
}
