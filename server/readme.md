# eShop Backend

This is the backend service for the eShop application, built with Go (Golang).

## Features

- User authentication and authorization
- Product management (CRUD operations)
- Order processing
- Payment integration
- RESTful API

## Requirements

- Go 1.18 or higher
- PostgreSQL
- Redis

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/eshop-backend.git
    cd eshop-backend
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up environment variables:
    ```sh
    cp .env.example .env
    # Update .env with your configuration
    ```

4. Run the application:
    ```sh
    go run main.go
    ```

## API Endpoints

- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - User login
- `GET /api/v1/products` - Get all products
- `POST /api/v1/products` - Create a new product
- `PUT /api/v1/products/{id}` - Update a product
- `DELETE /api/v1/products/{id}` - Delete a product
- `POST /api/v1/orders` - Create a new order

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.