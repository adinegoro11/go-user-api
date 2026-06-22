# GEMINI.md

## Project Overview

**Project Name:** Book Store API

A production-grade REST API built with Go that supports:

* Authentication & Authorization
* User Management
* Book Catalog Management
* Inventory Management
* Shopping Cart
* Order Processing
* Payment Integration
* Reporting & Dashboard
* Audit Logging

The system follows Clean Architecture principles and is designed to be scalable, testable, and cloud-ready.

---

# Tech Stack

## Backend

* Go 1.24+
* Gin Framework
* GORM
* PostgreSQL
* Redis
* JWT Authentication
* Refresh Token
* Zap Logger

## Infrastructure

* Docker
* Docker Compose
* GitHub Actions

## Documentation

* Swagger/OpenAPI

## Testing

* Testify
* Mockery

---

# Architecture

Follow Clean Architecture.

```text
cmd/
    api/

internal/
    domain/
    repository/
    service/
    handler/
    middleware/
    dto/
    validator/

pkg/
    logger/
    database/
    jwt/
    cache/
    response/

migrations/

docs/

tests/
```

Dependency direction:

```text
Handler
    ↓
Service
    ↓
Repository
    ↓
Database
```

Handlers MUST NOT access the database directly.

Repositories MUST NOT contain business logic.

Services MUST contain business rules.

---

# User Roles

Only two roles exist:

```go
user
admin
```

## User Permissions

* Register
* Login
* Refresh Token
* View Own Profile
* Update Own Profile
* Browse Books
* Add To Cart
* Checkout
* View Own Orders

## Admin Permissions

Everything available to user plus:

* Manage Books
* Manage Inventory
* Manage Users
* View All Orders
* Manage Order Status
* Access Reports
* Access Dashboard

---

# Authentication

Authentication uses:

* JWT Access Token
* Refresh Token

Access token lifetime:

```text
15 minutes
```

Refresh token lifetime:

```text
7 days
```

Passwords MUST be hashed using bcrypt.

Never store plain-text passwords.

---

# Database Standards

Use PostgreSQL.

Every table must contain:

```sql
created_at
updated_at
deleted_at
```

Soft delete should be enabled whenever possible.

Use UUID for public identifiers.

Example:

```go
ID uuid.UUID
```

Never expose internal numeric IDs through APIs.

---

# Domain Models

## User

```go
User
```

Fields:

* ID
* Name
* Email
* Password
* Role

---

## Book

```go
Book
```

Fields:

* ID
* ISBN
* Title
* Author
* Description
* Category
* Price
* Stock

---

## Cart

```go
Cart
CartItem
```

---

## Order

```go
Order
OrderItem
```

---

## InventoryLog

```go
InventoryLog
```

---

# Order Lifecycle

Allowed transitions:

```text
PENDING_PAYMENT
    ↓
PAID
    ↓
PACKING
    ↓
SHIPPED
    ↓
DELIVERED
    ↓
COMPLETED
```

Cancellation:

```text
PENDING_PAYMENT
    ↓
CANCELLED
```

Any invalid status transition must return an error.

---

# Inventory Rules

Inventory cannot become negative.

Validation:

```go
book.Stock >= requestedQty
```

Checkout must fail if stock is insufficient.

Inventory updates must be logged.

Every inventory change creates:

```go
InventoryLog
```

Types:

```text
STOCK_IN
STOCK_OUT
ADJUSTMENT
```

---

# Checkout Rules

Checkout must be wrapped inside a database transaction.

Required steps:

1. Load cart
2. Validate stock
3. Create order
4. Create order items
5. Reduce stock
6. Create inventory logs
7. Clear cart
8. Commit transaction

Rollback on any failure.

---

# API Response Standard

Success:

```json
{
  "success": true,
  "message": "success",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "validation error",
  "errors": []
}
```

Never return raw database errors to clients.

---

# Error Handling

Create domain-specific errors.

Example:

```go
var (
    ErrUserNotFound
    ErrBookNotFound
    ErrInvalidCredential
    ErrInsufficientStock
)
```

Convert domain errors into proper HTTP responses.

---

# Logging

Use structured logging with Zap.

Example:

```go
logger.Info(
    "checkout completed",
    zap.String("order_id", orderID),
    zap.String("user_id", userID),
)
```

Never log:

* Passwords
* JWT Tokens
* Refresh Tokens
* Credit Card Data

---

# Caching

Redis should be used for:

* Book list cache
* Book detail cache
* Dashboard metrics
* JWT blacklist

Cache invalidation required after:

* Book update
* Book delete
* Inventory change

---

# Repository Rules

Repositories must:

* Only perform persistence operations
* Not contain business rules
* Accept context.Context

Example:

```go
FindByID(ctx context.Context, id uuid.UUID)
```

---

# Service Rules

Services contain:

* Validation
* Business logic
* Transaction orchestration

Services should be unit-testable.

Dependencies must be injected through constructors.

Example:

```go
func NewOrderService(
    orderRepo OrderRepository,
    bookRepo BookRepository,
) *OrderService
```

---

# Handler Rules

Handlers should:

* Validate request
* Call service
* Return response

Handlers should not:

* Use GORM directly
* Contain business logic

---

# Validation Rules

Use request DTOs.

Example:

```go
type CreateBookRequest struct {
    Title string `validate:"required"`
}
```

Validation occurs before service execution.

---

# Testing Standards

Minimum coverage:

```text
80%
```

Required tests:

* Service unit tests
* Repository integration tests
* Handler tests

Critical flows:

* Login
* Register
* Checkout
* Inventory adjustment
* Order status transition

must always be tested.

---

# Security Standards

Never trust user input.

Always:

* Validate requests
* Sanitize inputs
* Use parameterized queries

Enable:

* CORS
* Rate limiting
* Request size limits

Use HTTPS in production.

---

# Swagger Standards

Every endpoint must include:

* Summary
* Description
* Request example
* Response example
* Error responses

Swagger route:

```text
/swagger/index.html
```

---

# Background Jobs

Use Asynq.

Supported jobs:

* Send order confirmation email
* Send payment notification
* Generate sales report
* Low stock alert

Background jobs must be idempotent.

---

# Event Driven Design

Publish events:

```text
UserRegistered
BookCreated
BookUpdated
OrderCreated
OrderPaid
OrderCompleted
InventoryChanged
```

Events should be published after successful transactions.

---

# Coding Standards

Use:

```bash
gofmt
goimports
golangci-lint
```

Requirements:

* No global mutable state
* Dependency injection everywhere
* Context propagation everywhere
* Small functions
* Single responsibility principle

Maximum function size:

```text
50 lines
```

Maximum file size:

```text
500 lines
```

Refactor when limits are exceeded.

---

# Performance Targets

API response:

```text
< 300ms
```

Book search:

```text
< 500ms
```

Checkout:

```text
< 2 seconds
```

Dashboard:

```text
< 1 second
```

---

# AI Agent Instructions

When generating code:

1. Follow Clean Architecture.
2. Write tests alongside new features.
3. Prefer composition over inheritance.
4. Never place business logic inside handlers.
5. Never access database outside repositories.
6. Use transactions for checkout and inventory operations.
7. Use context.Context in all repository and service methods.
8. Return domain errors instead of generic errors.
9. Add Swagger annotations for new endpoints.
10. Maintain backward compatibility whenever possible.

When unsure, prioritize maintainability, testability, and security over brevity.
