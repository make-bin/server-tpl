# Datastore Layer

This package contains the datastore implementations for the application.

## Supported Datastores

### PostgreSQL
- Full-featured relational database implementation
- Uses GORM as the ORM
- Production-ready with ACID compliance

### OpenGauss
- OpenGauss database implementation
- Compatible with PostgreSQL driver
- Enterprise-grade database solution

### Memory
- In-memory storage implementation
- Useful for testing and development
- Thread-safe with mutex protection
- No persistence - data is lost on restart

## Factory Pattern

The `factory` package provides a simple factory pattern to create datastore instances based on configuration:

```go
factory := factory.NewSimpleFactory()
datastore, err := factory.CreateDatastore(config)
```

## Interface

All datastore implementations must implement the `DatastoreInterface`:

- Application CRUD operations
- Database management (Migrate, Close, HealthCheck)
- Context-aware operations
- Error handling with standardized errors

## Configuration

Configure the datastore type in your application configuration:

```yaml
database:
  type: "postgresql"  # or "opengauss" or "memory"
  host: "localhost"
  port: 5432
  username: "user"
  password: "password"
  database: "myapp"
  sslmode: "disable"
  timezone: "UTC"
```
