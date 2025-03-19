# Restaurant Management System

A backend system for managing restaurant operations, built using Golang and MongoDB.

## Features
- User authentication with JWT
- Menu management (CRUD operations)
- Order processing
- Table management
- Billing system
- Middleware for request validation and authentication

## Tech Stack
- **Golang** – Backend API
- **Gin** – HTTP framework for Go
- **MongoDB** – NoSQL database
- **JWT-Go** – Authentication handling

## Installation

### Prerequisites
- Go installed on your machine
- MongoDB instance running

### Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/YashBubna/restaurant_management.git
   cd restaurant_management
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Set up environment variables in a `.env` file:
   ```env
   MONGO_URI=mongodb://localhost:27017
   JWT_SECRET=your_secret_key
   ```
4. Run the application:
   ```sh
   go run main.go
   ```

## API Endpoints

### 1. User Authentication
**POST** `/register` – Register a new user
```json
{
  "username": "testuser",
  "password": "securepassword"
}
```

**POST** `/login` – Login user and get JWT token
```json
{
  "username": "testuser",
  "password": "securepassword"
}
```
Response:
```json
{
  "token": "your_jwt_token"
}
```

### 2. Menu Management
- **GET** `/menu` – Retrieve all menu items
- **POST** `/menu` – Add a new menu item

### 3. Orders
- **POST** `/order` – Place an order
- **GET** `/order/:id` – Get order details

### 4. Table Management
- **GET** `/tables` – Get available tables
- **POST** `/tables` – Reserve a table

## License
This project is licensed under the MIT License.
