# go-events

**go-events** is a backend microservice developed in Go that implements an event management and routing system using **gRPC**. It is designed following **Clean Architecture** principles to provide a highly testable, decoupled, and scalable business core.

The service exposes operations for both traditional resource management (`Event`) and Pub/Sub messaging patterns (`Subscription`, `PullMessages`, `AckMessage`).

---

## 🚀 Main Technologies

- **Language:** Go 1.26.4
- **Communications:** gRPC and Protocol Buffers (`protoc`)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Configuration:** Viper
- **Logging:** `slog` (with structured colored formatting in main)
- **Local Infrastructure:** Docker & Docker Compose
- **Quality and Security:** `golangci-lint`, Snyk, Gitleaks

---

## 🏗 Architecture and Project Structure

The project follows a domain-oriented distribution:

```text
├── bin/                 # Support scripts (.sh) for the Makefile
├── cmd/
│   └── app/             # Application entry point (main.go)
├── internal/
│   ├── domain/          # Core entities (Event, Queue, Subscription) and ports (interfaces)
│   ├── infrastructure/  # Input and output adapters (GORM, Viper, Logging, gRPC server)
│   │   ├── configuration/ # Configuration loading (app.env / environment variables)
│   │   ├── database/      # Repository implementation (Postgres)
│   │   ├── gapi/          # gRPC handlers and auto-generated code (.pb.go)
│   │   └── proto/         # gRPC interface definition (.proto)
│   └── testsuite/       # Integration tests and infrastructure testing
├── localhost/           # Local development environment (docker-compose.yaml and hooks)
├── Makefile             # Main orchestrator for commands and tasks
├── .golangci.yml        # Linter configuration for Go
└── go.mod               # Go module dependencies
```

---

## 🛠 Configuración

El servicio se configura primariamente a través de un archivo `app.env` en la raíz del proyecto o mediante **Variables de Entorno**. Las variables de entorno tienen prioridad.

| Variable | Descripción | Ejemplo |
| :--- | :--- | :--- |
| `DATABASE_DSN` | Cadena de conexión a PostgreSQL | `postgres://admin:admin@localhost:5432/goevents?sslmode=disable` |
| `GRPC_SERVER_ADDRESS` | Dirección y puerto del servidor gRPC | `0.0.0.0:9090` |

*(La aplicación se detendrá inmediatamente si no puede cargar la configuración o conectarse a la base de datos)*.

---

## 📡 Interfaz gRPC

La interfaz central está definida en `internal/infrastructure/proto/goevents.proto`. El servicio `Eventservice` provee los siguientes métodos RPC:

**Gestión de Eventos:**
- `CreateEvent`: Registra un evento nuevo (retorna el payload insertado).
- `GetEvent`: Recupera un evento por su ID.
- `DeleteEvent`: Elimina un evento de la base de datos.
- `AllByNameAndSource`: Lista eventos filtrando por nombre y fuente.

**Mensajería (Pub/Sub):**
- `CreateSubscription`: Crea una suscripción asociando un nombre de suscriptor, el nombre de un evento y su fuente.
- `PullMessages`: Extrae mensajes encolados asociados a un evento y fuente dados.
- `AckMessage`: Confirma que un mensaje encolado ha sido procesado exitosamente (Acknowledge).

> Si realizas cambios en el archivo `.proto`, debes regenerar el código ejecutando: `make proto`.

---

## 🚦 Guía de Inicio (Local)

### Requisitos Previos

- [Go](https://golang.org/) 1.26+
- [Docker](https://www.docker.com/) y [Docker Compose](https://docs.docker.com/compose/)
- *Opcional:* Compilador Protobuf (`protoc`) para regenerar código gRPC.

### Pasos para ejecutar

1. **Levantar la Infraestructura Local:**
   Levanta la instancia de PostgreSQL utilizando el entorno local:
   ```bash
   make db-start
   ```

2. **Crear la Base de Datos:**
   Una vez que el contenedor de la BD esté corriendo, crea la base de datos de la aplicación:
   ```bash
   make db-create
   ```

3. **Descargar Dependencias:**
   ```bash
   make tidy
   ```

4. **Ejecutar la Aplicación:**
   ```bash
   make start
   ```
   *Al iniciar, la aplicación ejecutará las automigraciones de GORM (creando/actualizando las tablas `events`, `queue_messages` y `subscriptions`) y escuchará peticiones en el puerto configurado.*

---

## 💻 Comandos de Desarrollo (Makefile)

El proyecto incluye un `Makefile` muy completo para simplificar el ciclo de desarrollo. Ejecuta `make` o `make help` para ver el listado interactivo.

### ⚙️ Aplicación Principal
| Comando | Descripción |
| :--- | :--- |
| `make start` | Inicia la aplicación localmente. |
| `make build` | Genera el binario final de la aplicación en la carpeta `dist`. |
| `make test` | Ejecuta la suite de pruebas de la aplicación. |
| `make proto` | Genera los archivos de código de Go a partir de los archivos `.proto`. |
| `make tidy` | Limpia y actualiza las dependencias de Go (`go mod tidy`). |

### 🗄️ Base de Datos
| Comando | Descripción |
| :--- | :--- |
| `make db-start` / `db-stop` | Inicia o detiene el contenedor de base de datos. |
| `make db-create` / `db-drop` | Crea o elimina por completo (drop) la base de datos. |

### 🔍 Código y Calidad
| Comando | Descripción |
| :--- | :--- |
| `make lint` | Analiza el código Go con `golangci-lint`. |
| `make lint-fix` | Formatea el código automáticamente (`gofmt`, `goimports`). |
| `make support-install-linter` | Instala la herramienta `golangci-lint`. |

### 🛡️ AppSec (Seguridad)
| Comando | Descripción |
| :--- | :--- |
| `make appsec-install` | Instala herramientas de seguridad (Snyk, Gitleaks). |
| `make appsec-test` | Ejecuta pruebas de vulnerabilidades y escaneo de secretos. |
