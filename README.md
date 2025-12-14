# WiseTech LMS REST API

This is the Go-based REST API foundation for the WiseTech LMS.

## Project Structure

The project follows the standard Go project layout:

- `cmd/api/`: Main application entry point.
- `internal/`: Private application code.
  - `config/`: Configuration loading and management.
  - `database/`: PostgreSQL database connection setup.
  - `models/`: Database model structs.
  - `server/`: HTTP server and routing.
  - `utils/`: Utility functions, including password hashing and validation.
  - `auth/`: Authentication utilities, including JWT token generation and validation.
- `pkg/`: (currently unused) Publicly-usable library code.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.20 or later)

## Setup

1.  **Clone the repository:**
    ```sh
    git clone <your-repository-url>
    cd WiseTech_LMS_Rest_API
    ```

2.  **Set up environment variables:**
    - Copy the example environment file:
      ```sh
      cp .env.example .env
      ```
    - Edit the `.env` file if you want to change the server port or the database file path.
      ```
      # Server Configuration
      SERVER_PORT=8080
      ENVIRONMENT=development
      JWT_SECRET=your-super-secret-key

      # Database Configuration
      DB_PATH=wisetech_lms.db
      ```

3.  **Install dependencies:**
    - The project uses Go modules. Dependencies like `golang.org/x/crypto/bcrypt` and `github.com/golang-jwt/jwt/v5` are automatically downloaded when you build or run the application. You can also install them manually:
      ```sh
      go mod tidy
      ```

## Running the Application

To start the API server, run the following command from the project root:

```sh
go run cmd/api/main.go
```

The server will start on the port specified in your `.env` file (default is `8080`). The first time you run it, a `wisetech_lms.db` file will be created with the necessary tables.

You can check if the server is running by accessing the health check endpoint:

```sh
curl http://localhost:8080/health
```

You should see the following response:

```json
{
  "status": "ok"
}
```
