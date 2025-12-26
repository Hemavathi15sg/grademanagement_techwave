# Development Standards for Grade Management API

## Model Conventions
- All models MUST include: ID (int), CreatedAt (time.Time), UpdatedAt (time.Time)
- Use snake_case for JSON tags
- Include omitempty for optional fields

## Validation Standards
- Status fields must use predefined constants
- All create/update handlers must validate input
- Return 400 Bad Request with clear error messages

## API Documentation
- Every endpoint must have:
  - Description
  - Request body schema
  - Response codes (200, 201, 400, 404, 500)
  - Example requests/responses

## Test Standards
- Use gomock for interface mocking
- Test data factories must support:
  - Default builders
  - Chaining for custom fields
  - Common scenario helpers (valid, invalid, edge cases)