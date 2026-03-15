# WMS Backend API

Backend service for the **Warehouse Management System (WMS)** built for the marketplace integration technical assessment.

This service is responsible for:

- Synchronizing orders from the marketplace API
- Managing warehouse lifecycle operations
- Handling marketplace OAuth and signed requests
- Processing webhook updates
- Providing APIs for the frontend dashboard

The backend integrates with the **Marketplace Mock API** while exposing internal APIs used by the WMS frontend.

---

# Tech Stack

### Language

- **Go**

### Web Framework

- **Fiber v3**

### Database

- **PostgreSQL**
- **Bun ORM**

### Other Libraries

- **pgdriver** – PostgreSQL driver for Bun
- **Fiber middleware** – logging, recovery, CORS

---

# Installation

### Prerequisite

This installation process requires **Task** (Taskfile) to be installed on your machine.

If you don't have it installed yet, please follow the official installation guide:
https://taskfile.dev/docs/installation

Make sure `task` is available in your terminal before continuing with the installation steps.

### 1. Clone the repository

```bash
git clone https://github.com/Corazon-17/wms-api.git
cd wms-api
```

### 2. Copy environment file

```
cp .env.example .env
```

### 3. Install dependencies

```
go mod tidy
```

### 4. Run the infrastructure (database)

```
task infra
```

### 5. Database migration

```
task migrate
```

### 6. Run the server

```
task dev
```

The API runs at:

```
http://localhost:3000
```

---

# Architecture

The backend follows a **layered architecture** to keep responsibilities clear and maintainable.

```
internal
 ├── config
 ├── database
 ├── domain
 ├── dto
 ├── handler
 ├── helper
 ├── middleware
 ├── service
 ├── repository
 ├── model
 ├── provider
 │     └── marketplace
 └── route
```

### Layer Responsibilities

| Layer                    | Responsibility                 |
| ------------------------ | ------------------------------ |
| **handler**              | HTTP request/response handling |
| **service**              | Business logic                 |
| **repository**           | Database queries               |
| **model**                | Database models                |
| **provider/marketplace** | Marketplace API client         |
| **config**               | Environment configuration      |
| **route**                | Route registration             |

Dependency direction:

```
handler → service → repository → database
              ↓
        marketplace provider
```

---

# Authentication Strategy

The WMS backend uses a **simple internal API key authentication strategy**.

### Why

The WMS is designed as an **internal operational system**, so:

- Only trusted clients access the API
- Authentication complexity can remain minimal

### Flow

```
Frontend
   ↓
API Request
   ↓
x-api-key header
   ↓
Middleware validation
   ↓
Authorized request
```

Example:

```
x-api-key: your-secret-key
```

The key is stored in environment variables and validated through middleware.

Marketplace credentials are **never exposed to the frontend**.

---

# Database Design

The system uses PostgreSQL with these main tables:

### Orders

```
orders
```

Stores synchronized marketplace orders and WMS state.

Important fields:

- order_sn
- marketplace_status
- shipping_status
- wms_status_id
- tracking_number
- total_amount
- raw_marketplace_payload

### Order Items

```
order_items
```

Stores items belonging to orders.

Relationship:

```
orders.id → order_items.order_id
```

### WMS Status Master Table

```
wms_statuses
```

Used to enforce valid warehouse lifecycle states.

---

# Order Lifecycle

Warehouse workflow:

```
READY_TO_PICK → PICKING → PACKED → SHIPPED
```

Operations:

| Action | Transition              |
| ------ | ----------------------- |
| Pickup | READY_TO_PICK → PICKING |
| Pack   | PICKING → PACKED        |
| Ship   | PACKED → SHIPPED        |

The backend determines allowed operations via:

```
allowed_actions
```

Example:

```json
{
  "orderSN": "SHP001",
  "wmsStatus": "PICKING",
  "allowedAction": "pack"
}
```

---

# Marketplace Integration

The backend integrates with the **Marketplace Mock API**.

### OAuth Authorization

Used to authorize the WMS with a marketplace shop.

### Signed Requests

Marketplace requests require HMAC signatures.

### Token Handling

The backend manages:

- access tokens
- refresh tokens

Tokens are stored securely and refreshed when needed.

### Order Synchronization

Orders are synchronized from the marketplace.

```
Marketplace
      ↓
Sync process
      ↓
Local database
```

### Webhooks

Marketplace updates are received through:

```
POST /webhook/order-status
POST /webhook/shipping-status
```

These update the order records.

---

# Error Handling

Marketplace integrations must handle unreliable APIs.

Handled scenarios:

### 401 Unauthorized

Token expired.

Action:

- refresh token
- retry request

### 429 Rate Limit

Retry with delay.

### Random 500 Errors

Retry logic prevents synchronization failures.

---

# API Endpoints

Example endpoints:

```
POST /auth/login

GET  /api/orders
GET  /api/orders/:order_sn

POST /api/orders/:order_sn/pick
POST /api/orders/:order_sn/pack
POST /api/orders/:order_sn/ship

GET  /api/orders/summary

POST  /api/webhook/order-status
POST  /api/webhook/shipping-status
```
