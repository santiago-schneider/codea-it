FROM golang:1.22-alpine AS build

RUN apk add --no-cache gcc musl-dev sqlite-dev

# Crear el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el módulo go.mod y go.sum y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o main ./cmd/ltp-app

# Etapa final
FROM alpine:latest

# Crear el directorio de trabajo dentro del contenedor
WORKDIR /app

RUN apk add --no-cache sqlite

# Copiar el ejecutable desde la etapa de compilación
COPY --from=build /app/main .

RUN mkdir -p /app/data

# Exponer el puerto que usa tu aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
