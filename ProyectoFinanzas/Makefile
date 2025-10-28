# Nombre del binario de la aplicación
APP_NAME=my-app
# Variables para la base de datos
DB_URL="postgres://postgres:12345@localhost:5432/proyectos?sslmode=disable"

# Esta línea detecta si 'docker compose' (v2) o 'docker-compose' (v1) está instalado
COMPOSE_CMD ?= $(shell if docker compose version >/dev/null 2>&1; then \
                    echo "docker compose"; \
                 elif docker-compose --version >/dev/null 2>&1; then \
                    echo "docker-compose"; \
                 else \
                    echo "docker-compose-not-found"; \
                 fi)

# Target por defecto
all: build

# Inicia el servidor
run: check-docker docker/up build
	@echo ">= Instalando Air..."
	@go install github.com/air-verse/air@latest
	@echo ">= Base de datos lista. Iniciando Air..."
	@air

# Genera código SQLC y Templ
generate:
	@echo ">= Generating SQLC code..."
	@sqlc generate
#@echo ">= Generating Templ code..."
#@templ generate

# Construye el binario de la aplicación
build: generate
	@echo ">= Building application..."
	@go build -o $(APP_NAME) .

#Base de Datos y Docker 
# Levanta la base de datos en modo 'detached' (-d)
docker/up: check-docker
	@echo ">= Iniciando base de datos con Docker Compose..."
	@$(COMPOSE_CMD) up -d

# Baja la base de datos
docker/down: check-docker
	@echo ">= Deteniendo base de datos..."
	@$(COMPOSE_CMD) down

# Muestra los logs de la base de datos
docker/logs: check-docker
	@$(COMPOSE_CMD) logs -f

# Crea una nueva migración de base de datos, te pide el nombre 
db/migrate:
	@read -p "Enter migration name: " name; \
	atlas migrate diff $name \
	--dir "file://db/migrations" \
	--to "file://db/schema/schema.sql" \
	--dev-url "docker://postgres/16/dev" 

# Construye el binario de la aplicación
clean:
	@rm -f $(APP_NAME)

# Esta regla comprueba si se encontró un comando de Docker Compose
check-docker:
	@if [ "$(COMPOSE_CMD)" = "docker-compose-not-found" ]; then \
		echo ""; \
		echo "ERROR -> No se encontró 'docker compose' (v2) ni 'docker-compose' (v1)."; \
		echo "Por favor, instala Docker para continuar."; \
		echo ""; \
		exit 1; \
	fi


