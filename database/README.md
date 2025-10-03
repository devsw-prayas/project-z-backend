# Database Setup Instructions

This document provides instructions for setting up the PostgreSQL database for the coding evaluation platform.

## Prerequisites

- PostgreSQL 12 or higher
- Go 1.19 or higher

## Database Setup

### 1. Install PostgreSQL

**On Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

**On macOS:**
```bash
brew install postgresql
brew services start postgresql
```

**On Windows:**
Download and install from [PostgreSQL official website](https://www.postgresql.org/download/windows/)

### 2. Create Database and User

```bash
# Connect to PostgreSQL as superuser
sudo -u postgres psql

# Create database
CREATE DATABASE database;

# Create user (optional - you can use default postgres user)
CREATE USER Parth WITH PASSWORD 'pass1';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE database TO Parth;

# Exit PostgreSQL
\q
```

### 3. Run Migration

Navigate to the project root and run the migration script:

```bash
# Using psql command
psql -h localhost -U postgres -d coding_platform -f database/migrate.sql

# Or if you created a custom user
psql -h localhost -U coding_user -d coding_platform -f database/migrate.sql
```

### 4. Environment Configuration

1. Copy the environment example file:
```bash
cp .env.example .env
```

2. Update the `.env` file with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres  # or coding_user if you created one
DB_PASSWORD=your_password_here
DB_NAME=coding_platform
DB_SSLMODE=disable
```

### 5. Required Go Dependencies

Add these dependencies to your `go.mod`:

```bash
go mod tidy
go get github.com/lib/pq
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
```

## Database Schema Overview

### Core Tables

1. **users** - User accounts and authentication
2. **problems** - Coding problems/challenges
3. **test_cases** - Test cases for each problem
4. **submissions** - User code submissions
5. **submission_test_results** - Results for individual test cases
6. **user_sessions** - Authentication sessions
7. **user_problem_attempts** - User attempt tracking

### Key Features

- **Auto-incrementing IDs** for all primary keys
- **Foreign key constraints** to maintain data integrity
- **Indexes** for optimal query performance
- **Triggers** for automatic timestamp updates
- **Check constraints** for data validation
- **Sample data** included for testing

### Security Considerations

- Passwords are stored as hashes (use bcrypt)
- JWT tokens are stored as hashes in sessions table
- User roles for authorization
- IP address tracking for sessions

## Testing the Setup

You can test the database connection using the provided Go database package:

```go
package main

import (
    "log"
    "your-project/src/database"
)

func main() {
    db, err := database.InitDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.Close(db)
    
    log.Println("Database connection successful!")
}
```

## Sample Queries

### Get all problems:
```sql
SELECT id, title, difficulty, total_submissions, accepted_submissions 
FROM problems 
WHERE is_public = true 
ORDER BY created_at DESC;
```

### Get user submission history:
```sql
SELECT s.id, p.title, s.language, s.status, s.score, s.submitted_at
FROM submissions s
JOIN problems p ON s.problem_id = p.id
WHERE s.user_id = $1
ORDER BY s.submitted_at DESC;
```

### Get problem statistics:
```sql
SELECT 
    p.title,
    p.difficulty,
    p.total_submissions,
    p.accepted_submissions,
    CASE 
        WHEN p.total_submissions > 0 
        THEN ROUND((p.accepted_submissions::decimal / p.total_submissions * 100), 2)
        ELSE 0 
    END as acceptance_rate
FROM problems p
WHERE p.is_public = true;
```

## Backup and Maintenance

### Create backup:
```bash
pg_dump -h localhost -U postgres coding_platform > backup.sql
```

### Restore from backup:
```bash
psql -h localhost -U postgres coding_platform < backup.sql
```

### Monitor database size:
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```