# Inventory Management API in Go

A robust and secure RESTful API for managing inventory, built with Go, Gin, and GORM. This project follows modern Clean Architecture principles for a scalable and maintainable codebase.

## Features

- **Product Management:** Full CRUD (Create, Read, Update, Delete) functionality for products.
- **Secure User Management:** User registration with strong password validation.
- **Professional Authentication:** A complete two-token system using JWTs (short-lived Access Tokens and long-lived Refresh Tokens).
- **Authorization:** Protected API endpoints via custom middleware, ensuring only authenticated users can access sensitive data.
- **Clean Architecture:** A clear separation of concerns using a Handler -> Service -> Repository pattern.

## Technology Stack

- **Language:** [Go](https://golang.org/)
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Authentication:** [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- **Password Hashing:** [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Configuration:** [godotenv](https://github.com/joho/godotenv)

## Getting Started

Follow these instructions to get the project set up and running on your local machine.

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Postman](https://www.postman.com/downloads/) or a similar API client for testing.

### Installation & Setup

1.  **Clone the repository:**

    ```bash
    git clone [https://github.com/RezaBG/Inventory-management-api.git](https://github.com/RezaBG/Inventory-management-api.git)
    cd Inventory-management-api
    ```

2.  **Create the environment file:**
    Create a file named `.env` in the root of the project and paste the following content. Fill in your PostgreSQL details.

    ```env
    # Server Port
    PORT=8080

    # PostgreSQL Database Connection
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_postgres_user
    DB_PASSWORD=your_postgres_password
    DB_NAME=inventory_db

    # JWT Configuration
    JWT_SECRET="your-super-long-and-random-secret-key"
    JWT_ISSUER="inventory-api"
    ACCESS_TOKEN_EXPIRATION_MINUTES=15
    REFRESH_TOKEN_EXPIRATION_HOURS=168
    ```

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

4.  **Run the application:**
    ```bash
    go run cmd/main.go
    ```
    The server will start, connect to the database, run migrations, and be ready to accept requests on the port you specified (e.g., `http://localhost:8080`).

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Path             | Description                                            |
| :----- | :--------------- | :----------------------------------------------------- |
| `POST` | `/register`      | Creates a new user account.                            |
| `POST` | `/login`         | Authenticates a user and returns tokens.               |
| `POST` | `/refresh_token` | Issues a new access token using a valid refresh token. |

**Example: `POST /register`**
Body:

```json
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "P@ssword123!"
}
```

## Protected Endpoints (Authentication Required)

To access these endpoints, you must include an `Authorization` header with a valid Access Token.
**Format:** `Authorization: Bearer <your_access_token>`

| Method   | Path             | Description                           |
| :------- | :--------------- | :------------------------------------ |
| `GET`    | `/products`      | Retrieves a list of all products.     |
| `POST`   | `/products`      | Creates a new product.                |
| `GET`    | `/products/{id}` | Retrieves a single product by its ID. |
| `PUT`    | `/products/{id}` | Updates an existing product.          |
| `DELETE` | `/products/{id}` | Deletes a product.                    |

**Example: `POST /products`**
Header: `Authorization: Bearer eyJhbGciOiJIUzI1Ni...`
Body:

```json
{
  "name": "Gaming Laptop",
  "description": "High-end Laptop with RTX 4090",
  "price": 2499.99,
  "quantity": 10
}
```

## Project Roadmap

- [x] **Phase 1: Foundation** - Project setup and Product CRUD.
- [x] **Phase 2: Core Functionality** - User management and a complete JWT authentication system.
- [ ] **Phase 3: Expanding the Domain** - Adding `Supplier` management and `Inventory Transaction` logic for a full audit trail.
- [ ] **Phase 4: Professional Grade** - API documentation (Swagger), advanced validation, and containerization (Docker).
