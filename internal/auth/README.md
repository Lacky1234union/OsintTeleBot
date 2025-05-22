# Authentication Microservice

This is a standalone authentication microservice that provides JWT-based authentication for your applications.

## Features

- User registration and login
- JWT token generation and validation
- Password hashing with bcrypt
- Role-based access control

## API Endpoints

### Register
```http
POST /auth/register
Content-Type: application/json

{
    "username": "user123",
    "email": "user@example.com",
    "password": "securepassword"
}
```

### Login
```http
POST /auth/login
Content-Type: application/json

{
    "username": "user123",
    "password": "securepassword"
}
```

### Validate Token
```http
GET /auth/validate
Authorization: <jwt_token>
```

## Setup

1. Set environment variables:
```bash
export JWT_SECRET=your-secret-key
export PORT=8080  # optional, defaults to 8080
```

2. Install dependencies:
```bash
go mod download
```

3. Run the service:
```bash
go run cmd/main.go
```

## Integration with Other Services

To use this authentication service with other microservices:

1. Send authentication requests to this service
2. Store the received JWT token
3. Include the token in the Authorization header for subsequent requests
4. Validate tokens using the `/auth/validate` endpoint

## Security Considerations

- Always use HTTPS in production
- Store JWT_SECRET securely
- Implement rate limiting
- Use strong password policies
- Consider implementing token refresh mechanism
- Monitor for suspicious activities

## Development

To add new features or modify existing ones:

1. Update the domain models in `internal/domain/`
2. Implement business logic in `internal/app/`
3. Add new handlers in `internal/interfaces/http/`
4. Update tests accordingly 