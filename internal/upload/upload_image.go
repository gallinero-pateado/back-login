package upload

import (
	"login/internal/database"
	"login/internal/models"
	"login/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadImageHandler maneja la subida de imágenes de perfil y actualiza el perfil del usuario
// @Summary Subir una imagen de perfil
// @Description Sube una imagen a Firebase Storage y actualiza el campo de foto de perfil del usuario autenticado
// @Tags upload
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "Imagen a subir"
// @Success 200 {object} map[string]string "URL de la imagen subida y mensaje de éxito"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error al subir la imagen"
// @Router /upload-image [post]
func UploadImageHandler(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ha proporcionado un archivo"})
		return
	}

	// Subir la imagen a Firebase Storage
	url, err := storage.UploadFileToFirebase(file, "ulink-sprint-1.appspot.com") // Reemplaza con tu bucket
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir la imagen", "details": err.Error()})
		return
	}

	// Actualizar el campo FotoPerfil en la base de datos
	var usuario models.Usuario
	result := database.DB.Model(&usuario).Where("firebase_usuario = ?", uid).Update("foto_perfil", url)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la foto de perfil en la base de datos"})
		return
	}

	// Responder con la URL de la imagen subida
	c.JSON(http.StatusOK, gin.H{
		"message": "Imagen subida y foto de perfil actualizada correctamente",
		"url":     url,
	})
}
