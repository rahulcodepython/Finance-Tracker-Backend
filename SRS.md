# System Requirements Specification (SRS): Finance Tracker

Version: 1.0

Date: October 5, 2025

Modified: October 25, 2025

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
	
- Maintaining detailed logs for future reference.

**Out-of-Scope:**

- The development of the client-side frontend application (e.g., Next.js web app, mobile apps).
    
- Direct integration with third-party banking APIs for automatic transaction syncing (Plaid, etc.).
    
- System administration user interface.
    
- Deployment infrastructure setup and management (CI/CD pipelines, cloud hosting).
    
- CSV data import functionality (marked as a future feature).
    

## 2. Overall Description

### 2.1. User Personas & Roles

- **Standard User:** The primary user of the application. This user can register, log in, manage their own profile, accounts, transactions, budgets, logs, and view their financial data through dashboards and reports. All data is scoped to their own profile.
    

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

```markdown
/finance-tracker-api
├── .dockerignore
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── README.md
├── SRS.md
├── .git/
├── api/
│   └── v1/
│       ├── account.handler.go
│       ├── auth.handler.go
│       ├── budget.handler.go
│       ├── category.handler.go
│       ├── dashboard.handler.go
│       ├── log.handler.go
│       ├── recurring.transaction.handler.go
│       ├── report.handler.go
│       └── transaction.handler.go
├── backend/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── database.go
│   │   └── migrations.go
│   ├── interfaces/
│   │   └── sql.interfaces.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logger.go
│   ├── models/
│   │   ├── account.go
│   │   ├── budget.go
│   │   ├── category.go
│   │   ├── jwt.token.go
│   │   ├── log.go
│   │   ├── recurring.transaction.go
│   │   ├── transaction.go
│   │   └── user.go
│   ├── pkg/
│   │   └── scheduler/
│   │       └── scheduler.go
│   ├── repository/
│   │   ├── account.repository.go
│   │   ├── budget.repository.go
│   │   ├── category.repository.go
│   │   ├── jwt.token.repository.go
│   │   ├── log.repository.go
│   │   ├── recurring.transaction.repository.go
│   │   ├── transaction.repository.go
│   │   └── user.repository.go
│   ├── routes/
│   │   └── routes.go
│   ├── services/
│   │   ├── account.service.go
│   │   ├── budget.service.go
│   │   ├── category.service.go
│   │   ├── dashboard.service.go
│   │   ├── log.service.go
│   │   ├── recurring.transaction.service.go
│   │   ├── report.service.go
│   │   ├── transaction.service.go
│   │   └── user.service.go
│   └── utils/
│       ├── db.transaction.go
│       ├── password.go
│       ├── ping.go
│       ├── response.go
│       ├── time.go
│       └── token.go
└── migrations/
    ├── drop.sql
    └── schema.sql
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
    
## Enumerated Types

| Type Name | Values | Description |
|-----------|--------|-------------|
| `account_type` | `checking`, `savings`, `credit_card`, `cash`, `investment`, `loan`, `upi` | Types of financial accounts |
| `transaction_type` | `income`, `expense` | Types of financial transactions |
| `auth_provider` | `email`, `google` | User authentication methods |
| `recurring_frequency` | `monthly`, `yearly` | Frequencies for recurring transactions |

## Tables

### Users Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique user identifier |
| `name` | VARCHAR(100) | NOT NULL | User's full name |
| `email` | VARCHAR(255) | UNIQUE, NOT NULL | User's email address |
| `password` | VARCHAR(255) | NULL | Hashed password (nullable for OAuth) |
| `provider` | auth_provider | NOT NULL, DEFAULT 'email' | Authentication provider |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Account creation timestamp |

### Accounts Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique account identifier |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Associated user |
| `name` | VARCHAR(100) | NOT NULL | Account name |
| `type` | account_type | NOT NULL | Type of account |
| `balance` | NUMERIC(19,4) | NOT NULL, DEFAULT 0.00 | Current account balance |
| `is_active` | BOOLEAN | NOT NULL, DEFAULT TRUE | Account status |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Account creation timestamp |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Last update timestamp |

### Categories Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique category identifier |
| `name` | VARCHAR(100) | NOT NULL | Category name |
| `type` | transaction_type | NOT NULL | Transaction type for this category |
| - | - | UNIQUE (name, type) | Ensures unique category per transaction type |

### Transactions Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique transaction identifier |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Associated user |
| `account_id` | UUID | NOT NULL, REFERENCES accounts(id) ON DELETE CASCADE | Source/destination account |
| `category_id` | UUID | REFERENCES categories(id) ON DELETE RESTRICT | Transaction category |
| `budget_id` | UUID | REFERENCES budgets(id) ON DELETE SET NULL | Associated budget |
| `description` | VARCHAR(255) | NOT NULL | Transaction description |
| `amount` | NUMERIC(19,4) | NOT NULL | Transaction amount |
| `type` | transaction_type | NOT NULL | Income or expense |
| `transaction_date` | DATE | NOT NULL | Date of transaction |
| `note` | TEXT | - | Additional notes |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Last update timestamp |

### Recurring Transactions Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique recurring transaction identifier |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Associated user |
| `account_id` | UUID | NOT NULL, REFERENCES accounts(id) ON DELETE CASCADE | Source/destination account |
| `category_id` | UUID | REFERENCES categories(id) ON DELETE RESTRICT | Transaction category |
| `budget_id` | UUID | REFERENCES budgets(id) ON DELETE SET NULL | Associated budget |
| `description` | VARCHAR(255) | NOT NULL | Transaction description |
| `amount` | NUMERIC(19,4) | NOT NULL | Transaction amount |
| `type` | transaction_type | NOT NULL | Income or expense |
| `note` | TEXT | - | Additional notes |
| `recurring_frequency` | recurring_frequency | NOT NULL | Monthly or yearly recurrence |
| `recurring_date` | INTEGER | NOT NULL | Day of month for recurring transactions |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Last update timestamp |

### Budgets Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique budget identifier |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Associated user |
| `name` | VARCHAR(100) | NOT NULL | Budget name |
| `amount` | NUMERIC(19,4) | NOT NULL | Budget amount |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Last update timestamp |

### JWT Tokens Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique token identifier |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Associated user |
| `token` | TEXT | UNIQUE, NOT NULL | JWT token string |
| `expires_at` | TIMESTAMPTZ | NOT NULL | Token expiration timestamp |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Creation timestamp |
| - | - | UNIQUE (user_id) | One active token per user |

### Logs Table
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Log entry identifier |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Associated user |
| `message` | TEXT | NOT NULL | Log message content |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Log creation timestamp |

## Indexes

| Index Name | Table | Columns | Description |
|------------|-------|---------|-------------|
| `idx_transactions_user_id_date` | transactions | (user_id, transaction_date DESC) | Optimizes user transaction queries by date |
| `idx_accounts_user_id` | accounts | (user_id) | Optimizes user account lookups |

## Key Relationships

- **Users → Accounts**: One-to-Many (CASCADE delete)
- **Users → Transactions**: One-to-Many (CASCADE delete)  
- **Users → Budgets**: One-to-Many (CASCADE delete)
- **Accounts → Transactions**: One-to-Many (CASCADE delete)
- **Categories → Transactions**: One-to-Many (RESTRICT delete)
- **Budgets → Transactions**: One-to-Many (SET NULL delete)
- **Users → JWT Tokens**: One-to-One (CASCADE delete)


## Functional Requirements

### 1. Authentication & Authorization Module
#### User Management
- **User Registration** with email/password or OAuth (Google)
- **User Login/Logout** with JWT token management
- **Password Management** - secure hashing, reset functionality
- **Session Management** - token expiration and refresh
- **Profile Management** - update user information

#### Security Features
- **Authentication Middleware** - protect routes and validate tokens
- **Provider-based Auth** - support multiple authentication methods

### 2. Account Management Module
#### Core Account Operations
- **Create/Update/Delete** financial accounts
- **Account Type Support** - checking, savings, credit cards, cash, investments, loans, UPI
- **Balance Management** - track current balances with precision (4 decimal places)
- **Account Status** - activate/deactivate accounts
- **Account Categorization** - organize by type and status

### 3. Transaction Management Module
#### Basic Transaction Operations
- **Record Transactions** - income and expense tracking
- **Transaction Categorization** - assign to predefined categories
- **Transaction Editing** - modify existing transactions
- **Transaction Deletion** - with proper constraints

#### Advanced Features
- **Transaction Search** - filter by date, category, amount, description

### 4. Category Management Module
- **Category CRUD** - create, read, update, delete categories
- **Type-based Categories** - separate categories for income and expense
- **Category Validation** - ensure unique category names per type
- **Default Categories** - pre-defined common categories

### 5. Budget Management Module
#### Budget Planning
- **Budget Creation** - set spending limits by period
- **Budget Tracking** - monitor actual vs planned spending
- **Budget Categories** - associate transactions with budgets
- **Budget Alerts** - notifications when approaching limits

#### Advanced Features
- **Rollover Budgets** - handle unused amounts
- **Multiple Budget Periods** - weekly, monthly, yearly
- **Budget Templates** - reusable budget structures

### 6. Recurring Transactions Module
#### Automated Transactions
- **Recurring Setup** - configure automatic transaction generation
- **Frequency Support** - monthly and yearly recurrences
- **Date Management** - specific day of month for processing
- **Recurring Template Management** - create and modify templates

#### Advanced Features
- **Recurring Pattern Validation** - ensure valid recurrence rules
- **Auto-generation** - system-generated transactions
- **Recurring Budget Alignment** - integrate with budget planning

### 7. Reporting & Analytics Module
#### Financial Reports
- **Income Statement** - revenue vs expenses over time
- **Spending Analysis** - category-wise expenditure
- **Account Balances** - net worth tracking
- **Budget vs Actual** - performance reporting

#### Advanced Analytics
- **Trend Analysis** - spending patterns over time
- **Forecasting** - future financial projections
- **Custom Reports** - user-defined report generation
- **Data Visualization** - charts and graphs

### 8. Dashboard Module
#### Overview Features
- **Financial Snapshot** - current balances and recent activity
- **Quick Actions** - fast access to common operations
- **Alert Summary** - important notifications and warnings
- **Performance Metrics** - key financial indicators

### 9. Logging & Audit Module
- **Activity Logging** - track user actions and system events
- **Audit Trail** - compliance and debugging support
- **Error Tracking** - system error monitoring
- **Performance Logging** - response time and resource usage

## Advanced System Features

### 1. Data Management
- **Database Transactions** - ensure data consistency
- **Data Validation** - input sanitization and business rule enforcement

### 2. System Integration
- **RESTful API** - standardized API responses and error handling
- **Export Capabilities** - CSV, formats

## Non-Functional Requirements

### 1. Performance
- **Response Time**: API responses under 200ms for 95% of requests
- **Transaction Processing**: Handle 1000+ concurrent transactions
- **Database Queries**: Optimized queries with proper indexing
- **Throughput**: Support 10,000+ daily active users

### 2. Scalability
- **Horizontal Scaling**: Stateless architecture supporting multiple instances
- **Database Scaling**: Read replicas for reporting and analytics
- **Load Balancing**: Efficient request distribution
- **Resource Management**: Auto-scaling based on load patterns

### 3. Security
- **Data Encryption**: AES-256 for sensitive data at rest
- **TLS/SSL**: HTTPS for all communications
- **Authentication**: JWT with short expiration and secure refresh
- **Authorization**: Role-based access control (RBAC)
- **Input Validation**: SQL injection and XSS protection
- **API Security**: Rate limiting and DDoS protection

### 4. Reliability
- **Uptime**: 99.9% availability SLA
- **Data Consistency**: ACID compliance for financial transactions
- **Error Handling**: Graceful degradation and informative error messages
- **Backup Strategy**: Automated daily backups with point-in-time recovery
- **Disaster Recovery**: Multi-region deployment capability

### 5. Usability
- **Response Consistency**: Standardized API response format
- **Error Messages**: User-friendly, actionable error information
- **API Documentation**: Comprehensive OpenAPI/Swagger documentation
- **Mobile Responsive**: Responsive design for various devices
- **Accessibility**: WCAG 2.1 compliance for web interfaces

### 6. Maintainability
- **Code Quality**: Comprehensive test coverage (unit, integration, e2e)
- **Documentation**: API docs, architecture decisions, deployment guides
- **Monitoring**: Application performance monitoring (APM)
- **Logging**: Structured logging with correlation IDs
- **CI/CD**: Automated testing and deployment pipelines

### 9. Data Integrity
- **Referential Integrity**: Database constraints and cascading rules
- **Audit Trail**: Immutable transaction history
- **Data Validation**: Business rule enforcement at multiple layers
- **Consistency Checks**: Regular data integrity verification

Based on the API structure and route grouping information provided, here's the updated API endpoints documentation:

## 6. API Endpoints

### Base Route
- `GET /api/v1/` - **Public** - API welcome and health check

### Authentication Module
- `POST /api/v1/auth/register` - **Public** - User registration
- `POST /api/v1/auth/login` - **Public** - User login
- `GET /api/v1/auth/profile` - **Authenticated** - Get user profile (Own data only)
- `POST /api/v1/auth/change-password` - **Authenticated** - Change password (Own data only)
- `GET /api/v1/auth/google/login` - **Public** - Initiate Google OAuth flow
- `GET /api/v1/auth/google/callback` - **Public** - Google OAuth callback

### Account Management Module
- `POST /api/v1/accounts/create` - **Authenticated** - Create financial account (User-owned accounts)
- `GET /api/v1/accounts/` - **Authenticated** - Get all user accounts (User-owned accounts)
- `PATCH /api/v1/accounts/update/:id` - **Authenticated** - Update account (User-owned accounts)
- `DELETE /api/v1/accounts/delete/:id` - **Authenticated** - Delete account (User-owned accounts)
- `GET /api/v1/accounts/total-balance` - **Authenticated** - Get total balance (User-owned accounts)

### Transaction Management Module
- `POST /api/v1/transactions/create` - **Authenticated** - Create transaction (User-owned transactions)
- `GET /api/v1/transactions/` - **Authenticated** - Get all transactions (User-owned transactions)
- `PATCH /api/v1/transactions/update/:id` - **Authenticated** - Update transaction (User-owned transactions)
- `DELETE /api/v1/transactions/delete/:id` - **Authenticated** - Delete transaction (User-owned transactions)
- `GET /api/v1/transactions/aggregate` - **Authenticated** - Get aggregated transaction data (User-owned transactions)

### Dashboard Module
- `GET /api/v1/dashboard/` - **Authenticated** - Get financial overview and analytics (User data aggregation)

### Reporting Module
- `GET /api/v1/reports/` - **Authenticated** - Generate financial reports (User data only)
- `GET /api/v1/reports/export` - **Authenticated** - Export transactions (User data only)

### Category Management Module
- `POST /api/v1/categories/create` - **Authenticated** - Create category (System + user categories)
- `GET /api/v1/categories/` - **Authenticated** - Get all categories (System + user categories)
- `PATCH /api/v1/categories/update/:id` - **Authenticated** - Update category (System + user categories)
- `DELETE /api/v1/categories/delete/:id` - **Authenticated** - Delete category (System + user categories)

### Budget Management Module
- `POST /api/v1/budgets/create` - **Authenticated** - Create budget (User-owned budgets)
- `GET /api/v1/budgets/` - **Authenticated** - Get all budgets (User-owned budgets)
- `PATCH /api/v1/budgets/update/:id` - **Authenticated** - Update budget (User-owned budgets)
- `DELETE /api/v1/budgets/delete/:id` - **Authenticated** - Delete budget (User-owned budgets)

### Recurring Transactions Module
- `POST /api/v1/recurring-transactions/create` - **Authenticated** - Create recurring transaction (User-owned recurring transactions)
- `GET /api/v1/recurring-transactions/` - **Authenticated** - Get all recurring transactions (User-owned recurring transactions)
- `PATCH /api/v1/recurring-transactions/update/:id` - **Authenticated** - Update recurring transaction (User-owned recurring transactions)
- `DELETE /api/v1/recurring-transactions/delete/:id` - **Authenticated** - Delete recurring transaction (User-owned recurring transactions)

### System Logs Module
- `GET /api/v1/logs/` - **Authenticated** - Get user activity logs (User activity logs)

**Note**: All authenticated endpoints require the `DeserializeUser` middleware and enforce data scope restrictions to ensure users can only access their own data.

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
          }
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
          }
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
              "name": "John Doe",
              "email": "john.doe@example.com"
            }
          }
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
          "message": "Password changed successfully"
        }
        ```

- **Endpoint: `GET /api/v1/auth/google/login`**

    - **Description:** Initiates Google OAuth 2.0 login flow. Redirects the user to the Google login page.

- **Endpoint: `GET /api/v1/auth/google/callback`**

    - **Description:** Handles the callback from Google OAuth 2.0. Exchanges the authorization code for an access token, fetches user info, and then logs in or creates a new user.
    
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
              "provider": "google",
              "createdAt": "2025-10-09T10:00:00Z"
            },
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
          }
        }
        ```

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
            "userId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "name": "My Savings Account",
            "type": "savings",
            "balance": 1000.00,
            "isActive": true,
            "createdAt": "2025-10-09T10:00:00Z",
            "updatedAt": "2025-10-09T10:00:00Z"
          }
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
              "userId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
              "name": "My Savings Account",
              "type": "savings",
              "balance": 1000.00,
              "isActive": true,
              "createdAt": "2025-10-09T10:00:00Z",
              "updatedAt": "2025-10-09T10:00:00Z"
            }
          ]
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
          "isActive": false
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Account updated successfully",
          "data": {
            "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "userId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "name": "My Updated Savings Account",
            "type": "savings",
            "balance": 1000.00,
            "isActive": false,
            "createdAt": "2025-10-09T10:00:00Z",
            "updatedAt": "2025-10-09T10:05:00Z"
          }
        }
        ```

- **Endpoint: `DELETE /api/v1/accounts/delete/:id`**

    - **Description:** Deletes a financial account.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Account deleted successfully"
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
          }
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
          "accountId": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
          "categoryId": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
          "budgetId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6", // Optional
          "description": "Groceries",
          "amount": 75.50,
          "date": "2025-10-09",
          "note": "Weekly grocery shopping"
        }
        ```
    - **Success Response (201 Created):**
        ```json
        {
          "success": true,
          "message": "Transaction created successfully",
          "data": {
            "id": "c1d2e3f4-g5h6-i7j8-k9l0-m1n2o3p4q5r6",
            "userId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "accountId": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "categoryId": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
            "budgetId": {
                "UUID": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
                "Valid": true
            },
            "description": "Groceries",
            "amount": 75.50,
            "type": "expense",
            "transactionDate": "2025-10-09T00:00:00Z",
            "note": {
                "String": "Weekly grocery shopping",
                "Valid": true
            },
            "createdAt": "2025-10-09T10:00:00Z",
            "updatedAt": "2025-10-09T10:00:00Z"
          }
        }
        ```

- **Endpoint: `GET /api/v1/transactions/`**

    - **Description:** Retrieves all transactions for the authenticated user.
    - **Authorization:** Authenticated User
    - **Query Parameters:**
        - `page` (int, optional): Page number (default: 1)
        - `limit` (int, optional): Number of items per page (default: 10)
        - `description` (string, optional): Filter by description
        - `category` (string, optional): Filter by category ID
        - `account` (string, optional): Filter by account ID
        - `budget` (string, optional): Filter by budget ID
        - `startDate` (string, optional): Filter by start date (YYYY-MM-DD)
        - `endDate` (string, optional): Filter by end date (YYYY-MM-DD)
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Transactions retrieved successfully",
          "data": [
            {
              "id": "c1d2e3f4-g5h6-i7j8-k9l0-m1n2o3p4q5r6",
              "userId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
              "accountId": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
              "categoryId": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
              "budgetId": {
                "UUID": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
                "Valid": true
              },
              "description": "Groceries",
              "amount": 75.50,
              "type": "expense",
              "transactionDate": "2025-10-09T00:00:00Z",
              "note": {
                "String": "Weekly grocery shopping",
                "Valid": true
              },
              "createdAt": "2025-10-09T10:00:00Z",
              "updatedAt": "2025-10-09T10:00:00Z"
            }
          ]
        }
        ```

- **Endpoint: `PATCH /api/v1/transactions/update/:id`**

    - **Description:** Updates a transaction.
    - **Authorization:** Authenticated User
    - **Request Body:**
        ```json
        {
          "accountId": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
          "categoryId": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
          "budgetId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6", // Optional
          "description": "Weekly Groceries",
          "amount": 80.00,
          "date": "2025-10-09",
          "note": "Updated note"
        }
        ```
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Transaction updated successfully",
          "data": {
            "id": "c1d2e3f4-g5h6-i7j8-k9l0-m1n2o3p4q5r6",
            "userId": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
            "accountId": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
            "categoryId": "b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6",
            "budgetId": {
                "UUID": "f1g2h3i4-j5k6-l7m8-n9o0-p1q2r3s4t5u6",
                "Valid": true
            },
            "description": "Weekly Groceries",
            "amount": 80.00,
            "type": "expense",
            "transactionDate": "2025-10-09T00:00:00Z",
            "note": {
                "String": "Updated note",
                "Valid": true
            },
            "createdAt": "2025-10-09T10:00:00Z",
            "updatedAt": "2025-10-09T10:10:00Z"
          }
        }
        ```

- **Endpoint: `DELETE /api/v1/transactions/delete/:id`**

    - **Description:** Deletes a transaction.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Transaction deleted successfully"
        }
        ```

- **Endpoint: `GET /api/v1/transactions/aggregate`**

    - **Description:** Retrieves aggregate data for transactions.
    - **Authorization:** Authenticated User
    - **Query Parameters:**
        - `startDate` (string, optional): Start date (YYYY-MM-DD)
        - `endDate` (string, optional): End date (YYYY-MM-DD)
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Aggregate data retrieved successfully",
          "data": {
            "total_income": 5000.00,
            "total_expense": 1250.75,
            "net_income": 3749.25
          }
        }
        ```

---

### **`/api/v1/dashboard`**

- **Endpoint: `GET /api/v1/dashboard/`**

    - **Description:** Retrieves a summary of the user's financial data for the dashboard.
    - **Authorization:** Authenticated User
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Dashboard data retrieved successfully",
          "data": {
            "summary": {
              "totalBalance": {
                "total_balance": 5000.00
               },
              "monthlyIncome": 2000.00,
              "monthlyExpenses": 850.50,
              "monthlySavings": 1149.50
            },
            "graphs": {
              "incomeVsExpense": [],
              "spendingByCategory": [
                {
                  "category": "Groceries",
                  "total": 3000.00
                }
              ],
              "earningByCategory": [
                {
                  "category": "Salary",
                  "total": 5000.00
                }
              ]
            },
            "recentTransactions": [
              {
                "id": "d1e2f3g4-h5i6-j7k8-l9m0-n1o2p3q4r5s6",
                "description": "Salary",
                "amount": 2000.00,
                "type": "income",
                "transaction_date": "2025-10-01"
              }
            ]
          }
        }
        ```

---

### **`/api/v1/reports`**

- **Endpoint: `GET /api/v1/reports/`**

    - **Description:** Generates a financial report.
    - **Authorization:** Authenticated User
    - **Query Parameters:**
        - `from` (string, optional): Start date (YYYY-MM-DD)
        - `to` (string, optional): End date (YYYY-MM-DD)
    - **Success Response (200 OK):**
        ```json
        {
          "success": true,
          "message": "Report generated successfully",
          "data": {
            "summary": {
              "total_income": 24000.00,
              "total_expense": 15000.00,
              "net_income": 9000.00
            },
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
          }
        }
        ```

- **Endpoint: `GET /api/v1/reports/export`**

    - **Description:** Exports transaction data to a CSV file.
    - **Authorization:** Authenticated User
    - **Query Parameters:**
        - `from` (string, optional): Start date (YYYY-MM-DD)
        - `to` (string, optional): End date (YYYY-MM-DD)
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

The system implements a secure, stateless JWT-based authentication system with enhanced token management and refresh capabilities.

**Authentication Flow:**

1. **Credential Submission**: User submits email/password or initiates OAuth flow with Google
2. **Credential Validation**: Server validates credentials against database or OAuth provider
3. **Token Check**: System queries `jwt_tokens` table for existing valid tokens for the user
4. **Token Management**:
   - If valid token exists: return existing token
   - If token expired: remove expired token and generate new one
   - If no valid token: generate new JWT token
5. **Token Generation**: New JWT signed with `JWT_SECRET` containing:
   ```json
   {
     "user_id": "uuid",
     "exp": 168h_from_issue,
   }
   ```
6. **Token Storage**: New token stored in `jwt_tokens` table with user association
7. **Client Storage**: JWT securely stored in HttpOnly cookies or secure local storage
8. **Request Authentication**: Client includes JWT in `Authorization: Bearer <token>` header
9. **Middleware Validation**: Server validates JWT signature, expiration, and database existence
10. **Refresh Mechanism**: Short-lived access tokens (168h)

### 7.2. Authorization Strategy

The system employs Role-Based Access Control (RBAC) with resource-level ownership validation.

**Roles & Permissions:**

| Role | Description | Access Scope |
|------|-------------|--------------|
| **Public** | Unauthenticated users | Authentication endpoints only |
| **Authenticated User** | Verified system users | Full access to owned resources |

**Endpoint Authorization Matrix:**

| Module | Endpoint Pattern | Required Role | Data Scope | Description |
|--------|------------------|---------------|------------|-------------|
| **Authentication** | `/api/v1/auth/*` | Public | N/A | Registration, login, OAuth flows |
| **User Management** | `/api/v1/users/me` | Authenticated | Own data only | Profile management operations |
| **Account Management** | `/api/v1/accounts/*` | Authenticated | User-owned accounts | Full CRUD on user's financial accounts |
| **Transaction Management** | `/api/v1/transactions/*` | Authenticated | User-owned transactions | Complete transaction lifecycle management |
| **Category Management** | `/api/v1/categories/*` | Authenticated | System + user categories | Category setup and management |
| **Budget Management** | `/api/v1/budgets/*` | Authenticated | User-owned budgets | Budget planning and tracking |
| **Recurring Transactions** | `/api/v1/recurring/*` | Authenticated | User-owned recurring transactions | Automated transaction management |
| **Dashboard** | `/api/v1/dashboard/*` | Authenticated | User data aggregation | Financial overview and analytics |
| **Reporting** | `/api/v1/reports/*` | Authenticated | User data only | Financial reporting and insights |
| **System Logs** | `/api/v1/logs/*` | Authenticated | User activity logs | Audit trail and activity monitoring |

### 7.3. Security Implementation Details

**JWT Configuration:**
- **Access Token Expiry**: 1 hour (enhanced security)
- **Token Storage**: Database-persisted for revocation capability
- **Signature Algorithm**: HS256 with 32-character minimum secret

**Database-Level Security:**
- **CASCADE DELETE**: User deletion removes all associated data
- **RESTRICT DELETE**: Protected category deletions
- **SET NULL**: Optional relationships maintain data integrity
- **UUID Primary Keys**: Obfuscated resource identifiers

**API Security Measures:**
- **Rate Limiting**: 100 requests per minute per user
- **CORS Protection**: Configurable client origins
- **Input Validation**: Comprehensive request sanitization
- **SQL Injection Protection**: Parameterized queries throughout

## 8. Environment Configuration

The system configuration is managed through environment variables with the following structure:

```env
# =====================================
# Application Core Configuration
# =====================================
# Network binding configuration
HOST=localhost
PORT=8080

# Runtime environment and behavior
APP_ENV=development  # development|production|staging

# Cross-Origin Resource Sharing
CLIENT_ORIGIN=http://localhost:3000

# =====================================
# Database Configuration (PostgreSQL)
# =====================================
# Connection parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_postgres_password
DB_NAME=finance_tracker

# Connection security and performance
DB_SSL_MODE=disable  # require|verify-full|disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Connection string example:
# psql -U ${DB_USER} -d ${DB_NAME} -h ${DB_HOST} -p ${DB_PORT}

# =====================================
# Security & JWT Configuration
# =====================================
# JWT Signing and Validation
JWT_SECRET=minimum_32_character_super_secure_random_string
JWT_EXPIRES_IN=168h  # 7 days for refresh tokens
JWT_ACCESS_EXPIRES_IN=1h  # 1 hour for access tokens

# Token refresh configuration
JWT_REFRESH_ENABLED=true
JWT_REFRESH_EXPIRES_IN=168h

# =====================================
# External OAuth Services
# =====================================
# Google OAuth 2.0 Configuration
GOOGLE_CLIENT_ID=your_google_oauth_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_oauth_client_secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

# =====================================
# Rate Limiting & Performance
# =====================================
# General API rate limiting
RATE_LIMITER_MAX=100
RATE_LIMITER_DURATION_MINUTES=1

# Authentication-specific rate limits
AUTH_RATE_LIMITER_MAX=10
AUTH_RATE_LIMITER_DURATION_MINUTES=1

# =====================================
# Logging & Monitoring
# =====================================
# Log level and output configuration
LOG_LEVEL=info  # debug|info|warn|error
LOG_FORMAT=json  # json|text

# Monitoring and observability
METRICS_ENABLED=true
HEALTH_CHECK_ENDPOINT=/health

# =====================================
# Advanced Features
# =====================================
# Scheduled task configuration
SCHEDULER_ENABLED=true
RECURRING_TRANSACTION_HOUR=2  # 2 AM daily processing

# Data export and backup
BACKUP_ENABLED=true
BACKUP_SCHEDULE=0 2 * * *  # 2 AM daily
```

### Configuration Notes:

**Security Enhancements:**
- Separate access and refresh token expiration for balanced security and usability
- Database connection pooling for optimal performance
- Environment-specific SSL modes for database connections

**Production Considerations:**
- Set `APP_ENV=production` for production deployments
- Enable `DB_SSL_MODE=require` or `verify-full` in production
- Use strong, randomly generated `JWT_SECRET` (32+ characters)
- Configure proper `CLIENT_ORIGIN` for your frontend application

**Development Setup:**
- Default to `development` mode with detailed logging
- Local database with SSL disabled
- Extended token expiration for testing convenience

This configuration provides a robust foundation for both development and production environments while maintaining security best practices and system performance.