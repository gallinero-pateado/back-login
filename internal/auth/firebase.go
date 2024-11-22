package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth" // Este es el paquete correcto para autenticación
	"google.golang.org/api/option"
)

var authClient *auth.Client // Cambiado de firebase.Auth a auth.Client

// InitFirebase inicializa Firebase con el archivo de credenciales
func InitFirebase() error {
	opt := option.WithCredentialsFile("config/serviceAccountKey.json") // Ruta a tus credenciales
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error inicializando Firebase: %v", err)
	}

	// Inicializar el cliente de autenticación
	authClient, err = app.Auth(context.Background()) // Utilizar el método correcto para obtener el cliente de auth
	if err != nil {
		return fmt.Errorf("error creando el cliente de autenticación: %v", err)
	}

	log.Println("Firebase inicializado correctamente")
	return nil
}

// VerifyHandler maneja las solicitudes para verificar el token de autenticación
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el token del encabezado Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Falta el token de autorización", http.StatusUnauthorized)
		return
	}

	// Verificar el token con Firebase
	token, err := authClient.VerifyIDToken(context.Background(), authHeader)
	if err != nil {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Si el token es válido, responde con el UID del usuario
	fmt.Fprintf(w, "Token válido! UID: %s", token.UID)
}
