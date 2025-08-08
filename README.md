# Inventory Management API in Go

A robust and secure RESTful API for managing inventory, built with Go, Gin, and GORM. This project follows modern Clean Architecture principles for a scalable and maintainable codebase.

## Features

- **Product Management:** Full CRUD functionality for products with real-time, calculated stock quantities.
- **Supplier Management:** Full CRUD functionality for managing suppliers.
- **Inventory Transaction System:** A full audit trail (ledger) for every stock movement (stock-in, stock-out, adjustment).
- **Secure User Management:** User registration with strong password validation and secure `bcrypt` hashing.
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

1. **Clone the repository:**

   ```bash
   git clone [https://github.com/RezaBG/Inventory-management-api.git](https://github.com/RezaBG/Inventory-management-api.git)
   cd Inventory-management-api
   ```

2. **Create the environment file:**
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

3. **Install dependencies:**

   ```bash
   go mod tidy
   ```

4. **Run the application:**

   ```bash
   go run cmd/main.go
   ```

   The server will start, connect to the database, run migrations for all tables, and be ready for requests.

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Path             | Description                                            |
| :----- | :--------------- | :----------------------------------------------------- |
| `POST` | `/register`      | Creates a new user account.                            |
| `POST` | `/login`         | Authenticates a user and returns tokens.               |
| `POST` | `/refresh_token` | Issues a new access token using a valid refresh token. |

---

### Protected Endpoints (Authentication Required)

To access these endpoints, you must include an `Authorization` header with a valid Access Token.
**Format:** `Authorization: Bearer <your_access_token>`

#### Product Endpoints

| Method   | Path             | Description                                                                       |
| :------- | :--------------- | :-------------------------------------------------------------------------------- |
| `GET`    | `/products`      | Retrieves a list of all products with calculated quantities.                      |
| `POST`   | `/products`      | Creates a new product with an initial quantity of 0.                              |
| `GET`    | `/products/{id}` | Retrieves a single product by its ID with calculated quantity.                    |
| `PUT`    | `/products/{id}` | Updates a product's details (name, price, etc.). Quantity cannot be changed here. |
| `DELETE` | `/products/{id}` | Deletes a product.                                                                |

**Example: `POST /products`**
Header: `Authorization: Bearer eyJhbGciOiJIUzI1Ni...`
Body:

```json
{
  "name": "Gaming Laptop",
  "description": "High-end Laptop with RTX 4090",
  "price": 2499.99
}
```

#### Supplier Endpoints

| Method   | Path              | Description                            |
| :------- | :---------------- | :------------------------------------- |
| `GET`    | `/suppliers`      | Retrieves a list of all suppliers.     |
| `POST`   | `/suppliers`      | Creates a new supplier.                |
| `GET`    | `/suppliers/{id}` | Retrieves a single supplier by its ID. |
| `PUT`    | `/suppliers/{id}` | Updates an existing supplier.          |
| `DELETE` | `/suppliers/{id}` | Deletes a supplier.                    |

#### Inventory Endpoints

| Method | Path                      | Description                                           |
| :----- | :------------------------ | :---------------------------------------------------- |
| `POST` | `/inventory/transactions` | Creates a new inventory transaction (e.g., stock-in). |

**Example: `POST /inventory/transactions`**
Header: `Authorization: Bearer eyJhbGciOiJIUzI1Ni...`
Body:

```json
{
  "productID": 1,
  "type": "stock_in",
  "quantityChange": 50,
  "notes": "Received initial shipment."
}
```

## Project Roadmap

- [x] **Phase 1: Foundation** - Project setup and Product CRUD.
- [x] **Phase 2: Core Functionality** - User management and a complete JWT authentication system.
- [x] **Phase 3: Expanding the Domain** - Adding `Supplier` management and `Inventory Transaction` logic for a full audit trail.
- [ ] **Phase 4: Professional Grade** - API documentation (Swagger), advanced validation, and containerization (Docker).
