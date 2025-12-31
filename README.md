# Techwave Go Backend - Grade Management System

A comprehensive Go backend for managing students, grades, and course enrollments with RESTful APIs.

## Features

- Student and Grade domain with CRUD operations
- **NEW:** Student Enrollment Management with in-memory storage
- RESTful API using Gorilla Mux
- Repository and handler patterns
- In-memory persistence for all entities (students, grades, enrollments)
- Status-based enrollment workflow (pending → active → completed)
- Unit test samples
- Ready for Copilot-driven development and demo scenarios

## Prerequisites

- **Go 1.20+** installed

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

Enrollments use **in-memory storage** following the same pattern as students and grades.

### Enrollment Model

```json
{
  "id": 1,
  "student_id": 123,
  "course_id": 456,
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

Creates a new enrollment with auto-generated ID and timestamps.

- **POST** `/api/enrollments`
- **Request Body:**
  ```json
  {
    "student_id": 123,
    "course_id": 456,
    "status": "pending"
  }
  ```
- **Response:** 201 Created
  ```json
  {
    "id": 1,
    "student_id": 123,
    "course_id": 456,
    "enrollment_date": "2024-01-15T10:30:00Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
  ```
- **Status Codes:**
  - `201` - Enrollment created successfully
  - `400` - Validation error (missing fields, invalid status)

**cURL Example:**
```sh
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 123,
    "course_id": 456,
    "status": "pending"
  }'
```

#### 2. Get Enrollment by ID

Retrieves a specific enrollment by its ID.

- **GET** `/api/enrollments/{id}`
- **Response:** 200 OK
  ```json
  {
    "id": 1,
    "student_id": 123,
    "course_id": 456,
    "enrollment_date": "2024-01-15T10:30:00Z",
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-16T14:20:00Z"
  }
  ```
- **Status Codes:**
  - `200` - Enrollment found
  - `404` - Enrollment not found

**cURL Example:**
```sh
curl http://localhost:8080/api/enrollments/1
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
      "id": "1",
      "student_id": 123,
      "course_id": 456,
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
curl http://localhost:8080/api/enrollments?student_id=123

# Get enrollments for a specific course
curl http://localhost:8080/api/enrollments?course_id=456

# Get specific enrollment for student in course
curl "http://localhost:8080/api/enrollments?student_id=123&course_id=456"
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
    "id": "1",
    "student_id": 123,
    "course_id": 456,
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
curl -X PUT http://localhost:8080/api/enrollments/1 \
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
curl -X DELETE http://localhost:8080/api/enrollments/1
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

## Testing

### Run All Tests

```sh
go test ./... -v
```

### Run Specific Test Package

```sh
# Test students and grades
go test ./tests -v -run TestCreateStudent

# Test enrollments
go test ./tests -v -run TestCreateEnrollment
```

### Test Requirements

- **All Tests:** No external dependencies - all use in-memory storage
- Tests are self-contained and can run independently

## Build the Application

```sh
go build -o techwave-app .
./techwave-app
```

## Performance Considerations

### Enrollment Operations

All operations use in-memory storage with concurrent access protection via sync.RWMutex:

- **Create:** O(1) - Direct map insertion with ID generation
- **Get by ID:** O(1) - Direct map lookup
- **Get All:** O(n) - Iterates through all enrollments
- **Filter by Student/Course:** O(n) - Linear scan with filtering
- **Update:** O(1) - Direct map update with validation
- **Delete:** O(1) - Direct map deletion
- **Stats:** O(n) - Iterates through all enrollments to aggregate

### Common Issues

#### Status Transition Errors

**Problem:** `400 Bad Request` when updating enrollment status

**Cause:** Invalid status transition attempted

**Solution:** Review status transition rules:
- ✓ pending → active, completed
- ✓ active → completed
- ✗ completed → (no transitions)
- ✗ active → pending

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
│   ├── enrollment_repository.go  (In-memory)
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

## License

This project is for educational and demonstration purposes.