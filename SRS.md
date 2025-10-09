# System Requirements Specification (SRS): Finance Tracker

Version: 1.0

Date: October 5, 2025

Status: Draft

---

## 1. Introduction

### 1.1. Project Overview

The Finance Tracker is a secure, high-performance backend service designed to power a modern personal finance management application. Built with Golang's Fiber framework and a PostgreSQL database, this service provides a comprehensive RESTful API for all core financial management functionalities. The system is designed to be containerized using Docker for consistency and scalability, serving a companion Next.js frontend.

### 1.2. Purpose

The primary purpose of this project is to provide users with a robust and secure platform to manage their personal finances effectively. It aims to solve the problem of fragmented financial tracking by offering a centralized system for monitoring accounts, logging transactions, setting budgets, and analyzing spending habits. The API will serve as the single source of truth for the user's financial data, accessible through a client application.

### 1.3. Scope

**In-Scope:**

- A versioned RESTful API backend service.
    
- Secure user registration and authentication (Email/Password and Google OAuth).
    
- Complete CRUD functionality for user profiles, financial accounts, transactions, categories, and budgets.
    
- Data aggregation endpoints for dashboard visualization and advanced reporting.
    
- Implementation of security best practices, including JWT authorization, input sanitization, CORS, rate limiting, and password hashing.
    
- Containerization of the application using Docker for streamlined deployment and development.
    
- Management of recurring transactions via a background scheduler.
    
- Data export functionality (CSV).
    
- Generate a details README.md file containing all essential description, features, API endpoints, and other information about the project.
	

**Out-of-Scope:**

- The development of the client-side frontend application (e.g., Next.js web app, mobile apps).
    
- Direct integration with third-party banking APIs for automatic transaction syncing (Plaid, etc.).
    
- System administration user interface.
    
- Deployment infrastructure setup and management (CI/CD pipelines, cloud hosting).
    
- CSV data import functionality (marked as a future feature).
    

### 1.4. Target Audience

This document is intended for the following stakeholders:

- **Backend Developers:** To understand the system architecture, functional requirements, and API specifications for implementation.
    
- **Frontend Developers:** To understand the API contract, data models, and authentication flow for integrating the client application.
    
- **QA Engineers:** To create test plans and test cases based on the defined requirements.
    
- **Project Managers:** To oversee the project scope, features, and development timeline.
    

## 2. Overall Description

### 2.1. User Personas & Roles

- **Standard User:** The primary user of the application. This user can register, log in, manage their own profile, accounts, transactions, budgets, and view their financial data through dashboards and reports. All data is scoped to their own profile.
    
- **Administrator (System-level):** A potential future role responsible for system maintenance, monitoring application health, and managing system-wide settings. This role does not have access to individual user's financial data but can perform administrative tasks. For the current scope, all API endpoints are designed for the 'Standard User'.
    

### 2.2. Use Case Diagram

This diagram illustrates the main interactions a 'Standard User' has with the Finance Tracker system.

Code snippet

```
graph TD
    actor User as "Standard User"

    subgraph "Finance Tracker System"
        usecase UC1 as "Manage Profile"
        usecase UC2 as "Manage Accounts"
        usecase UC3 as "Manage Transactions"
        usecase UC4 as "Manage Categories (Future)"
        usecase UC5 as "Manage Budgets (Future)"
        usecase UC6 as "View Dashboard"
        usecase UC7 as "Generate Reports"
        usecase UC8 as "Authenticate (Login, Register, Logout)"
    end

    User --|> UC8
    User --|> UC1
    User --|> UC2
    User --|> UC3
    User --|> UC6
    User --|> UC7

    UC8 -.-> UC1 : <<include>>
    UC3 -.-> UC2 : <<include>>

```

## 3. System Architecture & Design

### 3.1. Proposed Technology Stack

- **Backend Language/Framework:** Golang v1.25+ with Fiber v2+
    
- **Database:** PostgreSQL v18+
    
- **Containerization:** Docker & Docker Compose
    
- **Authentication:** JWT (JSON Web Tokens), Google OAuth 2.0
    
- **Password Hashing:** `bcrypt`
    
- **API Specification:** OpenAPI 3.0 (Swagger) for documentation
    
- **Database Migrations:** `golang-migrate/migrate`
    
- **Testing:** Go's built-in testing package with `testify/suite` and `testify/assert`
    
- **Cloud Provider (Suggestion):** AWS (using ECS/Fargate for containers and RDS for PostgreSQL) or any similar provider like GCP or Azure.
    

### 3.2. Folder Structure

A logical and scalable folder structure for the Golang backend service:

Plaintext

```
/finance-tracker-api
├── /api
│   └── /v1              # API versioning, contains handlers/controllers
│       ├── auth_handler.go
│       ├── account_handler.go
│       ├── transaction_handler.go
│       └── dashboard_handler.go
├── /backend
│       ├── /config              # Configuration management
│       │   └── config.go
│       ├── /internal
│       │   ├── /database        # DB connection, setup, migrations, ping
│       │   ├── /middleware      # Custom Fiber middleware (auth, logging)
│       │   ├── /models          # Structs for DB entities and API requests/responses
│       │   ├── /repository      # Data Access Layer (interacts with the DB)
│       │   ├── /routes          # API route definitions
│       │   ├── /services        # Business logic layer
│       │   └── /utils           # Helpers (password, token, validation)
│       ├── /pkg
│           └── /scheduler       # Background jobs (recurring transactions)
├── main.go              # Main application entry point
├── .env.example         # Environment variable template
├── Dockerfile           # Docker build instructions
├── docker-compose.yml   # For local development environment
├── go.mod
└── go.sum
```

### 3.3. Data Flow Diagram (DFD)

A Level 0 DFD showing the context of the system.

Code snippet

```
graph TD
    subgraph "External Entities"
        U[User via Client App]
        G[Google OAuth Service]
    end

    subgraph "System"
        P(Finance Tracker API)
    end

    U -- "API Requests (Login, CRUD ops, Reports)" --> P
    P -- "API Responses (JWT, Data, Reports)" --> U
    P -- "OAuth Redirect" --> G
    G -- "User Profile & Token" --> P
    P -- "Login Confirmation (JWT)" --> U
```

### 3.4. Database Schema

- **Database Type:** PostgreSQL
    
- **Schema Design:** The following tables define the structure for storing user and financial data. The provided SQL script will be used as the basis for database migrations.
    
    ### 🧩 ENUM Types

#### `account_type`

|Value|
|---|
|checking|
|savings|
|credit_card|
|cash|
|investment|
|loan|
|upi|

#### `transaction_type`

|Value|
|---|
|income|
|expense|

#### `auth_provider`

|Value|
|---|
|email|
|google|

---
### 👤 `users` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|full_name|VARCHAR(100)|Not Null|
|email|VARCHAR(255)|Unique, Not Null|
|password_hash|VARCHAR(255)|Nullable|
|provider|auth_provider|Not Null, Default: `'email'`|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

---

### 🏦 `accounts` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|user_id|UUID|Foreign Key → `users(id)`, On Delete CASCADE|
|name|VARCHAR(100)|Not Null|
|type|account_type|Not Null|
|balance|NUMERIC(19,4)|Not Null, Default: `0.00`|
|is_active|BOOLEAN|Not Null, Default: `TRUE`|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|
|updated_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

---

### 🗂️ `categories` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|name|VARCHAR(100)|Not Null|
|type|transaction_type|Not Null|

> 🔒 Unique Constraint: `(name, type)` ensures category names are unique per transaction type.

---

### 💸 `transactions` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|user_id|UUID|Foreign Key → `users(id)`, On Delete CASCADE|
|account_id|UUID|Foreign Key → `accounts(id)`, On Delete CASCADE|
|category_id|UUID|Foreign Key → `categories(id)`, On Delete RESTRICT|
|description|VARCHAR(255)|Not Null|
|amount|NUMERIC(19,4)|Not Null|
|type|transaction_type|Not Null|
|transaction_date|DATE|Not Null|
|note|TEXT|Optional|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|
|updated_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

---

### ⚡ Indexes

| Index Name                    | Columns                            |
| ----------------------------- | ---------------------------------- |
| idx_transactions_user_id_date | `(user_id, transaction_date DESC)` |
| idx_accounts_user_id          | `(user_id)`                        |

---

## 4. Functional Requirements

#### Module 1: User Authentication & Profile

- **FR-01: User Registration:** The system shall allow a new user to register using their Full Name, Email, and a secure Password.
    
- **FR-02: User Login:** The system shall allow a registered user to log in using their Email and Password.
    
- **FR-03: Google OAuth Integration:** The system shall allow users to sign up or sign in using their Google account.
    
- **FR-04: JWT Generation:** Upon successful authentication, the system shall generate and return a secure JSON Web Token (JWT) to the user.
    
- **FR-05: Route Protection:** The system shall protect sensitive API endpoints, allowing access only to requests with a valid JWT.
    
- **FR-06: View Profile:** The system shall allow an authenticated user to view their profile information (Full Name, Email).
    
- **FR-07: Update Profile:** The system shall allow an authenticated user to update their Full Name.
    
- **FR-08: Change Password:** The system shall provide a secure endpoint for an authenticated user to change their password.
    
- **FR-09: Fetch User Statistics:** The system shall provide an endpoint to return user statistics, including total number of accounts, total transactions, and days since registration.
    

#### Module 2: Account Management

- **FR-10: Create Account:** The system shall allow an authenticated user to create a new financial account (e.g., Savings, Checking) by providing a name and type.
    
- **FR-11: Read Accounts:** The system shall allow a user to retrieve a list of all their financial accounts, including their current balance, type, and status.
    
- **FR-12: Update Account:** The system shall allow a user to update the details of an existing account (name, type, status).
    
- **FR-13: Delete Account:** The system shall allow a user to delete an account. Deletion should be blocked if the account has associated transactions.
    
- **FR-14: Calculate Total Balance:** The system shall provide an endpoint to calculate and return the sum of balances from all the user's active accounts.
    

#### Module 3: Transaction Management

- **FR-15: Create Transaction:** The system shall allow a user to add a new transaction (income or expense), linking it to a specific account and category.
    
- **FR-16: Read Transactions:** The system shall allow a user to retrieve a list of all their transactions.
    
- **FR-17: Paginate Transactions:** The transaction list endpoint must support server-side pagination (e.g., page number, page size).
    
- **FR-18: Filter/Search Transactions:** The transaction list endpoint must allow filtering by description, category, account, and date range.
    
- **FR-19: Update Transaction:** The system shall allow a user to update the details of an existing transaction.
    
- **FR-20: Delete Transaction:** The system shall allow a user to delete a transaction. This action must also update the corresponding account balance.
    
- **FR-21: View Aggregate Data:** The system shall provide endpoints to view total income, total expenses, and net income over a specified period.
    

#### Module 4 & 5: Dashboard & Reporting

- **FR-22: Dashboard Summary:** The system shall provide a single endpoint to fetch consolidated data for a dashboard view, including total balance, current month's income/expenses/savings.
    
- **FR-23: Recent Transactions:** The system shall provide an endpoint to fetch the 10 most recent transactions for the user.
    
- **FR-24: Historical Comparison Data:** The system shall provide data formatted for an "Income vs. Expenses" graph for the last 12 months.
    
- **FR-25: Categorical Spending Data:** The system shall provide data formatted for a "Monthly Spending by Category" pie chart for the current month.
    
- **FR-26: Custom Date-Range Reports:** The system shall provide a powerful reporting endpoint that accepts a start and end date to generate detailed financial summaries.
    
- **FR-27: Trend Analysis Data:** The reporting endpoint shall return data suitable for visualizing income/expense trends over the selected period.
    

#### Module 6, 7, 8, 9: Advanced Features

- **FR-28: Manage Budgets:** The system shall provide full CRUD functionality for users to set monthly budgets for specific expense categories.
    
- **FR-29: Manage Custom Categories:** The system shall provide full CRUD functionality for users to create, update, and delete their own custom transaction categories.
    
- **FR-30: Default Categories:** Upon registration, the system shall populate a user's account with a predefined set of default income and expense categories.
    
- **FR-31: Manage Recurring Transactions:** The system shall allow users to set up, view, update, and delete recurring transactions with a defined frequency.
    
- **FR-32: Process Recurring Transactions:** A background worker shall automatically create new transactions when a recurring transaction is due.
    
- **FR-33: Export Data:** The system shall allow users to export their transaction data within a specified date range to a CSV file.
    

## 5. Non-Functional Requirements

### 5.1. Performance

- **NFR-01 (Response Time):** 95% of all API requests must be processed and responded to in under 250ms.
    
- **NFR-02 (Complex Queries):** Reporting and dashboard endpoints with complex aggregations must respond in under 800ms.
    
- **NFR-03 (Database Performance):** Proper indexing must be applied to all frequently queried columns (e.g., foreign keys, dates) to ensure fast query execution.
    

### 5.2. Scalability

- **NFR-04 (Concurrent Users):** The system architecture must be designed to support at least 1,000 concurrent users during its initial launch phase, with the ability to scale horizontally.
    
- **NFR-05 (Statelessness):** The API must be stateless, allowing for easy horizontal scaling by adding more container instances behind a load balancer.
    

### 5.3. Security

- **NFR-06 (Data Encryption):** All data in transit must be encrypted using TLS 1.2 or higher.
    
- **NFR-07 (Password Security):** User passwords must be hashed using a strong, adaptive algorithm like `bcrypt`. Passwords must never be stored in plaintext.
    
- **NFR-08 (Authentication):** All endpoints, except for registration and login, must be protected and require a valid JWT.
    
- **NFR-09 (SQL Injection):** All database queries must be executed using parameterized statements or an ORM that provides this protection by default to prevent SQL injection attacks.
    
- **NFR-10 (XSS Protection):** The API must handle user-generated content by properly sanitizing or encoding output to prevent Cross-Site Scripting (XSS) vulnerabilities on the client side.
    
- **NFR-11 (CORS Policy):** A strict Cross-Origin Resource Sharing (CORS) policy must be implemented to allow requests only from the authorized frontend domain.
    
- **NFR-12 (Rate Limiting):** Sensitive endpoints (e.g., login, password reset requests) must be protected by variable rate limiting to mitigate brute-force attacks.
    

### 5.4. Reliability

- **NFR-13 (Uptime):** The system must achieve a minimum of 99.9% uptime, excluding planned maintenance windows.
    
- **NFR-14 (Data Integrity):** The system must use database transactions for operations that involve multiple writes (e.g., creating a transaction and updating an account balance) to ensure atomicity and data integrity.
    
- **NFR-15 (Logging):** The application must generate structured logs for key events, errors, and access patterns to facilitate monitoring and debugging.
    

### 5.5. Usability

- **NFR-16 (API Documentation):** A clear and comprehensive API documentation (e.g., using OpenAPI/Swagger) must be maintained.
    
- **NFR-17 (Error Handling):** The API must provide consistent and meaningful error messages and HTTP status codes to aid frontend development and debugging.
    

## 6. API Specification (Revised)

All API responses will adhere to the following structure:

JSON

```
{
  "success": true, // or false
  "message": "Descriptive message",
  "data": { ... }, // or null
  "error": "Error details if success is false" // or null
}
```

---

### **`/api/v1/auth`**

- **Endpoint: `POST /login`**
    
    - **Description:** Authenticates a user and returns a JWT.
        
    - **Authorization:** Public
        
    - **Request Body:**
        
        JSON
        
        ```
        {
          "email": "john.doe@example.com",
          "password": "aVeryStrongPassword123!"
        }
        ```
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Login successful",
          "data": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expiresAt": "2025-10-05T10:30:00Z"
          },
          "error": null
        }
        ```
        
- **Endpoint: `POST /register`**
    
    - **Description:** Registers a new user and returns a JWT.
        
    - **Authorization:** Public
        
    - **Request Body:**
        
        JSON
        
        ```
        {
          "fullName": "Jane Doe",
          "email": "jane.doe@example.com",
          "password": "anotherSecurePassword456!"
        }
        ```
        
    - **Success Response (201 Created):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "User registered successfully",
          "data": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expiresAt": "2025-10-05T10:31:00Z"
          },
          "error": null
        }
        ```
        
- **Endpoint: `GET /profile`**
    
    - **Description:** Retrieves the profile and stats for the authenticated user.
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Profile retrieved successfully",
          "data": {
            "personal": {
              "fullName": "Jane Doe",
              "email": "jane.doe@example.com"
            },
            "stats": {
              "activeAccounts": 3,
              "totalTransactions": 152,
              "daysActive": 90
            }
          },
          "error": null
        }
        ```
        
- **Endpoint: `POST /change-password`**
    
    - **Description:** Allows an authenticated user to change their password.
        
    - **Authorization:** Authenticated User
        
    - **Request Body:**
        
        JSON
        
        ```
        {
          "currentPassword": "anotherSecurePassword456!",
          "newPassword": "aNewerEvenStrongerPassword789!"
        }
        ```
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Password changed successfully",
          "data": null,
          "error": null
        }
        ```
        
- **Endpoint: `GET /logout`**
    
    - **Description:** Logs out the user. (Note: For stateless JWT, this is typically handled client-side by deleting the token. A server-side implementation would involve a token blocklist.)
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Logged out successfully",
          "data": null,
          "error": null
        }
        ```
        

---

### **`/api/v1/dashboard`**

- **Endpoint: `GET /`**
    
    - **Description:** Retrieves aggregated data for the main dashboard view for the current month.
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Dashboard data retrieved successfully",
          "data": {
            "summary": {
              "totalBalance": "12560.75",
              "monthlyIncome": "6000.00",
              "monthlyExpenses": "-2150.25",
              "monthlySavings": "3849.75"
            },
            "graphs": {
              "incomeVsExpense": [
                { "month": "Aug", "income": 5800, "expense": 2200 },
                { "month": "Sep", "income": 6000, "expense": 2500 },
                { "month": "Oct", "income": 6000, "expense": 2150.25 }
              ],
              "spendingByCategory": [
                { "category": "shopping", "amount": 850.50 },
                { "category": "food", "amount": 650.00 },
                { "category": "transport", "amount": 350.75 }
              ]
            },
            "recentTransactions": [
              { "description": "Weekly Groceries", "date": "2025-10-03", "type": "expense", "amount": "120.50" },
              { "description": "Client Payment", "date": "2025-10-01", "type": "income", "amount": "1500.00" }
            ]
          },
          "error": null
        }
        ```
        

---

### **`/api/v1/accounts`**

- **Endpoint: `GET /`**
    
    - **Description:** Retrieves all of the user's financial accounts.
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Accounts retrieved successfully",
          "data": {
            "totalBalance": "12560.75",
            "accounts": [
              { "id": "uuid-1", "name": "Main Checking", "isActive": true, "balance": "5430.50", "type": "checking" },
              { "id": "uuid-2", "name": "Savings Account", "isActive": true, "balance": "7130.25", "type": "savings" }
            ]
          },
          "error": null
        }
        ```
        
- **Endpoints: `POST /create`, `PATCH /update/:id`, `DELETE /delete/:id`**
    
    - These endpoints will perform standard CRUD operations on accounts, accepting and returning account objects in the `data` field of the response.
        

---

### **`/api/v1/transactions`**

- **Endpoint: `GET /?page=1&limit=20&category=food&description=coffee`**
    
    - **Description:** Retrieves a paginated and filtered list of transactions.
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Transactions retrieved successfully",
          "data": {
            "stats": {
              "totalIncome": "6000.00",
              "totalExpenses": "-2150.25",
              "netIncome": "3849.75"
            },
            "accounts": [
              { "id": "uuid-1", "name": "Main Checking" },
              { "id": "uuid-2", "name": "Savings Account" }
            ],
            "transactions": [
              { "id": "txn-uuid-1", "type": "expense", "amount": "4.50", "description": "Morning coffee", "category": "food", "accountId": "uuid-1", "date": "2025-10-04", "note": "From Starbucks" }
            ],
            "pagination": { "currentPage": 1, "totalPages": 5 }
          },
          "error": null
        }
        ```
        
- **Endpoints: `POST /create`, `PATCH /update/:id`, `DELETE /delete/:id`**
    
    - These endpoints will perform standard CRUD operations on transactions, accepting and returning transaction objects in the `data` field of the response. The `category` and `type` fields must correspond to the predefined values.
        

---

### **`/api/v1/reports`**

- **Endpoint: `GET /?from=2025-09-01&to=2025-09-30`**
    
    - **Description:** Generates a detailed report for a given date range. Defaults to the last month if no params are provided.
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        JSON
        
        ```
        {
          "success": true,
          "message": "Report generated successfully",
          "data": {
            "summary": {
              "totalIncome": "6000.00",
              "totalExpenses": "-2500.00"
            },
            "graphs": {
              "incomeVsExpense": [
                { "date": "2025-09-01", "income": 1500, "expense": 50 },
                { "date": "2025-09-15", "income": 1500, "expense": 800 }
              ],
              "spendingByCategory": [
                { "category": "shopping", "amount": 1200 },
                { "category": "food", "amount": 800 }
              ]
            },
            "incomeSources": [
              { "type": "salary", "name": "Salary", "percentage": 83.33 },
              { "type": "freelance", "name": "Freelance", "percentage": 16.67 }
            ],
            "spendingBreakdown": [
              { "category": "shopping", "amount": "1200.00", "percentage": 48.00 },
              { "category": "food", "amount": "800.00", "percentage": 32.00 }
            ]
          },
          "error": null
        }
        ```

## 7. Authentication & Authorization

### 7.1. Authentication Strategy

The system employs a JWT-based authentication strategy for its stateless API.

**Flow:**

1. A user submits credentials (email/password) or an OAuth token (from Google).
    
2. The server validates the credentials/token.
    
3. If valid, the server generates a signed JWT containing a payload with the `user_id` and an expiration timestamp (`exp`).
    
4. The server sends this JWT back to the client.
    
5. The client application stores the JWT securely (e.g., in an HttpOnly cookie or secure storage).
    
6. For all subsequent requests to protected endpoints, the client must include the JWT in the `Authorization` header with the `Bearer` scheme (`Authorization: Bearer <token>`).
    
7. A middleware on the server intercepts each request, validates the JWT's signature and expiration, and if valid, extracts the `user_id` to process the request within the user's scope.
    
8. Tokens will have a short lifespan (e.g., 1 hour), and a refresh token mechanism will be implemented for a seamless user experience.
    

### 7.2. Authorization Strategy

Authorization is managed via Role-Based Access Control (RBAC). For the current scope, the roles are simple.

**Roles:**

- **Public:** Any unauthenticated user.
    
- **Authenticated User:** Any user who has successfully logged in and possesses a valid JWT.
    

Permissions:

| Endpoint Group          | Required Role      | Description                                    |
| ----------------------- | ------------------ | ---------------------------------------------- |
| /api/v1/auth/*          | Public             | User registration and login.                   |
| /api/v1/users/me        | Authenticated User | All profile management operations.             |
| /api/v1/accounts/**     | Authenticated User | Full CRUD on the user's own accounts.          |
| /api/v1/transactions/** | Authenticated User | Full CRUD on the user's own transactions.      |
| /api/v1/categories/**   | Authenticated User | Full CRUD on the user's own categories.        |
| /api/v1/budgets/**      | Authenticated User | Full CRUD on the user's own budgets.           |
| /api/v1/dashboard/**    | Authenticated User | Access to dashboard and summary data.          |
| /api/v1/reports/**      | Authenticated User | Access to generate personal financial reports. |

## 8. Environment Configuration

The following `.env.example` file provides a template for all necessary configuration variables.

```
# -------------------------------------
# Application Configuration
# -------------------------------------
# The port the application will run on (e.g., 8080)
PORT=8080
# Application environment ('development', 'production', 'staging')
APP_ENV=development
# The allowed origin for CORS policy (e.g., http://localhost:3000)
CLIENT_ORIGIN=

# -------------------------------------
# Database Configuration (PostgreSQL)
# -------------------------------------
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=finance_tracker
DB_SSL_MODE=disable # Use 'require' in production

# -------------------------------------
# Security Configuration
# -------------------------------------
# A very strong secret key for signing JWTs
JWT_SECRET=a_super_secret_string_32_chars_long
# JWT token expiration time (e.g., 1h, 15m, 7d)
JWT_EXPIRES_IN=1h
# JWT refresh token expiration time
JWT_REFRESH_EXPIRES_IN=7d

# -------------------------------------
# External Services (Google OAuth)
# -------------------------------------
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

# -------------------------------------
# Rate Limiting Configuration
# -------------------------------------
# Max requests per minute for general endpoints
RATE_LIMITER_MAX=100
# Duration of the rate limit window in minutes
RATE_LIMITER_DURATION_MINUTES=1
```