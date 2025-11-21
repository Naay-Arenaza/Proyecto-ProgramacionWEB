# Proyecto-ProgramacionWEB
# PROYECTO FINANZAS
Este proyecto es un servidor web escrito en Go con el objetivo de registrar las finanzas del usuario. 

# Requsitos:
    Tener instalado:
        - docker compose
            Pueden encontrar informacion en https://docs.docker.com/engine/install/
        - sqlc  
            go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
        - air 
            go install github.com/air-verse/air@latest 
        - go
            Se puede descargar desde [golang.org](https://golang.org/dl/)
        - hurl
            Se puede descargar desde (https://hurl.dev/docs/installation.html)
        - go install github.com/a-h/templ/cmd/templ@latest
          go get github.com/a-h/templ
          templ generate

## Estructura del proyecto
    ProyectoFinanzas/
    ├── db                          # Lógica de la base de datos
    │   └── queries                 # Consultas SQL 
    │   │   └── movimiento.sql      # Consultas de movimientos financieros
    │   │   └── usuario.sql         # Consultas de usuarios
    │   └── schema                  # Estructura de la base de datos
    │   │   └── schema.sql          # Script para crear tablas y relaciones
    │   └── sqlc                    # Archivos generados por sqlc para acceso a datos
    │   |   └── db.go               # Funciones de conexión y manejo de la base de datos
    │   |   └── models.go           # Modelos de datos usados en la aplicación
    │   |   └── movimiento.sql.go   # Código Go generado para movimientos
    │   |   └── usuario.sql.Go      # Código Go generado para usuarios
    |   └── connectDB.go            # Conexion a la BD
    ├── handlers/                   # Capa de Presentación Web 
    |   └── handlers.go             # Controladores Web 
    ├── logic/                   
    |   └── logicMov.go             # Codigo con la Logica de Negocio (movimientos)
    ├── views/                      # Archivos estáticos para la web (HTML, CSS, JS, imágenes)
    |   └── conteiner.templ         # Contenedor Principal de la Página
    |   └── layout.templ            # Plantilla de la Página
    |   └── mod.templ               # Formulario de Modificación/Edición (update)
    |   └── movimiento_form.templ   # Formulario de Creación
    |   └── movimiento_list.templ   # Listado de Movimientos
    ├── static/                     # Archivos Servidos Directamente
    |   └── styles.css              # Front-end estilo
    ├── tmp/                        # Archivos Temporales
    ├── .air.toml                   # Recarga en vivo de app
    ├── db_test.go                  # Test de operaciones CRUD
    ├── docker-compose.yml          # Configuracion servicios docker
    ├── go.mod                      # Configuración y dependencias del módulo Go
    ├── go.sum                      # Configuración y dependencias del módulo Go
    ├── main.go                     # Punto de entrada y servidor principal
    ├── Makefile                    # Configuracion para automatizar tareas
    ├── my-app                      # Manipulacion de DOM, reaccion a eventos, manipulacion de formulario
    ├── requests.hurl               # Test de operaciones HURL 
    ├── sqlc.yaml                   # Configuracion de SQLC

## Ejecución
1. **Instalar requisitos previos**: 

2. **Abrir proyecto**: 
    Abre la carpeta del proyecto en tu editor de preferencia, por ejemplo, Visual Studio Code.

3. **Ejecutar el servidor**:
    1. Clona el proyecto --> git clone https://github.com/Naay-Arenaza/Proyecto-ProgramacionWEB.git
    2. cd ./Proyecto-ProgramacionWEB
    3. Ejecuta el siguiente comando "make run"
        Esto levantara un volumen docker compose con la base de datos, actualizara dependencias, levantara air e inciara el servidor.

4. **Ejecutar Tests**:
    1. Abrir otra terminar en ./Proyecto-ProgramacionWEB
    2. Ejecuta el siguiente comando "make tests"
        Esto ejecutara en primera instancia el test de operaciones CRUD y luego la pruba de la API.

4. **Acceder a la web**:
    Una vez iniciado, abre tu navegador y visita: http://localhost:8080

## Para acceder al Front-end --> http://localhost:8080

## Comandos de interes Makefile
    make run --> "incializa todos los servicios"
    make docker/up  --> levanta un contenedor docker con la configuración del docker-compose.yml
    make docker/down    --> “Baja el contenedor”
    make docker/logs    --> “Muestra los logs de docker”
    make db/migrate --> “Realiza migraciones de la BD, es necesario migracion inicial y sacar inserts de ProyectoFinanzas/db/schema/schema.sql”    
    make test-Op-CRUD --> “Realiza test de operaciones CRUD queries TP2”
    make test-Prueba-HURL --> prueba de integración de la aplicación web (Prueba el SSR/HTML, validar flujo, ...)
    make tests --> Ejecuta test-Op-CRUD y test-Prueba-HURL
