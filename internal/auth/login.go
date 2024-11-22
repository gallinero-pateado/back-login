package auth

import (
	"bytes"
	"encoding/json"
	"login/internal/database"
	"login/internal/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoginRequest representa los datos de inicio de sesión del usuario
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// FirebaseLoginResponse representa la respuesta de Firebase
type FirebaseLoginResponse struct {
	IDToken string `json:"idToken"`
}

// LoginResponse representa la respuesta del inicio de sesión
type LoginResponse struct {
	Token string `json:"token"`
	UID   string `json:"uid"`
}

// ErrorResponse representa la estructura de un error
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginHandler maneja el inicio de sesión
// @Summary Inicia sesión un usuario
// @Description Autentica al usuario utilizando Firebase y devuelve un token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Datos de inicio de sesión"
// @Success 200 {object} LoginResponse "Inicio de sesión exitoso"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 401 {object} ErrorResponse "Credenciales incorrectas"
// @Router /login [post]
// LoginHandler maneja el inicio de sesión
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	Email := strings.TrimSpace(strings.ToLower(req.Email))

	// Autenticar con Firebase
	token, err := SignInWithEmailAndPassword(Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Buscar al usuario en la tabla Usuario
	var usuario models.Usuario
	result := database.DB.Where("correo = ?", Email).First(&usuario)
	if result.Error != nil {
		// Si no encuentra el usuario en la tabla Usuario, buscar en la tabla Usuario_empresa
		var usuarioEmpresa models.Usuario_empresa
		resultEmpresa := database.DB.Where("correo_empresa = ?", Email).First(&usuarioEmpresa)
		if resultEmpresa.Error != nil {
			// Si no encuentra el usuario en ninguna de las tablas, retornar error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// Responder con el token JWT y el UID del usuario encontrado en Usuario_empresa
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"uid":   usuarioEmpresa.Firebase_usuario_empresa, // Asumiendo que Usuario_empresa también tiene Firebase_usuario
		})
		return
	}

	// Si se encuentra el usuario en la tabla Usuario, responder con el token JWT y el UID
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"uid":   usuario.Firebase_usuario,
	})
}

// SignInWithEmailAndPassword autentica al usuario con Firebase
func SignInWithEmailAndPassword(email, password string) (string, error) {
	apiKey := os.Getenv("FIREBASE_API_KEY") // Cargar la clave de API desde las variables de entorno
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + apiKey

	// Datos a enviar a Firebase
	loginPayload := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	jsonPayload, _ := json.Marshal(loginPayload)

	// Enviar la solicitud a Firebase
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decodificar la respuesta de Firebase
	var firebaseResp FirebaseLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&firebaseResp); err != nil {
		return "", err
	}

	// Retornar el token ID de Firebase
	return firebaseResp.IDToken, nil
}
