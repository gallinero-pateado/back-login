# Etapa de construcción
FROM golang:1.23.1-alpine AS build

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia dependencias 
COPY go.mod go.sum ./

# Descarga dependencias 
RUN go mod download

# Copia el código fuente 
COPY . .

# Compila el codigo
RUN go build -o main .

# Etapa de producción
FROM alpine:latest

# Copia el binario compilado desde la etapa de construcción
COPY --from=build /app/main .

# Exponer el puerto de la aplicación
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]
