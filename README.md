# Techwave Go Backend - Grade Management System

A comprehensive Go backend for managing students, grades, and course enrollments with RESTful APIs.

## Features

- Student and Grade domain with CRUD operations
- **NEW:** Student Enrollment Management with Redis storage
- RESTful API using Gorilla Mux
- Repository and handler patterns
- Redis-backed persistence for enrollments
- Status-based enrollment workflow (pending → active → completed)
- Unit test samples
- Ready for Copilot-driven development and demo scenarios

## Prerequisites

- **Go 1.20+** installed
- **Redis** server running (for enrollment features)
- **Docker** (optional, for running Redis in a container)

## Installation and Setup

### 1. Clone the repository

```sh
git clone https://github.com/Hemavathi15sg/grademanagement_techwave.git
cd grademanagement_techwave
```

### 2. Install dependencies

```sh
go mod download
```

### 3. Start Redis Server

**Option A: Using Docker**
```sh
docker run -d -p 6379:6379 redis:latest
```

**Option B: Using local Redis installation**
```sh
# Ubuntu/Debian
sudo apt-get install redis-server
sudo service redis-server start

# macOS with Homebrew
brew install redis
brew services start redis
```

### 4. Configure Environment Variables

The application supports the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_ADDR` | Redis server address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password (if required) | Empty |

**Example:**
```sh
export REDIS_ADDR="localhost:6379"
export REDIS_PASSWORD=""
```

## Run the Server

```sh
go run main.go
```

The server will start on **http://localhost:8080**

## API Documentation

### Student Endpoints

Students use in-memory storage.

#### Create Student
- **POST** `/students`
- **Request Body:**
  ```json
  {
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "grade": "A"
  }
  ```
- **Response:** 201 Created
  ```json
  {
    "id": 1,
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "grade": "A"
  }
  ```

#### Get All Students
- **GET** `/students`
- **Response:** 200 OK - Array of students

#### Get Student by ID
- **GET** `/students/{id}`
- **Response:** 200 OK or 404 Not Found

#### Update Student
- **PUT** `/students/{id}`
- **Request Body:** Student object
- **Response:** 200 OK or 404 Not Found

#### Delete Student
- **DELETE** `/students/{id}`
- **Response:** 204 No Content or 404 Not Found

### Grade Endpoints

Grades use in-memory storage.

#### Create Grade
- **POST** `/grades`
- **Request Body:**
  ```json
  {
    "student_id": 1,
    "value": "A+",
    "subject": "Mathematics"
  }
  ```
- **Response:** 201 Created

#### Get All Grades
- **GET** `/grades`
- **Response:** 200 OK - Array of grades

#### Get Grade by ID
- **GET** `/grades/{id}`
- **Response:** 200 OK or 404 Not Found

#### Update Grade
- **PUT** `/grades/{id}`
- **Request Body:** Grade object
- **Response:** 200 OK or 404 Not Found

#### Delete Grade
- **DELETE** `/grades/{id}`
- **Response:** 204 No Content or 404 Not Found

---

## Enrollment Management API

Enrollments use **Redis** for persistence with advanced querying capabilities.

### Enrollment Model

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "student_id": "student-123",
  "course_id": "course-456",
  "enrollment_date": "2024-01-15T10:30:00Z",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Enrollment Status Values

The enrollment system supports three status values with enforced transitions:

| Status | Description | Can Transition To |
|--------|-------------|-------------------|
| `pending` | Enrollment awaiting approval or processing | `active`, `completed` |
| `active` | Student is actively enrolled in the course | `completed` |
| `completed` | Student has completed the course (final state) | _(none - final state)_ |

**Status Transition Rules:**
- **pending** → **active** or **completed** ✓
- **active** → **completed** ✓
- **completed** → _(no transitions allowed)_ ✗
- Attempting invalid transitions returns **400 Bad Request**

### Enrollment Endpoints

#### 1. Create Enrollment

Creates a new enrollment with auto-generated UUID and timestamps.

- **POST** `/api/enrollments`
- **Request Body:**
  ```json
  {
    "student_id": "student-123",
    "course_id": "course-456",
    "status": "pending"
  }
  ```
- **Response:** 201 Created
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "student-123",
    "course_id": "course-456",
    "enrollment_date": "2024-01-15T10:30:00Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
  ```
- **Status Codes:**
  - `201` - Enrollment created successfully
  - `400` - Validation error (missing fields, invalid status)
  - `500` - Internal server error

**cURL Example:**
```sh
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "student-123",
    "course_id": "course-456",
    "status": "pending"
  }'
```

#### 2. Get Enrollment by ID

Retrieves a specific enrollment by its UUID.

- **GET** `/api/enrollments/{id}`
- **Response:** 200 OK
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "student-123",
    "course_id": "course-456",
    "enrollment_date": "2024-01-15T10:30:00Z",
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-16T14:20:00Z"
  }
  ```
- **Status Codes:**
  - `200` - Enrollment found
  - `404` - Enrollment not found
  - `500` - Internal server error

**cURL Example:**
```sh
curl http://localhost:8080/api/enrollments/550e8400-e29b-41d4-a716-446655440000
```

#### 3. Get All Enrollments (with filtering)

Retrieves enrollments with optional query parameters for filtering.

- **GET** `/api/enrollments`
- **Query Parameters:**
  - `student_id` - Filter by student ID
  - `course_id` - Filter by course ID
  - Both parameters can be combined to get a specific enrollment

- **Response:** 200 OK - Array of enrollments
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "student_id": "student-123",
      "course_id": "course-456",
      "enrollment_date": "2024-01-15T10:30:00Z",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
  ```
- **Status Codes:**
  - `200` - Success (returns empty array if no matches)
  - `500` - Internal server error

**cURL Examples:**
```sh
# Get all enrollments
curl http://localhost:8080/api/enrollments

# Get enrollments for a specific student
curl http://localhost:8080/api/enrollments?student_id=student-123

# Get enrollments for a specific course
curl http://localhost:8080/api/enrollments?course_id=course-456

# Get specific enrollment for student in course
curl "http://localhost:8080/api/enrollments?student_id=student-123&course_id=course-456"
```

#### 4. Update Enrollment

Updates an existing enrollment with status transition validation.

- **PUT** `/api/enrollments/{id}`
- **Request Body:** (all fields optional, only provided fields are updated)
  ```json
  {
    "status": "active"
  }
  ```
- **Response:** 200 OK
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "student-123",
    "course_id": "course-456",
    "enrollment_date": "2024-01-15T10:30:00Z",
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-16T14:20:00Z"
  }
  ```
- **Status Codes:**
  - `200` - Enrollment updated successfully
  - `400` - Validation error (invalid status transition)
  - `404` - Enrollment not found
  - `500` - Internal server error

**cURL Example:**
```sh
curl -X PUT http://localhost:8080/api/enrollments/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"status": "active"}'
```

#### 5. Delete Enrollment

Removes an enrollment and all its index entries from Redis.

- **DELETE** `/api/enrollments/{id}`
- **Response:** 204 No Content
- **Status Codes:**
  - `204` - Enrollment deleted successfully
  - `404` - Enrollment not found
  - `500` - Internal server error

**cURL Example:**
```sh
curl -X DELETE http://localhost:8080/api/enrollments/550e8400-e29b-41d4-a716-446655440000
```

#### 6. Get Enrollment Statistics

Returns aggregate statistics about enrollments grouped by status.

- **GET** `/api/enrollments/stats`
- **Response:** 200 OK
  ```json
  {
    "total_enrollments": 10,
    "by_status": {
      "pending": 3,
      "active": 5,
      "completed": 2
    }
  }
  ```
- **Status Codes:**
  - `200` - Success
  - `500` - Internal server error

**cURL Example:**
```sh
curl http://localhost:8080/api/enrollments/stats
```

### Error Response Format

All endpoints return errors in a consistent format:

```json
HTTP/1.1 400 Bad Request
Content-Type: text/plain

status must be one of: pending, active, completed
```

**Common Error Scenarios:**

| Status Code | Scenario |
|-------------|----------|
| `400` | Missing required fields (student_id, course_id, status) |
| `400` | Invalid status value (not pending/active/completed) |
| `400` | Invalid status transition (e.g., completed → active) |
| `404` | Enrollment ID not found |
| `500` | Redis connection error or internal server error |

## Redis Data Structure

Enrollments are stored in Redis using the following key structure:

### Primary Data Keys
- **`enrollment:{id}`** - Stores enrollment JSON data
  ```
  enrollment:550e8400-e29b-41d4-a716-446655440000 → {"id":"550e8400-...", ...}
  ```

### Index Sets (for fast queries)
- **`enrollments:all`** - Set of all enrollment IDs
- **`student:enrollments:{student_id}`** - Set of enrollment IDs for a student
- **`course:enrollments:{course_id}`** - Set of enrollment IDs for a course

### Example Redis Operations

View all enrollment IDs:
```sh
redis-cli SMEMBERS enrollments:all
```

View enrollments for a specific student:
```sh
redis-cli SMEMBERS student:enrollments:student-123
```

View a specific enrollment:
```sh
redis-cli GET enrollment:550e8400-e29b-41d4-a716-446655440000
```

## Testing

### Run All Tests

```sh
go test ./... -v
```

### Run Specific Test Package

```sh
# Test students and grades
go test ./tests -v -run TestCreateStudent

# Test enrollments (requires Redis)
go test ./tests -v -run TestCreateEnrollment
```

### Test Requirements

- **Student/Grade Tests:** No external dependencies (in-memory)
- **Enrollment Tests:** Requires Redis server running on `localhost:6379`
  - Tests use database 15 to avoid conflicts with production data
  - Database is flushed before each test

## Build the Application

```sh
go build -o techwave-app .
./techwave-app
```

## Troubleshooting

### Redis Connection Error

**Problem:** `connection refused` or `cannot connect to Redis`

**Solutions:**
1. Verify Redis is running:
   ```sh
   redis-cli ping
   # Should return: PONG
   ```

2. Check Redis address configuration:
   ```sh
   export REDIS_ADDR="localhost:6379"
   ```

3. Start Redis if not running:
   ```sh
   # Docker
   docker run -d -p 6379:6379 redis:latest
   
   # Or use your system's service manager
   sudo service redis-server start
   ```

### Enrollment Test Failures

**Problem:** Tests fail with Redis errors

**Solutions:**
1. Ensure Redis is running before running tests
2. Check that port 6379 is available
3. Verify firewall settings allow local Redis connections

### Status Transition Errors

**Problem:** `400 Bad Request` when updating enrollment status

**Cause:** Invalid status transition attempted

**Solution:** Review status transition rules:
- ✓ pending → active, completed
- ✓ active → completed
- ✗ completed → (no transitions)
- ✗ active → pending

## Performance Considerations

### Enrollment Operations

- **Create:** O(1) - Uses Redis SET and SADD operations
- **Get by ID:** O(1) - Direct Redis GET
- **Get All:** O(n) - Retrieves all enrollments
- **Filter by Student/Course:** O(n) - Set membership lookup + data fetch
- **Update:** O(1) - Single Redis SET operation
- **Delete:** O(1) - Uses Redis pipeline for atomic deletion
- **Stats:** O(n) - Fetches all enrollments and aggregates in memory

### Redis Best Practices

1. **Indexing:** Multiple indexes maintained for fast queries
2. **Atomic Operations:** Pipeline/transactions ensure data consistency
3. **Key Expiration:** Currently disabled (0 TTL) - add if needed for temporary enrollments
4. **Connection Pooling:** go-redis handles connection pooling automatically

## Project Structure

```
.
├── handlers/           # HTTP request handlers
│   ├── enrollment_handler.go
│   ├── grade_handler.go
│   └── student_handler.go
├── models/            # Data models
│   ├── enrollment.go
│   ├── grade.go
│   └── student.go
├── repository/        # Data access layer
│   ├── enrollment_repository.go  (Redis)
│   ├── grade_repository.go       (In-memory)
│   └── student_repository.go     (In-memory)
├── routes/            # Route registration
│   └── routes.go
├── tests/             # Unit tests
│   ├── enrollment_test.go
│   ├── grade_test.go
│   └── student_test.go
├── main.go           # Application entry point
├── go.mod            # Go module dependencies
└── README.md         # This file
```

## Dependencies

- **github.com/gorilla/mux** v1.8.0 - HTTP router
- **github.com/go-redis/redis/v8** v8.11.5 - Redis client
- **github.com/google/uuid** v1.6.0 - UUID generation

## License

This project is for educational and demonstration purposes.