# User Management API

A RESTful API for user management built with Go, Gin, GORM, PostgreSQL, and Redis.

## Features

- **CRUD Operations**: Complete Create, Read, Update, Delete operations for users
- **Database**: PostgreSQL with GORM ORM
- **Caching**: Redis for improved performance
- **Validation**: Request validation using go-playground/validator
- **Middleware**: Logging, recovery, and CORS support
- **Testing**: Comprehensive unit tests with mocks
- **Docker**: Multi-stage Docker build with health checks
- **Reverse Proxy**: Optional Nginx configuration

## Architecture

```
user_management/
├── config/             # Configuration management
├── controllers/        # HTTP request handlers
├── database/          # Database connection and migrations
├── middleware/        # HTTP middleware (logging, CORS, recovery)
├── models/           # Data models and DTOs
├── repository/       # Data access layer
├── routes/           # Route definitions
├── service/          # Business logic layer
├── tests/            # Unit tests
├── Dockerfile        # Multi-stage Docker build
├── docker-compose.yml # Container orchestration
├── Makefile          # Development commands
└── main.go           # Application entry point
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/users` | Create a new user |
| GET | `/api/v1/users` | Get all users (paginated) |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

## User Model

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "phone": "1234567890",
  "address": "123 Main St",
  "is_active": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

## Quick Start

### Prerequisites

- Go 1.21.5 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd user_management
   ```

2. **Run with Docker Compose**
   ```bash
   make docker-run
   # OR
   docker-compose up --build -d
   ```

3. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Create a user
   curl -X POST http://localhost:8080/api/v1/users \
     -H "Content-Type: application/json" \
     -d '{"name":"John Doe","email":"john@example.com","age":30}'
   ```

### Local Development

1. **Install dependencies**
   ```bash
   make deps
   # OR
   go mod download
   ```

2. **Start PostgreSQL and Redis** (using Docker)
   ```bash
   docker-compose up db redis -d
   ```

3. **Set environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   make run
   # OR
   go run main.go
   ```

## Development

### Available Make Commands

```bash
make help           # Show all available commands
make deps           # Download dependencies
make build          # Build the application
make run            # Run locally
make test           # Run unit tests
make test-coverage  # Run tests with coverage
make lint           # Run linting
make format         # Format code
make docker-build   # Build Docker image
make docker-run     # Run with Docker Compose
make docker-stop    # Stop Docker containers
make clean          # Clean build artifacts
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test file
go test -v ./tests/user_controller_test.go
```

### API Testing Examples

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "age": 28,
    "phone": "0987654321",
    "address": "456 Oak Ave"
  }'

# Get all users (with pagination)
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"

# Get user by ID
curl http://localhost:8080/api/v1/users/1

# Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Updated",
    "email": "jane.updated@example.com",
    "age": 29
  }'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Performance Features

### Caching with Redis
- User data is cached for 15 minutes
- Automatic cache invalidation on updates/deletes
- Graceful fallback when Redis is unavailable

### Database Optimization
- Connection pooling with configurable limits
- Indexed email field for fast lookups
- Soft deletes for data retention

### Nginx Reverse Proxy (Optional)
```bash
# Run with Nginx proxy
make docker-run-proxy
# OR
docker-compose --profile proxy up --build -d
```

Features:
- Rate limiting (10 requests/second with burst of 20)
- Gzip compression
- Security headers
- Load balancing ready

## Configuration

Environment variables can be set in `.env` file or Docker environment:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | Database host |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | password | Database password |
| `DB_NAME` | users_db | Database name |
| `DB_PORT` | 5432 | Database port |
| `SERVER_PORT` | 8080 | Server port |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `GIN_MODE` | debug | Gin mode (debug/release) |

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message",
  "details": "Additional error details (if available)"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `500` - Internal Server Error

## Project Structure Details

### Layered Architecture

1. **Controllers**: Handle HTTP requests and responses
2. **Services**: Contain business logic and validation
3. **Repository**: Data access layer with GORM
4. **Models**: Data structures and validation rules

### Key Design Patterns

- **Dependency Injection**: Services are injected into controllers
- **Interface Segregation**: Each layer depends on interfaces
- **Repository Pattern**: Abstracted data access
- **DTO Pattern**: Separate request/response models

## Testing

The project includes comprehensive unit tests:

- **Controller Tests**: Test HTTP handlers with mocked services
- **Service Tests**: Test business logic with mocked repositories
- **Mock Objects**: Using testify/mock for clean testing

### Test Coverage

Run `make test-coverage` to generate an HTML coverage report.

## Deployment

### Production Considerations

1. **Environment Variables**: Set `GIN_MODE=release`
2. **Database**: Use managed PostgreSQL service
3. **Cache**: Use managed Redis service
4. **Monitoring**: Add health checks and metrics
5. **Security**: Use HTTPS and proper authentication

### Docker Production Build

```bash
make prod-build
docker build -t user-management-api:prod .
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
