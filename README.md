<div align="center">
  <h1>🚀 go-events</h1>
  <p><strong>A robust, scalable event management and routing microservice built with Go and gRPC.</strong></p>

  <p>
    <img src="https://img.shields.io/badge/Go-1.26.4-00ADD8?style=flat-square&logo=go" alt="Go Version" />
    <img src="https://img.shields.io/badge/gRPC-Framework-244C5A?style=flat-square&logo=grpc" alt="gRPC" />
    <img src="https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat-square&logo=postgresql" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/Architecture-Clean-success?style=flat-square" alt="Clean Architecture" />
  </p>
</div>

---

## 📖 Overview

**go-events** is a backend microservice developed in Go that implements an event management and routing system via **gRPC**. It is meticulously designed following **Clean Architecture** principles to provide a highly testable, decoupled, and scalable business core.

The service exposes operations for both traditional resource management (`Event`) and Pub/Sub messaging patterns (`Subscription`, `PullMessages`, `AckMessage`).

---

## ✨ Key Features

- **Event Storage & Retrieval:** Store, fetch, and delete individual events, or list them by slug and source.
- **Pub/Sub Messaging System:** Create subscriptions for specific events and pull queued messages reliably.
- **Acknowledge Mechanism:** Safely acknowledge (`Ack`) processed messages to ensure they are handled properly by subscribers.
- **Clean Architecture:** Strict separation of concerns (Domain, Infrastructure, Application).
- **Rich Observability:** Implements structured and colored logging using `slog` and unary gRPC interceptors for request/response tracing.
- **Developer Experience:** Fully containerized local environment and a powerful `Makefile` for automated workflows.

---

## 🛠️ Tech Stack

- **Language:** Go `1.26.4`
- **Communication:** gRPC & Protocol Buffers (`protoc`)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Configuration Management:** Viper
- **Logging:** Go's standard `log/slog` (with custom structured & colored formatting)
- **Local Infrastructure:** Docker & Docker Compose
- **Code Quality & Security:** `golangci-lint`, Snyk, Gitleaks

---

## 🏗️ Architecture & Project Structure

The project follows a strictly domain-oriented directory layout to enforce Clean Architecture boundaries:

```text
├── bin/                 # Shell scripts for Makefile automation
├── cmd/
│   └── app/             # Application entry point (main.go)
├── internal/
│   ├── domain/          # Core business entities (Event, Queue, Subscription) and Ports (Interfaces)
│   ├── infrastructure/  # Input/Output Adapters (GORM, Viper, Logging, gRPC server)
│   │   ├── configuration/ # Configuration loader (config.yaml / Environment Variables)
│   │   ├── database/      # PostgreSQL Repository implementations
│   │   ├── gapi/          # gRPC Handlers and auto-generated code (.pb.go)
│   │   └── proto/         # gRPC interface definitions (.proto)
│   └── testsuite/       # Integration tests and infrastructure testing
├── localhost/           # Local development environment (docker-compose.yaml & git hooks)
├── Makefile             # Main orchestrator for commands and tasks
├── .golangci.yml        # Go linter configuration
└── go.mod               # Go module dependencies
```

---

## ⚙️ Configuration

The service is configured primarily through a `config.yaml` file located in the root directory, or via **Environment Variables**. Environment variables take precedence over the configuration file.

| Environment Variable | Description | Example |
| :--- | :--- | :--- |
| `DATABASE_DSN` | PostgreSQL connection string | `postgres://admin:admin@localhost:5432/goevents?sslmode=disable` |
| `GRPC_SERVER_ADDRESS` | Address and port for the gRPC server | `0.0.0.0:9090` |

> ⚠️ *The application will terminate immediately upon startup if it cannot load the configuration or connect to the database.*

---

## 📡 gRPC Interface

The core interface is defined in `internal/infrastructure/proto/govent.proto`. The `Eventservice` provides the following RPC methods:

### 📦 Event Management
- `CreateEvent`: Registers a new event and returns the inserted payload.
- `GetEvent`: Retrieves a specific event by its ID.
- `DeleteEvent`: Removes an event from the database.
- `AllBySlugAndSource`: Lists events filtered by their slug and source.

### 📨 Messaging (Pub/Sub)
- `CreateSubscription`: Creates a subscription associating a subscriber name with an event name and its source.
- `PullMessages`: Extracts queued messages associated with a given event and source.
- `AckMessage`: Confirms that a queued message has been processed successfully (Acknowledge).

> 💡 **Note:** If you make changes to the `.proto` file, you must regenerate the Go code by running: `make proto`.

---

## 🚦 Getting Started (Local Development)

### Prerequisites

- [Go](https://golang.org/) 1.26+
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- *(Optional)* Protobuf Compiler (`protoc`) to regenerate gRPC code.

### Running the Application

1. **Spin up the Local Infrastructure:**
   Start the PostgreSQL container using the local environment setup:
   ```bash
   make db-start
   ```

2. **Create the Database:**
   Once the DB container is running, create the application database:
   ```bash
   make db-create
   ```

3. **Install Dependencies:**
   ```bash
   make tidy
   ```

4. **Start the Application:**
   ```bash
   make start
   ```
   *Upon startup, the application will automatically run GORM migrations (creating/updating the `events`, `queues`, and `subscriptions` tables) and begin listening for gRPC requests on the configured port.*

---

## 💻 Development Commands (Makefile)

The project includes a comprehensive `Makefile` to simplify the development lifecycle. Run `make` or `make help` to see the interactive list of commands.

### ⚙️ Core Application
| Command | Description |
| :--- | :--- |
| `make start` | Starts the application locally. |
| `make build` | Genera the final application binary in the `dist` folder. |
| `make test` | Executes the application's test suite. |
| `make proto` | Generates Go code from the `.proto` files. |
| `make tidy` | Cleans and updates Go module dependencies (`go mod tidy`). |

### 🗄️ Database Management
| Command | Description |
| :--- | :--- |
| `make db-start` | Starts the local database container. |
| `make db-stop` | Stops the local database container. |
| `make db-create` | Creates the necessary database for the application. |
| `make db-drop` | Completely drops the application database. |

### 🔍 Code Quality & Formatting
| Command | Description |
| :--- | :--- |
| `make lint` | Analyzes Go code using `golangci-lint`. |
| `make lint-fix` | Automatically formats the code (`gofmt`, `goimports`). |
| `make support-install-linter` | Installs the `golangci-lint` tool locally. |

### 🛡️ AppSec (Security)
| Command | Description |
| :--- | :--- |
| `make appsec-install` | Installs security tools (Snyk, Gitleaks). |
| `make appsec-test` | Runs vulnerability tests and secret scanning. |

---

<div align="center">
  <i>Built with ❤️ using Go & gRPC</i>
</div>