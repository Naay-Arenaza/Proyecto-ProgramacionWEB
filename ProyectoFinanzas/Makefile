
APP_NAME=my-app

DB_URL="postgres://postgres:12345@localhost:5432/proyectos?sslmode=disable"

all: build

run: docker/up tidy
	@echo "--> Base de datos lista. Iniciando Air ..."
	@air

generate:
	@echo ">= Generando SQLC ..."
	@sqlc generate

tidy:
	@go mod tidy

build: generate
	@echo "--> Iniciando servidor ..."
	@go build -o $(APP_NAME) .

docker/up: 
	@echo ">= Iniciando base de datos con Docker Compose..."
	@docker compose up -d  # El -d es importante para que corra en segundo plano
	@echo ">= Esperando a que la base de datos PostgreSQL esté lista..."
	@until docker exec contenedor_ProyectoFinanzas pg_isready -U postgres -q; do \
        echo "Base de datos no está lista, reintentando en 1s..."; \
        sleep 1; \
    done

docker/down: 
	@echo "--> Deteniendo base de datos ..."
	@docker compose down

docker/logs: 
	@docker compose logs -f

db/migrate:
	@read -p "Enter migration name: " name; \
	atlas migrate diff $name \
	--dir "file://db/migrations" \
	--to "file://db/schema/schema.sql" \
	--dev-url "docker://postgres/16/dev" 

clean:
	@rm -f $(APP_NAME)

test-Op-CRUD: 
	@echo "--> Ejecutando test de las operaciones CRUD: "
	@go test -v 

test-Prueba-API: 
	@echo "--> Ejecutando pruebas API: "
	@chmod +x ./requests.sh
	@./requests.sh

tests: 
	@make test-Op-CRUD
	@make test-Prueba-API