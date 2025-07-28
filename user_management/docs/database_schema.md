# Database Schema Documentation

## Users Table

The `users` table stores user information with the following structure:

```sql
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE,
    age         INTEGER NOT NULL CHECK (age >= 0 AND age <= 150),
    phone       VARCHAR(20),
    address     VARCHAR(255),
    is_active   BOOLEAN DEFAULT true,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE NULL
);

-- Indexes
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_is_active ON users(is_active);
```

### Field Descriptions

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-incrementing unique identifier |
| `name` | VARCHAR(100) | NOT NULL | User's full name (2-100 characters) |
| `email` | VARCHAR(100) | NOT NULL, UNIQUE | User's email address |
| `age` | INTEGER | NOT NULL, 0-150 | User's age |
| `phone` | VARCHAR(20) | NULLABLE | User's phone number (10-20 characters) |
| `address` | VARCHAR(255) | NULLABLE | User's address (max 255 characters) |
| `is_active` | BOOLEAN | DEFAULT true | Whether the user is active |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| `updated_at` | TIMESTAMP | DEFAULT NOW() | Last update timestamp |
| `deleted_at` | TIMESTAMP | NULLABLE | Soft delete timestamp |

### Business Rules

1. **Email Uniqueness**: Each email address can only be used once
2. **Soft Deletes**: Records are marked as deleted rather than physically removed
3. **Age Validation**: Age must be between 0 and 150
4. **Name Requirements**: Name must be between 2 and 100 characters
5. **Phone Format**: Phone number can be 10-20 characters (optional)
6. **Active Status**: Users are active by default

### Performance Considerations

- **Primary Index**: On `id` field for fast lookups
- **Unique Index**: On `email` field for email validation and login
- **Soft Delete Index**: On `deleted_at` for filtering active records
- **Status Index**: On `is_active` for filtering active users

### Sample Data

```sql
INSERT INTO users (name, email, age, phone, address) VALUES
('John Doe', 'john@example.com', 30, '1234567890', '123 Main St'),
('Jane Smith', 'jane@example.com', 25, '0987654321', '456 Oak Ave'),
('Bob Johnson', 'bob@example.com', 35, NULL, NULL);
```

### Connection Pool Settings

```go
// Database connection pool configuration
sqlDB.SetMaxIdleConns(10)    // Maximum idle connections
sqlDB.SetMaxOpenConns(100)   // Maximum open connections
```

This configuration allows for:
- Up to 100 concurrent database connections
- 10 idle connections kept alive
- Automatic connection management by GORM
