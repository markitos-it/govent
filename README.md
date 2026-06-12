# Govent

Govent (also known internally as `markitos-it-svc-golden`) is a Go-based gRPC microservice. It provides a structured foundation for managing `Golden` resources using a clean architecture approach.

## 🚀 Technologies

* **Go**: 1.26.4
* **gRPC / Protocol Buffers**: For high-performance RPC communication.
* **PostgreSQL**: Relational database.
* **GORM**: The fantastic ORM library for Golang.
* **Viper**: Complete configuration solution.
* **Docker / Docker Compose**: For local infrastructure management (PostgreSQL).
* **golangci-lint**: Fast linters runner for Go.

## 📁 Project Structure

```text
├── bin/            # Shell scripts for building, testing, and running tasks
├── cmd/            # Main applications for this project
│   └── app/        # Entry point for the service (main.go)
├── internal/       # Private application and library code
│   ├── domain/         # Core business logic and types (entities, repository interfaces)
│   ├── infrastructure/ # External dependencies (gRPC server, database, config)
│   │   ├── configuration/ # Viper config loading
│   │   ├── database/      # GORM Postgres implementation
│   │   ├── gapi/          # gRPC API handlers and generated proto files
│   │   └── proto/         # Protocol Buffer definitions (.proto)
│   └── testsuite/      # Application test suite
├── localhost/      # Local development environment (docker-compose.yaml)
├── Makefile        # Task runner definitions
├── .golangci.yml   # Linter configuration
├── go.mod          # Go module dependencies
└── README.md       # Project documentation
```

## 🛠️ Configuration

The service can be configured via a `config.yaml` file located in the root directory or via environment variables (which take precedence over the config file).

### Environment Variables

| Variable | Description | Example |
|---|---|---|
| `DATABASE_DSN` | PostgreSQL connection string | `postgres://admin:admin@localhost:5432/govent?sslmode=disable` |
| `GRPC_SERVER_ADDRESS` | The address and port for the gRPC server | `0.0.0.0:9090` |

## 🚦 Getting Started

### Prerequisites

* Go 1.26 or higher
* Docker and Docker Compose
* Protocol Buffers Compiler (`protoc`) - Optional, only if you need to regenerate gRPC code.

### 1. Database Setup

Start the local PostgreSQL instance using Docker Compose:

```bash
make db-start
```

Once the database container is running, create the database:

```bash
make db-create
```

*(To stop the database, run `make db-stop`. To drop it, run `make db-drop`)*

### 2. Install Dependencies

Clean up and install Go modules:

```bash
make tidy
```

### 3. Run the Service

You can start the service locally:

```bash
make start
```

The service will automatically run database migrations on startup.

## 💻 Development Commands

The project includes a `Makefile` with helpful commands to speed up your development workflow:

| Command | Description |
|---|---|
| `make help` | Show all available commands |
| `make build` | Build the application binary into the `dist` folder |
| `make start` | Start the application locally |
| `make test` | Run the application test suite |
| `make proto` | Generate Go code from gRPC `.proto` files |
| `make tidy` | Clean and update Go dependencies (`go mod tidy`) |
| `make lint` | Analyze Go code with `golangci-lint` |
| `make lint-fix` | Automatically format Go code (`gofmt`, `goimports`) |

### AppSec (Application Security)

| Command | Description |
|---|---|
| `make appsec-install` | Install security tools (Snyk, Gitleaks) |
| `make appsec-test` | Run security tests |

*(You can also use `make appsec-uninstall`, `make appsec-pre-commit`, and `make appsec-pre-push`)*

## 📡 gRPC Interface

The service definition can be found in `internal/infrastructure/proto/govent.proto`. The `Goldenservice` exposes the following RPC methods:

* `CreateGolden`
* `GetGolden`
* `UpdateGolden`
* `DeleteGolden`
* `ListGoldens`
* `SearchGoldens`

To regenerate the gRPC Go code after modifying the `.proto` file, run:

```bash
make proto
```

## 🧪 Testing

To run the full test suite for the application:

```bash
make test
```