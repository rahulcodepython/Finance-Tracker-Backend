# Finance Tracker API

This is the backend API for the Finance Tracker application. It is a secure, high-performance backend service designed to power a modern personal finance management application. Built with Golang's Fiber framework and a PostgreSQL database, this service provides a comprehensive RESTful API for all core financial management functionalities.

## Features

- A versioned RESTful API backend service.
- Secure user registration and authentication (Email/Password and Google OAuth).
- Complete CRUD functionality for user profiles, financial accounts, transactions, categories, and budgets.
- Data aggregation endpoints for dashboard visualization and advanced reporting.
- Implementation of security best practices, including JWT authorization, input sanitization, CORS, rate limiting, and password hashing.
- Containerization of the application using Docker for streamlined deployment and development.
- Management of recurring transactions via a background scheduler. The scheduler runs daily to check for and create transactions based on monthly or yearly frequencies.
- Data export functionality (CSV).

## Tech Stack

- **Golang:** The core programming language used for the backend.
- **Fiber:** A fast and expressive web framework for Golang.
- **PostgreSQL:** The chosen relational database for data storage.
- **database/sql:** The standard Go package for interacting with SQL databases.
- **gocron:** A job scheduling package for Go.
- **Docker:** Used for containerization and to ensure a consistent development and deployment environment.

## API Endpoints

All API endpoints are prefixed with `/api/v1`.

### Authentication

- `POST /auth/register`: Register a new user.
- `POST /auth/login`: Log in a user.
- `GET /auth/profile`: Get the authenticated user's profile.
- `POST /auth/change-password`: Change the authenticated user's password.
- `GET /auth/google/login`: Initiate Google OAuth login.
- `GET /auth/google/callback`: Callback for Google OAuth login.

### Accounts

- `POST /accounts/create`: Create a new financial account.
- `GET /accounts`: Get all financial accounts.
- `PATCH /accounts/update/:id`: Update a financial account.
- `DELETE /accounts/delete/:id`: Delete a financial account.
- `GET /accounts/total-balance`: Get the total balance of all active accounts.

### Transactions

- `POST /transactions/create`: Create a new transaction.
- `GET /transactions`: Get all transactions.
- `PATCH /transactions/update/:id`: Update a transaction.
- `DELETE /transactions/delete/:id`: Delete a transaction.
- `GET /transactions/aggregate`: Get aggregate data for transactions.

### Categories

- `POST /categories/create`: Create a new transaction category.
- `GET /categories`: Get all transaction categories.
- `PATCH /categories/update/:id`: Update a transaction category.
- `DELETE /categories/delete/:id`: Delete a transaction category.

### Budgets

- `POST /budgets/create`: Create a new budget.
- `GET /budgets`: Get all budgets.
- `PATCH /budgets/update/:id`: Update a budget.
- `DELETE /budgets/delete/:id`: Delete a budget.

### Recurring Transactions

- `POST /recurring-transactions/create`: Create a new recurring transaction.
- `GET /recurring-transactions`: Get all recurring transactions.
- `PATCH /recurring-transactions/update/:id`: Update a recurring transaction.
- `DELETE /recurring-transactions/delete/:id`: Delete a recurring transaction.

### Dashboard

- `GET /dashboard`: Get a summary of the user's financial data for the dashboard.

### Reports

- `GET /reports`: Generate a financial report.
- `GET /reports/export`: Export transaction data to a CSV file.

## Getting Started

1.  Clone the repository.
2.  Create a `.env` file from the `.env.example` file and fill in the required environment variables.
3.  Run `go run main.go` to start the server.