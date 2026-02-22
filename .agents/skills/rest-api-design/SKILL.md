---
name: rest-api-design
description: Design RESTful APIs following best practices for resource modeling, HTTP methods, status codes, versioning, and documentation. Use when creating new APIs, designing endpoints, or improving existing API architecture.
---

# REST API Design

## Overview

Design REST APIs that are intuitive, consistent, and follow industry best practices for resource-oriented architecture.

## When to Use

- Designing new RESTful APIs
- Creating endpoint structures
- Defining request/response formats
- Implementing API versioning
- Documenting API specifications
- Refactoring existing APIs

## Instructions

### 1. **Resource Naming**

Resource naming is only relevant when you are working with collection of data.  e.g. users, products, orders, etc.  It is not relevant when you are working with actions that do not directly relate to a resource.  e.g. login, logout, etc.

Action/methods names that have more than  one work must be hyphenated.  e.g. `reset-password`, `send-email`, etc.

```
✅ Good Resource Names (Nouns, Plural)
GET    /api/users
GET    /api/users/123
GET    /api/users/123/orders
POST   /api/products
DELETE /api/products/456

❌ Bad Resource Names (Verbs, Inconsistent)
GET    /api/getUsers
POST   /api/createProduct
GET    /api/user/123  (inconsistent singular/plural)
```

### 2. **HTTP Methods & Operations**

```http
# CRUD Operations
GET    /api/users          # List all users (Read collection)
GET    /api/users/123      # Get specific user (Read single)
POST   /api/users          # Create new user (Create)
PUT    /api/users/123      # Replace user completely (Update)
PATCH  /api/users/123      # Partial update user (Partial update)
DELETE /api/users/123      # Delete user (Delete)

# Nested Resources
GET    /api/users/123/orders       # Get user's orders
POST   /api/users/123/orders       # Create order for user
GET    /api/users/123/orders/456   # Get specific order
```

### 3. **Request Examples**

#### Creating a Resource
```http
POST /api/users
Content-Type: application/json

{
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "role": "admin"
}

Response: 201 Created
Location: /api/users/789
{
  "id": "789",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "role": "admin",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

#### Updating a Resource
```http
PATCH /api/users/789
Content-Type: application/json

{
  "firstName": "Jonathan"
}

Response: 200 OK
{
  "id": "789",
  "email": "john@example.com",
  "firstName": "Jonathan",
  "lastName": "Doe",
  "role": "admin",
  "updatedAt": "2025-01-15T11:00:00Z"
}
```

### 4. **Query Parameters**

```http
# Filtering
GET /api/products?$filter=price gt 100&$filter=category eq 'electronics'

# Sorting
GET /api/users?$order-by=createdAt desc,email asc

# Pagination
GET /api/users?$top=20&$skip=40

# Expand Selection
GET /api/users?$expand=id,email,firstName

# Search
GET /api/products?q=laptop

# Multiple filters combined
GET /api/orders?$filter=status eq 'pending'&$filter=customer eq 123&$order-by=createdAt desc&$top=50
```

See sections 11, 12, and 13 for more details on pagination, sorting, and filtering.

### 5. **Response Formats**

#### Success Response
```json
{
  "value": {
    "id": "123",
    "email": "user@example.com",
    "firstName": "John"
  },
  "ok": true,
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "version": "1.0"
  }
}
```

#### Collections with Pagination

As data grows, so do collections.
Planning for pagination is important for all services.
Therefore, when multiple pages are available, the serialization payload MUST contain the opaque URL for the next page as appropriate.
Refer to the paging guidance for more details.

Clients MUST be resilient to collection data being either paged or nonpaged for any given request.

```json
{
  "value":[
    { "id": "Item 1","price": 99.95,"sizes": null},
    { … },
    { … },
    { "id": "Item 99","price": 59.99,"sizes": null}
  ],
  "ok": true,
  "@nextLink": "{opaqueUrl}"
}
```

#### Error Response
```json
{
  "ok": false,
  "error": {
    "code": "validation-error",
    "target": "email",
    "message": "Invalid input data",
    "traceId": "abc-123-def",
    "spanId": "span-456-ghi",
    "innerError":{
        "code": "nil-value",
        "target": "email",
        "message": "Email is required",
    },
    
    "details": [
      {
        "code": "invalid-format",
        "target": "email",
        "message": "Email format is invalid"
      },
      {
        "code": "invalid-value",
        "target": "age",
        "message": "Must be at least 18"
      }
    ]
  },
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "requestId": "abc-123-def"
  }
}
```

### 6. **HTTP Status Codes**

```
Success:
200 OK              - Successful GET, PATCH, DELETE
201 Created         - Successful POST (resource created)
204 No Content      - Successful DELETE (no response body)

Client Errors:
400 Bad Request     - Invalid request format/data
401 Unauthorized    - Missing or invalid authentication
403 Forbidden       - Authenticated but not authorized
404 Not Found       - Resource doesn't exist
409 Conflict        - Resource conflict (e.g., duplicate email)
422 Unprocessable   - Validation errors
429 Too Many Requests - Rate limit exceeded

Server Errors:
500 Internal Server Error - Generic server error
503 Service Unavailable   - Temporary unavailability
```

### 7. **API Versioning**

```http
# URL Path Versioning (Recommended)
GET /api/v1/users
GET /api/v2/users

# Header Versioning
GET /api/users
Accept: application/vnd.myapi.v1+json

# Query Parameter (Not recommended)
GET /api/users?version=1
```

### 8. **Authentication & Security**

```http
# JWT Bearer Token
GET /api/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# API Key
GET /api/users
X-API-Key: your-api-key-here

# Always use HTTPS in production
https://api.example.com/v1/users
```

### 9. **Rate Limiting Headers**

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 995
X-RateLimit-Reset: 1642262400
```

### 10. **OpenAPI Documentation**

```yaml
openapi: 3.0.0
info:
  title: User API
  version: 1.0.0
  description: User management API

paths:
  /users:
    get:
      summary: List all users
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/User'

    post:
      summary: Create a new user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserInput'
      responses:
        '201':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '400':
          description: Invalid input
        '409':
          description: Email already exists

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
        email:
          type: string
          format: email
        firstName:
          type: string
        lastName:
          type: string
        createdAt:
          type: string
          format: date-time

    UserInput:
      type: object
      required:
        - email
        - firstName
        - lastName
      properties:
        email:
          type: string
          format: email
        firstName:
          type: string
        lastName:
          type: string
```

11. **Pagination**

ESTful APIs that return collections MAY return partial sets.
Consumers of these services MUST expect partial result sets and correctly page through to retrieve an entire set.

There are two forms of pagination that MAY be supported by RESTful APIs.
Server-driven paging mitigates against denial-of-service attacks by forcibly paginating a request over multiple response payloads.
Client-driven paging enables clients to request only the number of resources that it can use at a given time.

Sorting and Filtering parameters MUST be consistent across pages, because both client- and server-side paging is fully compatible with both filtering and sorting.

#### 111.1 Server-driven paging
Paginated responses MUST indicate a partial result by including a continuation token in the response.
The absence of a continuation token means that no additional pages are available.

Clients MUST treat the continuation URL as opaque, which means that query options may not be changed while iterating over a set of partial results.

Example:

```http
GET http://api.contoso.com/v1.0/people HTTP/1.1
Accept: application/json

HTTP/1.1 200 OK
Content-Type: application/json

{
  ...,
  "value": [...],
  "@nextLink": "{opaqueUrl}"
}
```

#### 11.2 Client-driven paging
Clients MAY use _$top_ and _$skip_ query parameters to specify a number of results to return and an offset.
The server SHOULD honor the values specified by the client; however, clients MUST be prepared to handle responses that contain a different page size or contain a continuation token.

Note: If the server can't honor _$top_ and/or _$skip_, the server MUST return an error to the client informing about it instead of just ignoring the query options.
This will avoid the risk of the client making assumptions about the data returned.

Example:

```http
GET http://api.contoso.com/v1.0/people?$top=5&$skip=2 HTTP/1.1
Accept: application/json

HTTP/1.1 200 OK
Content-Type: application/json

{
  ...,
  "value": [...]
}
```

### 12 Sorting

The results of a collection query MAY be sorted based on property values.
The property is determined by the value of the _$orderBy_ query parameter.

The value of the _$orderBy_ parameter contains a comma-separated list of expressions used to sort the items.
A special case of such an expression is a property path terminating on a primitive property.

The expression MAY include the suffix "asc" for ascending or "desc" for descending, separated from the property name by one or more spaces.
If "asc" or "desc" is not specified, the service MUST order by the specified property in ascending order.

NULL values MUST sort as "less than" non-NULL values.

Items MUST be sorted by the result values of the first expression, and then items with the same value for the first expression are sorted by the result value of the second expression, and so on.
The sort order is the inherent order for the type of the property.

For example:

```http
GET https://api.contoso.com/v1.0/people?$orderBy=name
```

Will return all people sorted by name in ascending order.

For example:

```http
GET https://api.contoso.com/v1.0/people?$orderBy=name desc
```

Will return all people sorted by name in descending order.

Sub-sorts can be specified by a comma-separated list of property names with OPTIONAL direction qualifier.

For example:

```http
GET https://api.contoso.com/v1.0/people?$orderBy=name desc,hireDate
```

Will return all people sorted by name in descending order and a secondary sort order of hireDate in ascending order.

Sorting MUST compose with filtering such that:

```http
GET https://api.contoso.com/v1.0/people?$filter=name eq 'david'&$orderBy=hireDate
```

Will return all people whose name is David sorted in ascending order by hireDate.

#### 12.1 Interpreting a sorting expression
Sorting parameters MUST be consistent across pages, as both client and server-side paging is fully compatible with sorting.

If a service does not support sorting by a property named in a _$orderBy_ expression, the service MUST respond with an error message as defined in the Responding to Unsupported Requests section.

## 13 Filtering
The _$filter_ querystring parameter allows clients to filter a collection of resources that are addressed by a request URL.
The expression specified with _$filter_ is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response.
Resources for which the expression evaluates to false or to null, or which reference properties that are unavailable due to permissions, are omitted from the response.

Example: return all Products whose Price is less than $10.00

```http
GET https://api.contoso.com/v1.0/products?$filter=price lt 10.00
```

The value of the _$filter_ option is a Boolean expression.

#### 13.1 Filter operations
Services that support _$filter_ SHOULD support the following minimal set of operations.

Operator             | Description           | Example
-------------------- | --------------------- | -----------------------------------------------------
Comparison Operators |                       |
eq                   | Equal                 | city eq 'Redmond'
ne                   | Not equal             | city ne 'London'
gt                   | Greater than          | price gt 20
ge                   | Greater than or equal | price ge 10
lt                   | Less than             | price lt 20
le                   | Less than or equal    | price le 100
Logical Operators    |                       |
and                  | Logical and           | price le 200 and price gt 3.5
or                   | Logical or            | price le 3.5 or price gt 200
not                  | Logical negation      | not price le 3.5
Grouping Operators   |                       |
( )                  | Precedence grouping   | (priority eq 1 or city eq 'Redmond') and price gt 100

#### 13.2 Operator examples
The following examples illustrate the use and semantics of each of the logical operators.

Example: all products with a name equal to 'Milk'

```http
GET https://api.contoso.com/v1.0/products?$filter=name eq 'Milk'
```

Example: all products with a name not equal to 'Milk'

```http
GET https://api.contoso.com/v1.0/products?$filter=name ne 'Milk'
```

Example: all products with the name 'Milk' that also have a price less than 2.55:

```http
GET https://api.contoso.com/v1.0/products?$filter=name eq 'Milk' and price lt 2.55
```

Example: all products that either have the name 'Milk' or have a price less than 2.55:

```http
GET https://api.contoso.com/v1.0/products?$filter=name eq 'Milk' or price lt 2.55
```

Example 54: all products that have the name 'Milk' or 'Eggs' and have a price less than 2.55:

```http
GET https://api.contoso.com/v1.0/products?$filter=(name eq 'Milk' or name eq 'Eggs') and price lt 2.55
```

#### 13.3 Operator precedence
Services MUST use the following operator precedence for supported operators when evaluating _$filter_ expressions.
Operators are listed by category in order of precedence from highest to lowest.
Operators in the same category have equal precedence:

Group           | Operator | Description
--------------- | -------- | ---------------------
Grouping        | ( )      | Precedence grouping
Unary           | not      | Logical Negation
Relational      | gt       | Greater Than
                | ge       | Greater than or Equal
                | lt       | Less Than
                | le       | Less than or Equal
Equality        | eq       | Equal
                | ne       | Not Equal
Conditional AND | and      | Logical And
Conditional OR  | or       | Logical Or

## Best Practices

### ✅ DO
- Use nouns for resources, not verbs
- Use plural names for collections
- Be consistent with naming conventions
- Return appropriate HTTP status codes
- Include pagination for collections
- Provide filtering and sorting options
- Version your API
- Document thoroughly with OpenAPI
- Use HTTPS
- Implement rate limiting
- Provide clear error messages
- Use ISO 8601 for dates

### ❌ DON'T
- Use verbs in endpoint names
- Return 200 for errors
- Expose internal IDs unnecessarily
- Over-nest resources (max 2 levels)
- Use inconsistent naming
- Forget authentication
- Return sensitive data
- Break backward compatibility without versioning

## Complete Example: Express.js

```javascript
const express = require('express');
const app = express();

app.use(express.json());

// List users with pagination
app.get('/api/v1/users', async (req, res) => {
  try {
    const page = parseInt(req.query.page) || 1;
    const limit = parseInt(req.query.limit) || 20;
    const offset = (page - 1) * limit;

    const users = await User.findAndCountAll({
      limit,
      offset,
      attributes: ['id', 'email', 'firstName', 'lastName']
    });

    res.json({
      data: users.rows,
      pagination: {
        page,
        limit,
        total: users.count,
        totalPages: Math.ceil(users.count / limit)
      }
    });
  } catch (error) {
    res.status(500).json({
      error: {
        code: 'INTERNAL_ERROR',
        message: 'An error occurred while fetching users'
      }
    });
  }
});

// Get single user
app.get('/api/v1/users/:id', async (req, res) => {
  try {
    const user = await User.findByPk(req.params.id);

    if (!user) {
      return res.status(404).json({
        error: {
          code: 'NOT_FOUND',
          message: 'User not found'
        }
      });
    }

    res.json({ data: user });
  } catch (error) {
    res.status(500).json({
      error: {
        code: 'INTERNAL_ERROR',
        message: 'An error occurred'
      }
    });
  }
});

// Create user
app.post('/api/v1/users', async (req, res) => {
  try {
    const { email, firstName, lastName } = req.body;

    // Validation
    if (!email || !firstName || !lastName) {
      return res.status(400).json({
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Missing required fields',
          details: [
            !email && { field: 'email', message: 'Email is required' },
            !firstName && { field: 'firstName', message: 'First name is required' },
            !lastName && { field: 'lastName', message: 'Last name is required' }
          ].filter(Boolean)
        }
      });
    }

    const user = await User.create({ email, firstName, lastName });

    res.status(201)
       .location(`/api/v1/users/${user.id}`)
       .json({ data: user });
  } catch (error) {
    if (error.name === 'SequelizeUniqueConstraintError') {
      return res.status(409).json({
        error: {
          code: 'CONFLICT',
          message: 'Email already exists'
        }
      });
    }
    res.status(500).json({
      error: {
        code: 'INTERNAL_ERROR',
        message: 'An error occurred'
      }
    });
  }
});

app.listen(3000);
```
