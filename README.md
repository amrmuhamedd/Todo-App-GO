# Todo API

A RESTful API for managing todos with user authentication built using Go and modern web technologies.

## Technology Stack

- **Go** - Main programming language
- **Gin** - Web framework
- **GORM** - ORM for database operations
- **JWT** - Authentication mechanism
- **PostgreSQL** - Database
- **Swagger** - API documentation
- **bcrypt** - Password hashing

## Project Architecture

```
├── cmd/                  # Application entry points
├── internal/            # Private application code
│   ├── auth/           # Authentication logic
│   ├── config/         # Configuration management
│   ├── database/       # Database connections and migrations
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware components
│   ├── models/         # Database models
│   └── routes/         # Route definitions
├── docs/               # Swagger documentation
└── main.go            # Main application entry point
```

### Key Components

1. **Authentication**
   - JWT-based authentication
   - Secure password hashing with bcrypt
   - Protected routes with middleware

2. **Database**
   - GORM for database operations
   - PostgreSQL for data persistence
   - Models for Users and Todos

3. **API Endpoints**
   - `/api/auth/signup` - User registration
   - `/api/auth/login` - User authentication
   - `/api/todos` - Todo CRUD operations
   - Protected routes with JWT middleware

4. **Documentation**
   - Swagger UI for API documentation
   - Auto-generated API specs

## Getting Started

### Prerequisites

1. Go 1.19 or later
2. PostgreSQL

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd todo-api
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables (create a .env file)
```env
JWT_SECRET=your_jwt_secret_key
DATABASE_URL=postgresql://username:password@localhost:5432/dbname?sslmode=disable
```

4. Run the application
```bash
go run main.go
```

The server will start at `http://localhost:8080`

### API Documentation

Access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Development

### Generate Swagger Documentation
```bash
swag init
```

### Run Tests
```bash
go test ./...
```

### Build for Production
```bash
go build -o todo-api
```

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Register a new user
- `POST /api/auth/login` - Login and receive JWT token

### Todos
- `GET /api/todos` - List all todos
- `POST /api/todos` - Create a new todo
- `GET /api/todos/:id` - Get a specific todo
- `PUT /api/todos/:id` - Update a todo
- `DELETE /api/todos/:id` - Delete a todo

## Security

- All passwords are hashed using bcrypt
- JWT tokens are required for protected endpoints
- Environment variables for sensitive data
- Input validation on all endpoints
