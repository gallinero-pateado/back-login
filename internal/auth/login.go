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

// LoginRequest representa los datos de inicio de sesión
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

// UserLoginHandler maneja el inicio de sesión para usuarios
// @Summary Inicia sesión un usuario
// @Description Autentica al usuario y devuelve un token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Datos de inicio de sesión"
// @Success 200 {object} LoginResponse "Inicio de sesión exitoso"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 401 {object} ErrorResponse "Credenciales incorrectas"
// @Router /login/user [post]
func UserLoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Autenticar con Firebase
	token, err := SignInWithEmailAndPassword(email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Buscar al usuario en la tabla Usuario
	var usuario models.Usuario
	result := database.DB.Where("correo = ?", email).First(&usuario)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Responder con el token JWT y el UID del usuario
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"uid":   usuario.Firebase_usuario,
	})
}

// CompanyLoginHandler maneja el inicio de sesión para empresas
// @Summary Inicia sesión una empresa
// @Description Autentica a una empresa y devuelve un token
// @Tags auth
// @Accept json
// @Produce json
// @Param company body LoginRequest true "Datos de inicio de sesión"
// @Success 200 {object} LoginResponse "Inicio de sesión exitoso"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 401 {object} ErrorResponse "Credenciales incorrectas"
// @Router /login/company [post]
func CompanyLoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Autenticar con Firebase
	token, err := SignInWithEmailAndPassword(email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Buscar a la empresa en la tabla Usuario_empresa
	var usuarioEmpresa models.Usuario_empresa
	result := database.DB.Where("correo_empresa = ?", email).First(&usuarioEmpresa)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Empresa no encontrada"})
		return
	}

	// Responder con el token JWT y el UID de la empresa
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"uid":   usuarioEmpresa.Firebase_usuario_empresa,
	})
}

// SignInWithEmailAndPassword autentica al usuario con Firebase
func SignInWithEmailAndPassword(email, password string) (string, error) {
	apiKey := os.Getenv("FIREBASE_API_KEY")
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + apiKey

	loginPayload := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	jsonPayload, _ := json.Marshal(loginPayload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var firebaseResp FirebaseLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&firebaseResp); err != nil {
		return "", err
	}

	return firebaseResp.IDToken, nil
}
