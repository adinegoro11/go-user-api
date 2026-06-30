# ISP Billing System (Golang)

## Project Overview

Build a simple but production-inspired Internet Service Provider (ISP) backend using Golang.

The project should demonstrate a clean architecture while following a real ISP business process from customer registration until invoice generation.

This project is intended for learning purposes, so prioritize **clean code**, **maintainability**, and **extensibility** over unnecessary complexity.

---

# Tech Stack

- Go 1.24+
- Gin (REST API)
- GORM
- PostgreSQL
- Docker & Docker Compose
- Air (Hot Reload)
- Swagger
- golang-migrate
- Zap Logger
- Viper Configuration
- UUID
- JWT Authentication (future)
- RabbitMQ (future)
- Redis (future)

---

# Project Goals

Implement the following business flow.

```
Customer Registration
        │
        ▼
Product Selection
        │
        ▼
Welcome Email
        │
        ▼
Notify OTT Provider
        │
        ▼
Create Subscription
        │
        ▼
Generate Invoice
        │
        ▼
Send Invoice
```

---

# Business Flow

## 1. Register Customer

Endpoint

```
POST /customers
```

Input

```json
{
    "name": "John Doe",
    "email": "john@email.com",
    "phone": "08123456789",
    "address": "Jakarta"
}
```

Business Rules

- Email must be unique.
- Phone must be unique.
- Generate customer number.
- Initial status = PENDING.
- Save customer.
- Publish CustomerRegistered event.

---

## 2. Setup Product

Endpoint

```
POST /products
```

Example Products

| Code | Name | Type | Price |
|------|------|------|------:|
| INTERNET100 | Internet 100 Mbps | Monthly | 300000 |
| INTERNET300 | Internet 300 Mbps | Monthly | 500000 |
| INSTALL | Installation Fee | One Time | 150000 |
| STATICIP | Static IP | Monthly | 100000 |
| OTT_NETFLIX | Netflix OTT | Monthly | 80000 |

Business Rules

- Product code must be unique.
- Can be active/inactive.
- Monthly and One Time billing types.

---

## 3. Customer Subscribe Product

Endpoint

```
POST /subscriptions
```

Input

```json
{
    "customer_id":"uuid",
    "products":[
        "internet100",
        "install"
    ]
}
```

Business Rules

- Customer must exist.
- Product must exist.
- Create subscription.
- Status = WAITING_INSTALLATION.
- Publish SubscriptionCreated event.

---

## 4. Send Welcome Email

Triggered after CustomerRegistered event.

Responsibilities

- Send Welcome Email
- Include Customer Number
- Include Customer Name

Failure

- Save into email_logs.
- Retry later.

---

## 5. Notify OTT Provider

Triggered after CustomerRegistered.

Call

```
POST /api/customer
```

Payload

```json
{
    "customerNumber":"CUS000001",
    "name":"John Doe",
    "email":"john@email.com"
}
```

Business Rules

- Save request payload.
- Save response payload.
- Retry if failed.
- Maximum retry = 5.

---

## 6. Generate Invoice

Endpoint

```
POST /billing/generate/{customerId}
```

Business Rules

Invoice consists of

- Installation Fee
- Monthly Package
- Static IP
- OTT Package
- Tax
- Discount

Invoice Status

```
UNPAID
PAID
VOID
```

Generate Invoice Number

```
INV-202607-000001
```

---

## 7. Send Invoice

Triggered after invoice generated.

Responsibilities

- Generate PDF (future)
- Email customer
- Save email log

---

# Database Design

## customers

| Column | Type |
|---------|------|
| id | UUID |
| customer_number | varchar |
| name | varchar |
| email | varchar |
| phone | varchar |
| address | text |
| status | varchar |
| created_at | timestamp |
| updated_at | timestamp |

---

## products

| Column | Type |
|---------|------|
| id | UUID |
| code | varchar |
| name | varchar |
| billing_type | varchar |
| price | decimal |
| active | boolean |

---

## subscriptions

| Column | Type |
|---------|------|
| id | UUID |
| customer_id | UUID |
| status | varchar |
| activation_date | timestamp |

---

## subscription_products

| Column | Type |
|---------|------|
| id | UUID |
| subscription_id | UUID |
| product_id | UUID |

---

## invoices

| Column | Type |
|---------|------|
| id | UUID |
| invoice_number | varchar |
| customer_id | UUID |
| subtotal | decimal |
| tax | decimal |
| discount | decimal |
| grand_total | decimal |
| due_date | timestamp |
| status | varchar |

---

## invoice_items

| Column | Type |
|---------|------|
| id | UUID |
| invoice_id | UUID |
| product_name | varchar |
| quantity | integer |
| unit_price | decimal |
| subtotal | decimal |

---

## email_logs

| Column | Type |
|---------|------|
| id | UUID |
| recipient | varchar |
| template | varchar |
| status | varchar |
| error_message | text |
| sent_at | timestamp |

---

## ott_requests

| Column | Type |
|---------|------|
| id | UUID |
| customer_id | UUID |
| provider | varchar |
| request_payload | jsonb |
| response_payload | jsonb |
| status | varchar |
| retry_count | integer |

---

# REST APIs

## Customer

```
POST   /customers
GET    /customers
GET    /customers/:id
PUT    /customers/:id
DELETE /customers/:id
```

---

## Product

```
POST /products
GET  /products
GET  /products/:id
PUT  /products/:id
```

---

## Subscription

```
POST /subscriptions
GET  /subscriptions
GET  /subscriptions/:id
```

---

## Billing

```
POST /billing/generate/:customerId
GET  /billing/invoices
GET  /billing/invoices/:id
```

---

## Notification

```
POST /notifications/send
```

---

## OTT

```
POST /ott/customer
```

---

# Project Structure

```
isp-billing/

cmd/
    api/

configs/

docs/

migrations/

internal/

    customer/
        handler.go
        service.go
        repository.go
        model.go
        dto.go

    product/
        handler.go
        service.go
        repository.go
        model.go

    subscription/
        handler.go
        service.go
        repository.go
        model.go

    billing/
        handler.go
        service.go
        repository.go
        invoice_generator.go

    notification/
        email_service.go
        template.go

    ott/
        client.go
        service.go

    middleware/
        logger.go
        recovery.go

    database/

    router/

pkg/

    logger/
    validator/
    utils/
    response/
    event/

scripts/

docker-compose.yml

Dockerfile

go.mod

README.md
```

---

# Architecture

Follow **Clean Architecture**.

```
HTTP
    │
Handler
    │
Service
    │
Repository
    │
Database
```

Rules

- Handler contains HTTP logic only.
- Service contains business logic.
- Repository handles database access.
- DTOs are separated from entities.
- Dependency Injection should be used.
- Avoid global variables.

---

# Coding Standards

- Follow Go idioms.
- Keep functions under ~50 lines where practical.
- Use context.Context everywhere.
- Return wrapped errors.
- Avoid duplicated logic.
- Write unit-testable services.
- Log only meaningful information.
- Validate all requests.
- Never panic in business logic.

---

# Future Enhancements

- Authentication (JWT)
- RBAC
- Payment Gateway Integration
- Midtrans/Xendit
- Technician Scheduling
- Installation Workflow
- Fiber Port Management
- Network Provisioning
- Monthly Billing Cron
- Auto Invoice Generation
- RabbitMQ Event Bus
- Redis Cache
- Prometheus Metrics
- Grafana Dashboard
- OpenTelemetry Tracing
- Kubernetes Deployment
- CI/CD Pipeline
- Multi-Tenant Support

---

# Success Criteria

The application should allow a user to:

1. Register as a customer.
2. Browse available products.
3. Subscribe to one or more products.
4. Automatically receive a welcome email.
5. Automatically notify an external OTT provider.
6. Generate an invoice containing all subscribed services.
7. Send the invoice to the customer.
8. Expose Swagger documentation for all APIs.
9. Run locally using Docker Compose with a single command.
10. Be structured so future features can be added with minimal refactoring.

# Testing Requirements

All business logic must be covered by unit tests.

## Test Framework

- Go built-in `testing` package
- `testify` for assertions and mocks
- Table-driven tests whenever appropriate

## Coverage Goals

| Layer | Test Required |
|--------|---------------|
| Handler | Optional |
| Service | ✅ Required |
| Repository | Optional (prefer integration tests) |
| Utility | ✅ Required |

Target minimum coverage: **80%** for service packages.

---

## Unit Test Guidelines

- Every public service method must have corresponding tests.
- Test success and failure scenarios.
- Mock all external dependencies.
- Never access a real database in unit tests.
- Never call external HTTP APIs in unit tests.
- Tests should run with:

```bash
go test ./...
```

---

## Mocking

Use interfaces for dependencies.

Example:

```go
type CustomerRepository interface {
    Create(ctx context.Context, customer *Customer) error
}
```

Generate mocks using:

```bash
mockgen
```

or

```bash
mockery
```

---

## Example Test Cases

### Customer Service

RegisterCustomer()

- ✅ Success
- ✅ Duplicate email
- ✅ Duplicate phone
- ✅ Repository returns error
- ✅ Welcome email published
- ✅ OTT event published

---

### Billing Service

GenerateInvoice()

- ✅ Installation fee included
- ✅ Monthly package included
- ✅ Multiple products
- ✅ Tax calculated correctly
- ✅ Discount applied
- ✅ Empty subscription returns error

---

### Product Service

CreateProduct()

- ✅ Success
- ✅ Duplicate code
- ✅ Invalid billing type

---

## Folder Structure

```
internal/customer/

    service.go
    service_test.go

internal/billing/

    service.go
    service_test.go
```

---

## Testing Principles

- Keep tests deterministic.
- Avoid sleeps and timing-based assertions.
- One behavior per test.
- Prefer table-driven tests.
- Use descriptive test names.

Example:

```
TestRegisterCustomer_Success
TestRegisterCustomer_DuplicateEmail
TestGenerateInvoice_WithMultipleProducts
```