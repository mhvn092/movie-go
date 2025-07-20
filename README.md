# Movie-Go: A Simple Movie Database in Go

## Overview
Movie-Go is a minimalistic movie database application built using raw Go syntax with few external libraries. Designed for learning Go, it demonstrates core programming concepts, Domain-Driven Design (DDD), and web development patterns. The project has dynamic data entities within the codebase with SQL relationships rather than using static data.

## Features
- **Dynamic Data Entities**: Data entities that are defined in the project itself and their relationships are managed via SQL foreign keys rather than using static data.
- **Custom Libraries**:
  - **Router**: A custom-built HTTP router for handling requests.
  - **Env Reader**: A lightweight package for reading `.env` files.
  - **Payload Validator**: Custom validation for request payloads.
  - **Migration CLI Tool**: A command-line tool for database migrations.
- **Domain-Driven Design (DDD)**: Structured for clear separation of concerns.
- **Authentication**:
  - JWT-based authentication with middleware.
  - Role-based authorization (Admin and Non-Admin roles).
- **Middleware**:
  - **Logger Middleware**: Logs HTTP requests/responses in color-formatted JSON.
  - **Panic Middleware**: Handles panics to prevent server crashes.
- **Database**: PostgreSQL with the `pgx` driver for efficient interactions.

## Purpose
This project is a learning tool for mastering Go syntax, DDD, and building RESTful APIs with minimal dependencies, emphasizing custom utility development.

## Tech Stack
- **Language**: Go 1.22
- **Database**: PostgreSQL (`github.com/jackc/pgx/v5`)
- **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Cryptography**: `golang.org/x/crypto` for secure hashing
- **Architecture**: DDD
- **Custom Utilities**: Router, Env Reader, Payload Validator, Migration CLI

## Prerequisites
- Go 1.22+
- PostgreSQL
- Make (for Makefile commands)

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/mhvn092/movie-go.git
   cd movie-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   Create a `.env` file:
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/movie_db
   JWT_SECRET=your_jwt_secret_key
   ```

4. Run migrations:
   ```bash
   make up
   ```

5. Build and run:
   ```bash
   make run
   ```

## Makefile Commands
- `make build`: Compiles to `/tmp/bin/movie.exe`.
- `make run`: Builds and runs the application.
- `make migrate`: Compiles the migration CLI to `/tmp/bin/migration`.
- `make create NAME="migration_name"`: Creates a new migration file.
- `make up`: Applies pending migrations.
- `make down`: Rolls back the latest migration.

## Project Structure
```
movie-go/
├── cmd/
│   ├── service/         # Main application entry point
│   └── migration/       # Database migration CLI tool
├── internal/            # DDD-based packages (domain, platform, transport)
├── migrations/          # Database migration files
├── pkg/                 # Custom packages written by my own 
├── .env                # Environment variables
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
└── Makefile            # Build and migration commands
```

## Dependencies
- `github.com/jackc/pgx/v5`: PostgreSQL driver
- `github.com/golang-jwt/jwt/v5`: JWT authentication
- `golang.org/x/crypto`: Cryptographic utilities
- Indirect dependencies for testing and database utilities (see `go.mod`)

## Usage
- **API Endpoints**: RESTful endpoints for managing movies, users, and authentication. (API docs TBD.)
- **Authentication**: JWT tokens secure protected routes; admin routes require admin role.
- **Logging**: Requests/responses logged in color-formatted JSON.

## Contributing
Contributions are welcome! Submit issues or pull requests to [GitHub](https://github.com/mhvn092/movie-go).

## License
MIT License
