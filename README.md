# go-events

**go-events** es un microservicio backend desarrollado en Go que implementa un sistema de gestión y enrutamiento de eventos mediante **gRPC**. Está diseñado siguiendo los principios de la **Arquitectura Limpia (Clean Architecture)** para ofrecer un núcleo de negocio altamente testable, desacoplado y escalable.

El servicio expone operaciones tanto para la gestión tradicional de recursos (`Event`), como para patrones de mensajería del tipo Pub/Sub (`Subscription`, `PullMessages`, `AckMessage`).

---

## 🚀 Tecnologías Principales

- **Lenguaje:** Go 1.26.4
- **Comunicaciones:** gRPC y Protocol Buffers (`protoc`)
- **Base de Datos:** PostgreSQL
- **ORM:** GORM
- **Configuración:** Viper
- **Logging:** `slog` (con formato coloreado estructurado en el main)
- **Infraestructura Local:** Docker & Docker Compose
- **Calidad y Seguridad:** `golangci-lint`, Snyk, Gitleaks

---

## 🏗 Arquitectura y Estructura del Proyecto

El proyecto sigue una distribución orientada a dominios:

```text
├── bin/                 # Scripts de soporte (.sh) para el Makefile
├── cmd/
│   └── app/             # Punto de entrada de la aplicación (main.go)
├── internal/
│   ├── domain/          # Entidades centrales (Event, Queue, Subscription) y puertos (interfaces)
│   ├── infrastructure/  # Adaptadores de salida y entrada (GORM, Viper, Logging, gRPC server)
│   │   ├── configuration/ # Carga de configuración (app.env / variables de entorno)
│   │   ├── database/      # Implementación del repositorio (Postgres)
│   │   ├── gapi/          # Handlers gRPC y código autogenerado (.pb.go)
│   │   └── proto/         # Definición de la interfaz gRPC (.proto)
│   └── testsuite/       # Tests de integración y pruebas de infraestructura
├── localhost/           # Entorno de desarrollo local (docker-compose.yaml y hooks)
├── Makefile             # Orquestador principal de comandos y tareas
├── .golangci.yml        # Configuración del linter para Go
└── go.mod               # Dependencias del módulo Go
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
