# Go HTTP Server Template

A modern Go HTTP server template with clean architecture, following best practices and design patterns.

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry
├── pkg/                   # Core application code
│   ├── server/            # Server implementation
│   ├── api/               # HTTP API layer
│   │   ├── dto/v1/        # Data Transfer Objects
│   │   ├── assembler/v1/  # DTO-Model conversion
│   │   ├── router/        # HTTP routing
│   │   ├── middleware/    # HTTP middleware
│   │   ├── application.go # API application layer
│   │   └── interface.go   # API interfaces
│   ├── domain/            # Domain layer
│   │   ├── model/         # Domain models
│   │   └── service/       # Business logic services
│   ├── infrastructure/    # Infrastructure layer
│   │   ├── datastore/     # Data persistence
│   │   └── middleware/    # External service middleware
│   ├── utils/             # Utility packages
│   │   ├── container/     # Dependency injection
│   │   ├── config/        # Configuration management
│   │   ├── logger/        # Logging utilities
│   │   ├── errors/        # Error handling
│   │   └── bcode/         # Business error codes
│   └── e2e/               # End-to-end tests
├── configs/               # Configuration files
├── docs/                  # Documentation
├── deploy/                # Deployment files
├── vendor/                # Vendor dependencies
├── go.mod                 # Go module file
├── go.sum                 # Go dependencies
├── Makefile               # Build automation
├── Dockerfile             # Docker configuration
├── .golangci.yml          # Linter configuration
└── README.md              # Project documentation
```

## Features

- **Clean Architecture**: Layered architecture with clear separation of concerns
- **Dependency Injection**: Simple DI container for managing dependencies
- **Multiple Database Support**: PostgreSQL, OpenGauss, and in-memory storage
- **Middleware Support**: CORS, logging, error handling, recovery, Prometheus metrics
- **Configuration Management**: Environment-based configuration
- **Error Handling**: Structured error handling with business codes
- **Logging**: Structured logging with different levels
- **Health Checks**: Built-in health check endpoints
- **Metrics**: Prometheus metrics integration
- **Redis Support**: Redis client for caching and session management
- **Docker Support**: Complete Docker configuration
- **Code Quality**: Comprehensive linting and code quality checks

## Quick Start

### Prerequisites

- Go 1.21 or later
- PostgreSQL (optional, can use in-memory storage)
- Redis (optional)
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/make-bin/server-tpl.git
cd server-tpl
```

2. Install dependencies:
```bash
make deps
```

3. Build the application:
```bash
make build
```

4. Run the application:
```bash
make run
```

The server will start on port 8080 by default.

### Configuration

Configure the application using environment variables:

```bash
export ENVIRONMENT=development
export PORT=8080
export LOG_LEVEL=info
export DB_TYPE=memory
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=
export DB_DATABASE=server_tpl
export REDIS_ADDRESS=localhost:6379
```

### Docker

Build and run with Docker:

```bash
make docker-build
make docker-run
```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics endpoint
- `GET /api/v1/applications/health` - Application health check

## Development

### Available Make Commands

- `make build` - Build the binary
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make lint` - Run linter
- `make fmt` - Format code
- `make vet` - Vet code
- `make security` - Run security checks
- `make dev-setup` - Setup development environment
- `make ci` - Run CI pipeline

### Architecture Layers

#### API Layer (`pkg/api/`)
- Handles HTTP requests and responses
- Parameter validation and error handling
- Route definition and management
- HTTP middleware implementation

#### Domain Layer (`pkg/domain/`)
- Defines business models and entities
- Implements core business logic
- Defines service interfaces
- Business rules and constraints

#### Infrastructure Layer (`pkg/infrastructure/`)
- Data persistence implementation
- External service integration
- External service middleware

#### Utils Layer (`pkg/utils/`)
- Common utility functions
- Dependency injection container
- Configuration management
- Logging and error handling

## Database Support

The application supports multiple database backends:

- **PostgreSQL**: Production-ready relational database
- **OpenGauss**: Enterprise-grade database solution
- **Memory**: In-memory storage for testing and development

See `pkg/infrastructure/datastore/README.md` for more details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
