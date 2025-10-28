# Nombre del binario de la aplicación
APP_NAME=my-app
# Variables para la base de datos
DB_URL="postgres://postgres:12345@localhost:5432/proyectos?sslmode=disable"
# Target por defecto
all: build
# Inicia el servidor de desarrollo con recarga en vivo
run: docker/up
	@echo ">= Base de datos lista. Iniciando Air..."
	@air
# Genera código SQLC y Templ
generate:
	@echo ">= Generating SQLC code..."
	@sqlc generate
#@echo ">= Generating Templ code..."
#@templ generate
# Levanta la base de datos en modo 'detached' (-d)
docker/up:
	@echo ">= Iniciando base de datos con Docker Compose..."
	@docker compose up -d
# Baja la base de datos
docker/down:
	@echo ">= Deteniendo base de datos..."
	@docker compose down
# Muestra los logs de la base de datos
docker/logs:
	@docker compose logs -f
# Crea una nueva migración de base de datos, te pide el nombre 
db/migrate:
	@read -p "Enter migration name: " name; \
	atlas migrate diff $name \
	--dir "file://db/migrations" \
	--to "file://db/schema/schema.sql" \
	--dev-url "docker://postgres/16/dev" \
# Construye el binario de la aplicación
build: generate
	@echo ">= Building application..."
	@go build -o $(APP_NAME) .
# Limpia los artefactos de construcción
clean:
	@rm -f $(APP_NAME)