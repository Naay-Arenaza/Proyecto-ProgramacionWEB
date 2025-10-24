# PROYECTO FINANZAS
Este proyecto es un servidor web escrito en Go con el objetivo de registrar las finanzas del usuario. 

## Estructura del proyecto
ProyectoFinanzas/
├── db                       # Lógica de la base de datos
│   └── queries              # Consultas SQL 
│   │   └── movimiento.sql   # Consultas de movimientos financieros
│   │   └── usuario.sql      # Consultas de usuarios
│   └── schema               # Estructura de la base de datos
│   │   └── schema.sql       # Script para crear tablas y relaciones
│   └── sqlc                 # Archivos generados por sqlc para acceso a datos
│       └── db.go            # Funciones de conexión y manejo de la base de datos
│       └── models.go        # Modelos de datos usados en la aplicación
│       └── movimiento.sql.go# Código Go generado para movimientos
│       └── usuario.sql.Go   # Código Go generado para usuarios
├── static/                  # Archivos estáticos para la web (HTML, CSS, JS, imágenes)
│   └── index.html           # Página principal del sitio web
├── go.mod                   # Configuración y dependencias del módulo Go
├── main.go                  # Punto de entrada y servidor principal

## Ejecución
1. **Instalar Go**: 
    Tener Go instalado en tu sistema (Se puede descargar desde [golang.org](https://golang.org/dl/))

2. **Abrir proyecto**: 
    Abre la carpeta del proyecto en tu editor de preferencia, por ejemplo, Visual Studio Code.

3. **Ejecutar el servidor**:
   Abre una terminal en la raíz del proyecto y ejecuta el siguiente comando "go run ."
   Esto iniciará el servidor web localmente.

4. **Acceder a la web**:
    Una vez iniciado, abre tu navegador y visita: http://localhost:8080
