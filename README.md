# Invoice API Documentation

## Overview

REST API untuk mengelola invoice dengan autentikasi JWT. API ini menyediakan operasi CRUD (Create, Read, Update, Delete) untuk invoice dan item-itemnya.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Semua endpoint yang dilindungi memerlukan JWT token di header:

```
Authorization: Bearer <your_jwt_token>
```

---

## Endpoints

### 1. Create Invoice

**POST** `/invoices`

Membuat invoice baru untuk user yang sedang login.

#### Headers

```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

#### Request Body

```json
{
  "invoiceNumber": "INV-001",
  "date": "2024-12-01",
  "fromName": "PT. ABC Company",
  "fromEmail": "billing@abc.com",
  "toName": "PT. XYZ Client",
  "toEmail": "finance@xyz.com",
  "taxRate": 10.0,
  "subtotal": 1000000.0,
  "taxAmount": 100000.0,
  "total": 1100000.0,
  "items": [
    {
      "description": "Web Development Service",
      "quantity": 1,
      "unitPrice": 500000.0,
      "total": 500000.0
    },
    {
      "description": "SEO Optimization",
      "quantity": 1,
      "unitPrice": 500000.0,
      "total": 500000.0
    }
  ]
}
```

#### Response

**Status: 201 Created**

```json
{
  "id": 1,
  "userID": 123,
  "number": "INV-001",
  "date": "2024-12-01",
  "fromName": "PT. ABC Company",
  "fromEmail": "billing@abc.com",
  "toName": "PT. XYZ Client",
  "toEmail": "finance@xyz.com",
  "taxRate": 10.0,
  "subtotal": 1000000.0,
  "taxAmount": 100000.0,
  "total": 1100000.0,
  "items": [
    {
      "id": 1,
      "invoiceID": 1,
      "description": "Web Development Service",
      "quantity": 1,
      "unitPrice": 500000.0,
      "total": 500000.0
    },
    {
      "id": 2,
      "invoiceID": 1,
      "description": "SEO Optimization",
      "quantity": 1,
      "unitPrice": 500000.0,
      "total": 500000.0
    }
  ],
  "createdAt": "2024-12-01T10:00:00Z",
  "updatedAt": "2024-12-01T10:00:00Z"
}
```

---

### 2. Get All Invoices

**GET** `/invoices`

Mengambil semua invoice milik user yang sedang login.

#### Headers

```
Authorization: Bearer <jwt_token>
```

#### Response

**Status: 200 OK**

```json
[
  {
    "id": 1,
    "userID": 123,
    "number": "INV-001",
    "date": "2024-12-01",
    "fromName": "PT. ABC Company",
    "fromEmail": "billing@abc.com",
    "toName": "PT. XYZ Client",
    "toEmail": "finance@xyz.com",
    "taxRate": 10.0,
    "subtotal": 1000000.0,
    "taxAmount": 100000.0,
    "total": 1100000.0,
    "items": [
      {
        "id": 1,
        "invoiceID": 1,
        "description": "Web Development Service",
        "quantity": 1,
        "unitPrice": 500000.0,
        "total": 500000.0
      }
    ],
    "createdAt": "2024-12-01T10:00:00Z",
    "updatedAt": "2024-12-01T10:00:00Z"
  }
]
```

---

### 3. Get Invoice by ID

**GET** `/invoices/{id}`

Mengambil detail invoice berdasarkan ID.

#### Headers

```
Authorization: Bearer <jwt_token>
```

#### Path Parameters

- `id` (integer): ID invoice yang akan diambil

#### Response

**Status: 200 OK**

```json
{
  "id": 1,
  "userID": 123,
  "number": "INV-001",
  "date": "2024-12-01",
  "fromName": "PT. ABC Company",
  "fromEmail": "billing@abc.com",
  "toName": "PT. XYZ Client",
  "toEmail": "finance@xyz.com",
  "taxRate": 10.0,
  "subtotal": 1000000.0,
  "taxAmount": 100000.0,
  "total": 1100000.0,
  "items": [
    {
      "id": 1,
      "invoiceID": 1,
      "description": "Web Development Service",
      "quantity": 1,
      "unitPrice": 500000.0,
      "total": 500000.0
    }
  ],
  "createdAt": "2024-12-01T10:00:00Z",
  "updatedAt": "2024-12-01T10:00:00Z"
}
```

**Status: 404 Not Found**

```json
{
  "error": "invoice not found"
}
```

---

### 4. Update Invoice

**PUT** `/invoices/{id}`

Mengupdate invoice berdasarkan ID. Semua item lama akan diganti dengan item baru.

#### Headers

```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

#### Path Parameters

- `id` (integer): ID invoice yang akan diupdate

#### Request Body

```json
{
  "invoiceNumber": "INV-001-UPDATED",
  "date": "2024-12-01",
  "fromName": "PT. ABC Company Updated",
  "fromEmail": "billing@abc.com",
  "toName": "PT. XYZ Client",
  "toEmail": "finance@xyz.com",
  "taxRate": 11.0,
  "subtotal": 1200000.0,
  "taxAmount": 132000.0,
  "total": 1332000.0,
  "items": [
    {
      "description": "Web Development Service - Updated",
      "quantity": 1,
      "unitPrice": 600000.0,
      "total": 600000.0
    },
    {
      "description": "Mobile App Development",
      "quantity": 1,
      "unitPrice": 600000.0,
      "total": 600000.0
    }
  ]
}
```

#### Response

**Status: 200 OK**

```json
{
  "id": 1,
  "userID": 123,
  "number": "INV-001-UPDATED",
  "date": "2024-12-01",
  "fromName": "PT. ABC Company Updated",
  "fromEmail": "billing@abc.com",
  "toName": "PT. XYZ Client",
  "toEmail": "finance@xyz.com",
  "taxRate": 11.0,
  "subtotal": 1200000.0,
  "taxAmount": 132000.0,
  "total": 1332000.0,
  "items": [
    {
      "id": 3,
      "invoiceID": 1,
      "description": "Web Development Service - Updated",
      "quantity": 1,
      "unitPrice": 600000.0,
      "total": 600000.0
    },
    {
      "id": 4,
      "invoiceID": 1,
      "description": "Mobile App Development",
      "quantity": 1,
      "unitPrice": 600000.0,
      "total": 600000.0
    }
  ],
  "createdAt": "2024-12-01T10:00:00Z",
  "updatedAt": "2024-12-01T11:00:00Z"
}
```

**Status: 404 Not Found**

```json
{
  "error": "invoice not found"
}
```

---

### 5. Delete Invoice

**DELETE** `/invoices/{id}`

Menghapus invoice berdasarkan ID.

#### Headers

```
Authorization: Bearer <jwt_token>
```

#### Path Parameters

- `id` (integer): ID invoice yang akan dihapus

#### Response

**Status: 200 OK**

```json
{
  "message": "invoice deleted"
}
```

**Status: 404 Not Found**

```json
{
  "error": "invoice not found"
}
```

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "invalid request body or parameters"
}
```

### 401 Unauthorized

```json
{
  "error": "invalid token"
}
```

### 404 Not Found

```json
{
  "error": "invoice not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "internal server error"
}
```

---

## Data Models

### Invoice

```json
{
  "id": "integer",
  "userID": "integer",
  "number": "string",
  "date": "string (YYYY-MM-DD)",
  "fromName": "string",
  "fromEmail": "string",
  "toName": "string",
  "toEmail": "string",
  "taxRate": "float64",
  "subtotal": "float64",
  "taxAmount": "float64",
  "total": "float64",
  "items": "array of InvoiceItem",
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)"
}
```

### InvoiceItem

```json
{
  "id": "integer",
  "invoiceID": "integer",
  "description": "string",
  "quantity": "integer",
  "unitPrice": "float64",
  "total": "float64"
}
```

---

## Authentication Endpoints

### Register

**POST** `/register`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Login

**POST** `/login`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 123,
    "email": "user@example.com"
  }
}
```

---

## Usage Examples

### cURL Examples

#### 1. Login to get token

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### 2. Create invoice

```bash
curl -X POST http://localhost:8080/api/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "invoiceNumber": "INV-001",
    "date": "2024-12-01",
    "fromName": "PT. ABC Company",
    "fromEmail": "billing@abc.com",
    "toName": "PT. XYZ Client",
    "toEmail": "finance@xyz.com",
    "taxRate": 10.0,
    "subtotal": 1000000.0,
    "taxAmount": 100000.0,
    "total": 1100000.0,
    "items": [
      {
        "description": "Web Development Service",
        "quantity": 1,
        "unitPrice": 500000.0,
        "total": 500000.0
      }
    ]
  }'
```

#### 3. Get all invoices

```bash
curl -X GET http://localhost:8080/api/invoices \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 4. Get invoice by ID

```bash
curl -X GET http://localhost:8080/api/invoices/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 5. Update invoice

```bash
curl -X PUT http://localhost:8080/api/invoices/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "invoiceNumber": "INV-001-UPDATED",
    "date": "2024-12-01",
    "fromName": "PT. ABC Company Updated",
    "fromEmail": "billing@abc.com",
    "toName": "PT. XYZ Client",
    "toEmail": "finance@xyz.com",
    "taxRate": 11.0,
    "subtotal": 1200000.0,
    "taxAmount": 132000.0,
    "total": 1332000.0,
    "items": [
      {
        "description": "Web Development Service - Updated",
        "quantity": 1,
        "unitPrice": 600000.0,
        "total": 600000.0
      }
    ]
  }'
```

#### 6. Delete invoice

```bash
curl -X DELETE http://localhost:8080/api/invoices/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Notes

1. **Date Format**: Semua tanggal menggunakan format `YYYY-MM-DD`
2. **Currency**: Semua nilai mata uang dalam format float64 (contoh: 1000000.0)
3. **Authentication**: Token JWT expired sesuai konfigurasi server
4. **User Isolation**: Setiap user hanya bisa mengakses invoice miliknya sendiri
5. **Transaction**: Update invoice menggunakan database transaction untuk konsistensi data
