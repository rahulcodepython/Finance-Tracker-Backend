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
│       ├── auth.handler.go
│       ├── account.handler.go
│       ├── transaction.handler.go
│       ├── dashboard.handler.go
│       ├── budget.handler.go
│       ├── category.handler.go
│       ├── recurring_transaction.handler.go
│       └── report.handler.go
├── /backend
│       ├── /config              # Configuration management
│       │   └── config.go
│       ├── /database        # DB connection, setup, migrations, ping
│       │   └── database.go
│       ├── /middleware      # Custom Fiber middleware (auth, logging)
│       │   └── auth.go
│       ├── /models          # Structs for DB entities and API requests/responses
│       │   ├── account.go
│       │   ├── budget.go
│       │   ├── category.go
│       │   ├── recurring_transaction.go
│       │   ├── transaction.go
│       │   └── user.go
│       ├── /repository      # Data Access Layer (interacts with the DB)
│       │   ├── account.repository.go
│       │   ├── category.repository.go
│       │   ├── transaction.repository.go
│       │   └── user.repository.go
│       ├── /routes          # API route definitions
│       │   └── routes.go
│       ├── /services        # Business logic layer
│       │   ├── account.service.go
│       │   ├── budget.service.go
│       │   ├── category.service.go
│       │   ├── dashboard.service.go
│       │   ├── recurring_transaction.service.go
│       │   ├── report.service.go
│       │   ├── transaction.service.go
│       │   └── user.service.go
│       ├── /utils           # Helpers (password, token, validation)
│       │   ├── password.go
│       │   ├── ping.go
│       │   ├── response.go
│       │   └── token.go
│       └── /pkg
│           └── /scheduler       # Background jobs (recurring transactions)
│               └── scheduler.go
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

#### `recurring_frequency`

|Value|
|---|
|monthly|
|yearly|

---
### 👤 `users` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|name|VARCHAR(100)|Not Null|
|email|VARCHAR(255)|Unique, Not Null|
|password|VARCHAR(255)|Nullable|
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
|category_id|UUID|Foreign Key → `categories(id)`, On Delete RESTRICT, Nullable|
|budget_id|UUID|Foreign Key → `budgets(id)`, On Delete SET NULL, Nullable|
|description|VARCHAR(255)|Not Null|
|amount|NUMERIC(19,4)|Not Null|
|type|transaction_type|Not Null|
|transaction_date|DATE|Not Null|
|note|TEXT|Optional|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|
|updated_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

---

### 🔁 `recurring_transactions` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|user_id|UUID|Foreign Key → `users(id)`, On Delete CASCADE|
|account_id|UUID|Foreign Key → `accounts(id)`, On Delete CASCADE|
|category_id|UUID|Foreign Key → `categories(id)`, On Delete RESTRICT|
|description|VARCHAR(255)|Not Null|
|amount|NUMERIC(19,4)|Not Null|
|type|transaction_type|Not Null|
|recurring_frequency|recurring_frequency|Not Null|
|recurring_date|INTEGER|Not Null|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|
|updated_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

---

### 💰 `budgets` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|user_id|UUID|Foreign Key → `users(id)`, On Delete CASCADE|
|name|VARCHAR(100)|Not Null|
|amount|NUMERIC(19,4)|Not Null|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|
|updated_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

---

### ⚡ Indexes

| Index Name | Columns |
| --- | --- |
| idx_transactions_user_id_date | `(user_id, transaction_date DESC)` |
| idx_accounts_user_id | `(user_id)` |

---

### 🔐 `jwt_tokens` Table

|Column|Type|Constraints / Default|
|---|---|---|
|id|UUID|Primary Key, Default: `gen_random_uuid()`|
|user_id|UUID|Foreign Key → `users(id)`, On Delete CASCADE, Unique|
|token|TEXT|Unique, Not Null|
|expires_at|TIMESTAMPTZ|Not Null|
|created_at|TIMESTAMPTZ|Not Null, Default: `NOW()`|

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
    

## 6. API Endpoints

### Base Route
- `GET /api/v1/`

### Auth Routes
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/profile`
- `POST /api/v1/auth/change-password`
- `GET /api/v1/auth/google/login`
- `GET /api/v1/auth/google/callback`

### Accounts Routes
- `POST /api/v1/accounts/create`
- `GET /api/v1/accounts/`
- `PATCH /api/v1/accounts/update/:id`
- `DELETE /api/v1/accounts/delete/:id`
- `GET /api/v1/accounts/total-balance`

### Transactions Routes
- `POST /api/v1/transactions/create`
- `GET /api/v1/transactions/`
- `PATCH /api/v1/transactions/update/:id`
- `DELETE /api/v1/transactions/delete/:id`
- `GET /api/v1/transactions/aggregate`

### Dashboard Routes
- `GET /api/v1/dashboard/`

### Reports Routes
- `GET /api/v1/reports/`
- `GET /api/v1/reports/export`

### Categories Routes
- `POST /api/v1/categories/create`
- `GET /api/v1/categories/`
- `PATCH /api/v1/categories/update/:id`
- `DELETE /api/v1/categories/delete/:id`

### Budgets Routes
- `POST /api/v1/budgets/create`
- `GET /api/v1/budgets/`
- `PATCH /api/v1/budgets/update/:id`
- `DELETE /api/v1/budgets/delete/:id`

### Recurring Transactions Routes
- `POST /api/v1/recurring-transactions/create`
- `GET /api/v1/recurring-transactions/`
- `PATCH /api/v1/recurring-transactions/update/:id`
- `DELETE /api/v1/recurring-transactions/delete/:id`

## 7. API Specification (Revised)

All API responses will adhere to the following structure:

```json
{
  "success": true, // or false
  "message": "Descriptive message",
  "data": { ... }, // or null
  "error": "Error details if success is false" // or null
}
```

---

### **`/api/v1/auth`**

- **Endpoint: `POST /api/v1/auth/register`**
    
    - **Description:** Registers a new user.
        
    - **Authorization:** Public
        
    - **Request Body:**
        
        ```json
        {
          "name": "John Doe",
          "email": "john.doe@example.com",
          "password": "aVeryStrongPassword123!"
        }
        ```
        
    - **Success Response (201 Created):**
        
        ```json
        {
          "success": true,
          "message": "User registered successfully",
          "data": {
            "user": {
              "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
              "name": "John Doe",
              "email": "john.doe@example.com",
              "provider": "email",
              "createdAt": "2025-10-09T10:00:00Z"
            },
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
          },
          "error": null
        }
        ```

- **Endpoint: `POST /api/v1/auth/login`**
    
    - **Description:** Authenticates a user and returns a JWT.
        
    - **Authorization:** Public
        
    - **Request Body:**
        
        ```json
        {
          "email": "john.doe@example.com",
          "password": "aVeryStrongPassword123!"
        }
        ```
        
    - **Success Response (200 OK):**
        
        ```json
        {
          "success": true,
          "message": "Login successful",
          "data": {
            "user": {
              "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
              "name": "John Doe",
              "email": "john.doe@example.com",
              "provider": "email",
              "createdAt": "2025-10-09T10:00:00Z"
            },
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/auth/profile`**
    
    - **Description:** Retrieves the profile for the authenticated user.
        
    - **Authorization:** Authenticated User
        
    - **Success Response (200 OK):**
        
        ```json
        {
          "success": true,
          "message": "Profile retrieved successfully",
          "data": {
            "personal": {
              "fullName": "John Doe",
              "email": "john.doe@example.com"
            }
          },
          "error": null
        }
        ```

- **Endpoint: `POST /api/v1/auth/change-password`**
    
    - **Description:** Allows an authenticated user to change their password.
        
    - **Authorization:** Authenticated User
        
    - **Request Body:**
        
        ```json
        {
          "currentPassword": "aVeryStrongPassword123!",
          "newPassword": "aNewerEvenStrongerPassword789!"
        }
        ```
        
    - **Success Response (200 OK):**
        
        ```json
        {
          "success": true,
          "message": "Password changed successfully",
          "data": null,
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/auth/google/login`**

    - **Description:** Initiates Google OAuth 2.0 login flow.

- **Endpoint: `GET /api/v1/auth/google/callback`**

    - **Description:** Handles the callback from Google OAuth 2.0.

---

### **`/api/v1/accounts`**

- **Endpoint: `POST /api/v1/accounts/create`**

    - **Description:** Creates a new financial account.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "name": "My Savings Account",
          "type": "savings",
          "balance": 1000.00
        }
        ```
    - **Success Response (201 Created):**
        ```json
        {
          "success": true,
          "message": "Account created successfully",
          "data": {
            "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "name": "My Savings Account",
            "type": "savings",
            "balance": 1000.00,
            "is_active": true,
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:00:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/accounts/`**

    - **Description:** Retrieves all financial accounts for the authenticated user.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Accounts retrieved successfully",
          "data": [
            {
              "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
              "name": "My Savings Account",
              "type": "savings",
              "balance": 1000.00,
              "is_active": true,
              "created_at": "2025-10-09T10:00:00Z",
              "updated_at": "2025-10-09T10:00:00Z"
            }
          ],
          "error": null
        }
        ```

- **Endpoint: `PATCH /api/v1/accounts/update/:id`**

    - **Description:** Updates a financial account.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "name": "My Updated Savings Account",
          "type": "savings",
          "is_active": false
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Account updated successfully",
          "data": {
            "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "name": "My Updated Savings Account",
            "type": "savings",
            "balance": 1000.00,
            "is_active": false,
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:05:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `DELETE /api/v1/accounts/delete/:id`**

    - **Description:** Deletes a financial account.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Account deleted successfully",
          "data": null,
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/accounts/total-balance`**

    - **Description:** Retrieves the total balance of all active accounts.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Total balance retrieved successfully",
          "data": {
            "total_balance": 5000.00
          },
          "error": null
        }
        ```

---

### **`/api/v1/transactions`**

- **Endpoint: `POST /api/v1/transactions/create`**

    - **Description:** Creates a new transaction.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
          "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
          "budget_id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6", // Optional
          "description": "Groceries",
          "amount": 75.50,
          "type": "expense",
          "transaction_date": "2025-10-09"
        }
        ```
    - **Success Response (201 Created):**
        ```json
        {
          "success": true,
          "message": "Transaction created successfully",
          "data": {
            "id": "c1d2e3f4-g5h6-i7j8-k9l0-m1n2o3p4q5r6",
            "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
            "budget_id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "description": "Groceries",
            "amount": 75.50,
            "type": "expense",
            "transaction_date": "2025-10-09",
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:00:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/transactions/`**

    - **Description:** Retrieves all transactions for the authenticated user.
    - **Authorization:** Authenticated User
    - **Produce:**  json
    - **Param** page query int false "Page number"
    - **Param** limit query int false "Number of items per page"
    - **Param** description query string false "Filter by description"
    - **Param** category query string false "Filter by category ID"
    - **Param** account query string false "Filter by account ID"
    - **Param** startDate query string false "Filter by start date (YYYY-MM-DD)"
    - **Param** endDate query string false "Filter by end date (YYYY-MM-DD)"
    - **Param** budget query string false "Filter by budget ID"

- **Endpoint: `PATCH /api/v1/transactions/update/:id`**

    - **Description:** Updates a transaction.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
          "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
          "budget_id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6", // Optional
          "description": "Weekly Groceries",
          "amount": 80.00,
          "type": "expense",
          "transaction_date": "2025-10-09"
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Transaction updated successfully",
          "data": {
            "id": "c1d2e3f4-g5h6-i7j8-k9l0-m1n2o3p4q5r6",
            "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
            "budget_id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "description": "Weekly Groceries",
            "amount": 80.00,
            "type": "expense",
            "transaction_date": "2025-10-09",
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:10:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `DELETE /api/v1/transactions/delete/:id`**

    - **Description:** Deletes a transaction.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Transaction deleted successfully",
          "data": null,
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/transactions/aggregate`**

    - **Description:** Retrieves aggregate data for transactions.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Aggregate data retrieved successfully",
          "data": {
            "total_income": 5000.00,
            "total_expense": 1250.75,
            "net_income": 3749.25
          },
          "error": null
        }
        ```

---

### **`/api/v1/dashboard`**

- **Endpoint: `GET /api/v1/dashboard`**

    - **Description:** Retrieves a summary of the user's financial data for the dashboard.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Dashboard summary retrieved successfully",
          "data": {
            "total_balance": 5000.00,
            "monthly_income": 2000.00,
            "monthly_expense": 850.50,
            "monthly_savings": 1149.50,
            "recent_transactions": [
              {
                "id": "d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6",
                "description": "Salary",
                "amount": 2000.00,
                "type": "income",
                "transaction_date": "2025-10-01"
              }
            ]
          },
          "error": null
        }
        ```

---

### **`/api/v1/reports`**

- **Endpoint: `GET /api/v1/reports`**

    - **Description:** Generates a financial report.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Report generated successfully",
          "data": {
            "start_date": "2025-01-01",
            "end_date": "2025-12-31",
            "total_income": 24000.00,
            "total_expense": 15000.00,
            "net_income": 9000.00,
            "spending_by_category": [
              {
                "category": "Groceries",
                "total": 3000.00
              },
              {
                "category": "Rent",
                "total": 12000.00
              }
            ]
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/reports/export`**

    - **Description:** Exports transaction data to a CSV file.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        - **Content-Type:** `text/csv`
        - **Body:** (CSV file content)

---

### **`/api/v1/categories`**

- **Endpoint: `POST /api/v1/categories/create`**

    - **Description:** Creates a new transaction category.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "name": "Utilities",
          "type": "expense"
        }
        ```
    - **Success Response (201 Created):**
        ```json
        {
          "success": true,
          "message": "Category created successfully",
          "data": {
            "id": "e1f2g3h4-i5j6-k7l8-m9n0-o1p2q3r4s5t6",
            "name": "Utilities",
            "type": "expense"
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/categories`**

    - **Description:** Retrieves all transaction categories.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Categories retrieved successfully",
          "data": [
            {
              "id": "e1f2g3h4-i5j6-k7l8-m9n0-o1p2q3r4s5t6",
              "name": "Utilities",
              "type": "expense"
            }
          ],
          "error": null
        }
        ```

- **Endpoint: `PATCH /api/v1/categories/update/:id`**

    - **Description:** Updates a transaction category.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "name": "Home Utilities"
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Category updated successfully",
          "data": {
            "id": "e1f2g3h4-i5j6-k7l8-m9n0-o1p2q3r4s5t6",
            "name": "Home Utilities",
            "type": "expense"
          },
          "error": null
        }
        ```

- **Endpoint: `DELETE /api/v1/categories/delete/:id`**

    - **Description:** Deletes a transaction category.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Category deleted successfully",
          "data": null,
          "error": null
        }
        ```

---

### **`/api/v1/budgets`**

- **Endpoint: `POST /api/v1/budgets/create`**

    - **Description:** Creates a new budget.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "name": "Monthly Groceries",
          "amount": 500.00
        }
        ```
    - **Success Response (201 Created):**
        ```json
        {
          "success": true,
          "message": "Budget created successfully",
          "data": {
            "id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "name": "Monthly Groceries",
            "amount": 500.00,
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:00:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/budgets/`**

    - **Description:** Retrieves all budgets for the authenticated user.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Budgets retrieved successfully",
          "data": [
            {
              "id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
              "name": "Monthly Groceries",
              "amount": 500.00,
              "created_at": "2025-10-09T10:00:00Z",
              "updated_at": "2025-10-09T10:00:00Z"
            }
          ],
          "error": null
        }
        ```

- **Endpoint: `PATCH /api/v1/budgets/update/:id`**

    - **Description:** Updates a budget.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "name": "Updated Monthly Groceries",
          "amount": 550.00
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Budget updated successfully",
          "data": {
            "id": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "name": "Updated Monthly Groceries",
            "amount": 550.00,
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:15:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `DELETE /api/v1/budgets/delete/:id`**

    - **Description:** Deletes a budget.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Budget deleted successfully",
          "data": null,
          "error": null
        }
        ```

---

### **`/api/v1/recurring-transactions`**

- **Endpoint: `POST /api/v1/recurring-transactions/create`**

    - **Description:** Creates a new recurring transaction.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
          "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
          "description": "Netflix Subscription",
          "amount": 15.99,
          "type": "expense",
          "recurring_frequency": "monthly",
          "recurring_date": 15
        }
        ```
    - **Success Response (201 Created):**
        ```json
        {
          "success": true,
          "message": "Recurring transaction created successfully",
          "data": {
            "id": "g1h2i3j4-k5l6-m7n8-o9p0-q1r2s3t4u5v6",
            "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
            "description": "Netflix Subscription",
            "amount": 15.99,
            "type": "expense",
            "recurring_frequency": "monthly",
            "recurring_date": 15,
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:00:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `GET /api/v1/recurring-transactions`**

    - **Description:** Retrieves all recurring transactions for the authenticated user.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Recurring transactions retrieved successfully",
          "data": [
            {
              "id": "g1h2i3j4-k5l6-m7n8-o9p0-q1r2s3t4u5v6",
              "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
              "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
              "description": "Netflix Subscription",
              "amount": 15.99,
              "type": "expense",
              "recurring_frequency": "monthly",
              "recurring_date": 15,
              "created_at": "2025-10-09T10:00:00Z",
              "updated_at": "2025-10-09T10:00:00Z"
            }
          ],
          "error": null
        }
        ```

- **Endpoint: `PATCH /api/v1/recurring-transactions/update/:id`**

    - **Description:** Updates a recurring transaction.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "amount": 16.99
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Recurring transaction updated successfully",
          "data": {
            "id": "g1h2i3j4-k5l6-m7n8-o9p0-q1r2s3t4u5v6",
            "account_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "category_id": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
            "description": "Netflix Subscription",
            "amount": 16.99,
            "type": "expense",
            "recurring_frequency": "monthly",
            "recurring_date": 15,
            "created_at": "2025-10-09T10:00:00Z",
            "updated_at": "2025-10-09T10:20:00Z"
          },
          "error": null
        }
        ```

- **Endpoint: `DELETE /api/v1/recurring-transactions/delete/:id`**

    - **Description:** Deletes a recurring transaction.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Recurring transaction deleted successfully",
          "data": null,
          "error": null
        }
        ```


## 7. Authentication & Authorization

### 7.1. Authentication Strategy

The system employs a JWT-based authentication strategy for its stateless API, with enhanced token management.

**Flow:**

1. A user submits credentials (email/password) or an OAuth token (from Google).
    
2. The server validates the credentials/token.
    
3. If valid, the server checks the `jwt_tokens` table for an existing valid token for the user.
    
4. If an existing valid token is found, it is returned. If the existing token is expired, it is removed.
    
5. If no valid token exists (or the previous one was removed), a new signed JWT is generated containing a payload with the `user_id` and an expiration timestamp (`exp`). This new token is stored in the `jwt_tokens` table (replacing any old token for that user).
    
6. The server sends this JWT back to the client.
    
7. The client application stores the JWT securely (e.g., in an HttpOnly cookie or secure storage).
    
8. For all subsequent requests to protected endpoints, the client must include the JWT in the `Authorization` header with the `Bearer` scheme (`Authorization: Bearer <token>`).
    
9. A middleware on the server intercepts each request, validates the JWT's signature and expiration, and crucially, checks if the token exists and is valid in the `jwt_tokens` table. If valid, it extracts the `user_id` to process the request within the user's scope.
    
10. Tokens will have a short lifespan (e.g., 1 hour), and a refresh token mechanism will be implemented for a seamless user experience.
    

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