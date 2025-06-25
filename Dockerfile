# Usa la imagen oficial de Go
FROM golang:1.24-alpine AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia el archivo .env dentro del contenedor
COPY .env .env
# Asegúrate de que el archivo se copie al contenedor

# Copia los archivos de módulo y descarga dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto de tu código
COPY . .

# Compila la aplicación
RUN go build -o auth-aca ./cmd/api

# Usa una imagen mínima para ejecutar la aplicación
FROM alpine:latest

# Copia el binario construido al contenedor
COPY --from=builder /app/auth-aca /auth-aca
# Copia el archivo .env al contenedor final
COPY --from=builder /app/.env /app/.env

# Expone el puerto que utilizará la aplicación
EXPOSE 8090

# Comando que se ejecutará al iniciar el contenedor
CMD ["/auth-aca"]