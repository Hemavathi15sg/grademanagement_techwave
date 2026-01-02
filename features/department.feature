Feature: Department Management System
  As an academic administrator
  I want to manage university departments through a REST API
  So that I can maintain department information, budgets, and organizational structure

  Background:
    Given the Department Management API is running
    And the Redis cache is available
    And the database is initialized

  # CRUD Operations - Create

  Scenario: Create a new department with valid data
    Given I have department data with code "CS" and name "Computer Science"
    And the annual budget is $5,000,000
    And the department head is "Dr. Jane Smith"
    And the status is "Active"
    When I send a POST request to "/api/departments"
    Then the response status should be 201 Created
    And the response should contain the department with code "CS"
    And the department should have an ID assigned
    And the created_at timestamp should be set
    And the updated_at timestamp should be set

  Scenario: Reject department creation with invalid code format
    Given I have department data with code "cs" (lowercase)
    When I send a POST request to "/api/departments"
    Then the response status should be 400 Bad Request
    And the error message should indicate "Department code must be 2-6 uppercase letters"

  Scenario: Reject department creation with code too short
    Given I have department data with code "C" (1 letter)
    When I send a POST request to "/api/departments"
    Then the response status should be 400 Bad Request
    And the error message should indicate "Department code must be 2-6 uppercase letters"

  Scenario: Reject department creation with code too long
    Given I have department data with code "COMPSCI" (7 letters)
    When I send a POST request to "/api/departments"
    Then the response status should be 400 Bad Request
    And the error message should indicate "Department code must be 2-6 uppercase letters"

  Scenario: Reject department creation with duplicate code
    Given a department with code "CS" already exists
    When I send a POST request to "/api/departments" with code "CS"
    Then the response status should be 409 Conflict
    And the error message should indicate "Department code already exists"

  Scenario: Reject department creation with negative budget
    Given I have department data with valid code "MATH"
    And the annual budget is -100,000
    When I send a POST request to "/api/departments"
    Then the response status should be 400 Bad Request
    And the error message should indicate "Budget must be non-negative"

  Scenario: Reject department creation with budget exceeding maximum
    Given I have department data with valid code "MATH"
    And the annual budget is $15,000,000
    When I send a POST request to "/api/departments"
    Then the response status should be 400 Bad Request
    And the error message should indicate "Budget cannot exceed $10,000,000"

  Scenario: Reject department creation with invalid status
    Given I have department data with valid code "PHYS"
    And the status is "Pending"
    When I send a POST request to "/api/departments"
    Then the response status should be 400 Bad Request
    And the error message should indicate "Status must be Active or Inactive"

  Scenario: Create department with default status
    Given I have department data with code "BIOL" and name "Biology"
    And no status is specified
    When I send a POST request to "/api/departments"
    Then the response status should be 201 Created
    And the department status should be "Active" by default

  # CRUD Operations - Read

  Scenario: Get department by valid code
    Given a department with code "CS" exists in the system
    When I send a GET request to "/api/departments/CS"
    Then the response status should be 200 OK
    And the response should contain department with code "CS"
    And the response should include all department fields

  Scenario: Get department by ID
    Given a department with ID 1 exists in the system
    When I send a GET request to "/api/departments/1"
    Then the response status should be 200 OK
    And the response should contain department with ID 1

  Scenario: Get non-existent department
    When I send a GET request to "/api/departments/NONEXIST"
    Then the response status should be 404 Not Found
    And the error message should indicate "Department not found"

  Scenario: List all departments
    Given the following departments exist:
      | code  | name                | budget      | head            | status   |
      | CS    | Computer Science    | 5000000.00  | Dr. Jane Smith  | Active   |
      | MATH  | Mathematics         | 3000000.00  | Dr. John Doe    | Active   |
      | PHYS  | Physics             | 4000000.00  | Dr. Alice Brown | Inactive |
    When I send a GET request to "/api/departments"
    Then the response status should be 200 OK
    And the response should contain 3 departments
    And all departments should have complete data

  Scenario: Filter departments by status
    Given multiple departments exist with different statuses
    When I send a GET request to "/api/departments?status=Active"
    Then the response status should be 200 OK
    And all returned departments should have status "Active"

  # CRUD Operations - Update

  Scenario: Update department with valid data
    Given a department with code "CS" exists
    When I send a PUT request to "/api/departments/CS" with:
      | field            | value                          |
      | department_name  | Computer Science & Engineering |
      | annual_budget    | 6000000.00                     |
      | department_head  | Dr. Michael Chen               |
    Then the response status should be 200 OK
    And the department should be updated with new values
    And the updated_at timestamp should be newer than created_at

  Scenario: Update department status to Inactive
    Given a department with code "HIST" and status "Active" exists
    When I send a PUT request to "/api/departments/HIST" with status "Inactive"
    Then the response status should be 200 OK
    And the department status should be "Inactive"

  Scenario: Reject update with invalid budget
    Given a department with code "CS" exists
    When I send a PUT request to "/api/departments/CS" with budget $12,000,000
    Then the response status should be 400 Bad Request
    And the error message should indicate "Budget cannot exceed $10,000,000"

  Scenario: Update non-existent department
    When I send a PUT request to "/api/departments/NOTFOUND" with valid data
    Then the response status should be 404 Not Found

  # CRUD Operations - Delete

  Scenario: Delete existing department
    Given a department with code "OLD" exists
    When I send a DELETE request to "/api/departments/OLD"
    Then the response status should be 204 No Content
    And the department should no longer exist in the system

  Scenario: Delete non-existent department
    When I send a DELETE request to "/api/departments/NOTFOUND"
    Then the response status should be 404 Not Found

  # Redis Caching

  Scenario: First department lookup caches the result
    Given a department with code "CS" exists
    And the Redis cache is empty for "CS"
    When I send a GET request to "/api/departments/CS"
    Then the response status should be 200 OK
    And the department should be stored in Redis cache
    And the cache TTL should be 1 hour

  Scenario: Second department lookup retrieves from cache
    Given a department with code "CS" is cached in Redis
    When I send a GET request to "/api/departments/CS"
    Then the response status should be 200 OK
    And the response should include "X-Cache: HIT" header
    And the database should not be queried

  Scenario: Cache miss retrieves from database
    Given a department with code "MATH" exists in database
    And "MATH" is not in Redis cache
    When I send a GET request to "/api/departments/MATH"
    Then the response status should be 200 OK
    And the response should include "X-Cache: MISS" header
    And the database should be queried
    And the result should be cached for future requests

  Scenario: Update operation invalidates cache
    Given a department with code "CS" is cached in Redis
    When I send a PUT request to "/api/departments/CS" with updated data
    Then the response status should be 200 OK
    And the Redis cache for "CS" should be invalidated
    And the next GET request should show "X-Cache: MISS"

  Scenario: Delete operation invalidates cache
    Given a department with code "OLD" is cached in Redis
    When I send a DELETE request to "/api/departments/OLD"
    Then the response status should be 204 No Content
    And the Redis cache for "OLD" should be invalidated

  Scenario: Cache expires after TTL
    Given a department with code "CS" was cached 61 minutes ago
    When I send a GET request to "/api/departments/CS"
    Then the response should include "X-Cache: MISS" header
    And the department should be re-cached with fresh TTL

  # Budget Tracking

  Scenario: Track budget allocation for new department
    Given I create a department with code "ENG" and budget $7,500,000
    When I send a GET request to "/api/departments/ENG/budget-history"
    Then the response status should be 200 OK
    And the budget history should contain one entry
    And the entry should show initial budget of $7,500,000

  Scenario: Track budget changes over time
    Given a department with code "CS" and budget $5,000,000 exists
    When I update the budget to $5,500,000
    And I update the budget again to $6,000,000
    Then the budget history should contain 3 entries
    And the entries should show the progression: $5,000,000 → $5,500,000 → $6,000,000

  Scenario: Calculate budget utilization percentage
    Given a department with code "CS" has annual budget $5,000,000
    And the department has spent $3,750,000
    When I send a GET request to "/api/departments/CS/budget-utilization"
    Then the response status should be 200 OK
    And the utilization percentage should be 75%
    And the remaining budget should be $1,250,000

  Scenario: Generate budget report for all departments
    Given multiple departments exist with different budgets and utilization
    When I send a GET request to "/api/departments/budget-report"
    Then the response status should be 200 OK
    And the report should include total allocated budget
    And the report should include total utilization across departments
    And the report should list each department's budget status

  # Data Validation

  Scenario: Validate department code pattern
    Given I attempt to create departments with these codes:
      | code     | valid |
      | CS       | true  |
      | MATH     | true  |
      | ENG      | true  |
      | COMPUT   | true  |
      | cs       | false |
      | C        | false |
      | TOOLONG  | false |
      | CS-123   | false |
      | COM!     | false |
    Then only valid codes should be accepted
    And invalid codes should return 400 Bad Request

  Scenario: Validate budget constraints
    Given I attempt to create departments with these budgets:
      | budget        | valid |
      | 0             | true  |
      | 1000000.00    | true  |
      | 10000000.00   | true  |
      | -100.00       | false |
      | 10000000.01   | false |
      | 15000000.00   | false |
    Then only valid budgets should be accepted

  Scenario: Validate status values
    Given I attempt to create departments with these statuses:
      | status    | valid |
      | Active    | true  |
      | Inactive  | true  |
      | Pending   | false |
      | Archived  | false |
      | ""        | false |
    Then only "Active" and "Inactive" should be accepted

  # Error Handling

  Scenario: Handle database connection failure gracefully
    Given the database connection is lost
    When I send a GET request to "/api/departments/CS"
    Then the response status should be 503 Service Unavailable
    And the error message should indicate "Database unavailable"

  Scenario: Handle Redis connection failure gracefully
    Given Redis is unavailable
    And a department with code "CS" exists in database
    When I send a GET request to "/api/departments/CS"
    Then the response status should be 200 OK
    And the response should include "X-Cache: BYPASS" header
    And the department should be retrieved from database

  Scenario: Handle malformed JSON in request body
    When I send a POST request to "/api/departments" with invalid JSON
    Then the response status should be 400 Bad Request
    And the error message should indicate "Invalid JSON format"

  # Performance & Cache Warming

  Scenario: Warm cache for frequently accessed departments
    Given departments "CS", "MATH", and "ENG" are accessed most frequently
    When the cache warming process runs
    Then these departments should be pre-loaded into Redis cache
    And subsequent requests should show "X-Cache: HIT"

  Scenario: Handle concurrent department creation
    Given 10 concurrent requests to create department "CS"
    Then only 1 department should be created successfully
    And 9 requests should fail with 409 Conflict

  # API Documentation Scenarios

  Scenario: API returns proper Content-Type headers
    When I send any request to "/api/departments"
    Then the response should have "Content-Type: application/json" header

  Scenario: API returns proper CORS headers
    When I send an OPTIONS request to "/api/departments"
    Then the response should include appropriate CORS headers

  Scenario: API endpoints follow RESTful conventions
    Then the following endpoints should be available:
      | method | endpoint                              | description              |
      | POST   | /api/departments                      | Create department        |
      | GET    | /api/departments                      | List all departments     |
      | GET    | /api/departments/{code}               | Get by code              |
      | GET    | /api/departments/{id}                 | Get by ID                |
      | PUT    | /api/departments/{code}               | Update department        |
      | DELETE | /api/departments/{code}               | Delete department        |
      | GET    | /api/departments/{code}/budget-history| Get budget history       |
      | GET    | /api/departments/{code}/budget-utilization | Get budget utilization |
      | GET    | /api/departments/budget-report        | Get overall budget report|
