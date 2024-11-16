package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// InitFirebaseStorage inicializa el cliente de Firebase Storage
func InitFirebaseStorage() (*firebase.App, error) {
	// Cargar el archivo de credenciales de Firebase (serviceAccountKey.json)
	sa := option.WithCredentialsFile("config/serviceAccountKey.json") // Ruta al archivo de claves de servicio
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, fmt.Errorf("error inicializando Firebase: %v", err)
	}
	return app, nil
}

// UploadFileToFirebase sube un archivo a Firebase Storage y devuelve la URL pública
func UploadFileToFirebase(file *multipart.FileHeader, bucketName string) (string, error) {
	// Inicializar Firebase Storage
	app, err := InitFirebaseStorage()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("error obteniendo cliente de storage: %v", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("error obteniendo bucket: %v", err)
	}

	// Abrir el archivo para subir
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer src.Close()

	// Crear el writer para subir el archivo a Firebase
	wc := bucket.Object(file.Filename).NewWriter(ctx)
	if _, err := io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("error subiendo archivo a Firebase: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error cerrando el writer: %v", err)
	}

	// Retornar la URL pública del archivo subido
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, file.Filename)
	return url, nil
}
