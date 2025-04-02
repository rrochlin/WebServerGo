# WebServerGo

A guided project from [boot.dev](www.boot.dev) teaching the basics of http web servers.

## Overview

The project covered setting up a web server in go using the builtin http package
with a postgres backend database ran locally on the machine. The application is host
to a few one static HTML file but is largely designed as a REST API for managing
posts and user accounts. Chirpy acts as a social media platform that allows user
creation and access via hashed passwords. Users who are authenticated are provided
a 1 hour JWT and 60 day refresh token for adding and removing their posts.

### Concepts Covered

- Authentication
  - Via hashed password using bycrypt
  - Via short-lived JWTs generated and provided to users on a successful login
  - Longer Lived refresh tokens stored in the postgres DB that are revokable

- Routing
  - Basics of creating/using GET, POST, DELETE, PUT route handlers
  - Using Go's fileserver to serve static assests
  - routes with named parameters and/or optional queries

- Go Package management basics
  - installing and managing packages/dependencies for a golang project
  - referencing different packages within a go project

- Database management via sqlc
  - setting up database migrations with sqlc
  - creating database queries to server route handlers and perform data operations

- Webhooks
  - setting up a webhook reciever and authenticating via API key
  - linked user membership status to requests recieved via webhooks to simulate
    recieving subscription information from stripe via webhooks.

### Quick Start

1. Make sure to install [GO](https://go.dev/doc/install) as well as [Postgres](https://www.postgresql.org/download/) and set up a user for postgres
2. clone the repository and create a `.env` file in the root directory

    ```env
    DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable",
    PORT=8080
    HOST_URL='0.0.0.0'
    PLATFORM='dev'
    SECRET="tG5eBGAUiVSMR2pV5WD5DqIHuuq0ForXbFbWRxy/P+gTo3wRfg+pId9NV0cQDniiqlnpss02+v7TXCupZ+TI/w=="
    POLKA_API_KEY="f271c81ff7084ee5b99a5091b42d486e"
    ```

3. run `go build -o out && ./out` to run the project

## API Endpoints

### Health & Admin

#### GET /api/healthz
Check if the API is running.

**Response**: `200 OK` with text body "OK"

#### GET /admin/metrics
View server metrics.

**Response**: `200 OK` with HTML content showing visit count

#### POST /admin/reset
Reset the server state and delete all users (Development environment only).

**Response**: `200 OK` with text "Count Reset"

**Restrictions**: Only available in development environment

### User Management

#### POST /api/users
Create a new user account.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response**: `201 Created` with the created user object (excluding password)
```json
{
  "id": "uuid",
  "email": "user@example.com"
}
```

#### PUT /api/users
Update user information.

**Authentication**: Required

**Request Body**:
```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

**Response**: `200 OK` with updated user object (excluding password)
```json
{
  "id": "uuid",
  "email": "newemail@example.com"
}
```

### Chirps

#### POST /api/chirps
Create a new chirp.

**Authentication**: Required

**Request Body**:
```json
{
  "body": "This is my chirp message!"
}
```

**Response**: `201 Created` with the created chirp object
```json
{
  "id": "uuid",
  "body": "This is my chirp message!",
  "userId": "user-uuid",
  "createdAt": "2023-04-01T12:00:00Z"
}
```

**Notes**:
- Chirps are limited to 140 characters
- Certain words are automatically censored with "****"

#### GET /api/chirps
Get all chirps.

**Query Parameters**:
- `authorId` (optional): Filter chirps by author ID
- `sort` (optional): Sort order, either "asc" (default) or "desc"

**Response**: `200 OK` with array of chirp objects
```json
[
  {
    "id": "uuid1",
    "body": "This is chirp 1",
    "userId": "user-uuid",
    "createdAt": "2023-04-01T12:00:00Z"
  },
  {
    "id": "uuid2",
    "body": "This is chirp 2",
    "userId": "user-uuid",
    "createdAt": "2023-04-01T12:30:00Z"
  }
]
```

#### GET /api/chirps/:uuid
Get a specific chirp by ID.

**Parameters**:
- `uuid`: The UUID of the chirp

**Response**: `200 OK` with chirp object
```json
{
  "id": "uuid",
  "body": "This is my chirp message!",
  "userId": "user-uuid",
  "createdAt": "2023-04-01T12:00:00Z"
}
```

#### DELETE /api/chirps/:chirpID
Delete a chirp.

**Authentication**: Required

**Parameters**:
- `chirpID`: The ID of the chirp to delete

**Response**: `204 No Content`

**Restrictions**: Users can only delete their own chirps

### Authentication

#### POST /api/login
Log in to user account.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response**: `200 OK` with user info and tokens
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "token": "jwt-token",
  "refreshToken": "refresh-token"
}
```

#### POST /api/refresh
Get a new JWT token using a refresh token.

**Authentication**: Requires a valid refresh token in Authorization header

**Response**: `200 OK` with new token
```json
{
  "token": "new-jwt-token"
}
```

#### POST /api/revoke
Revoke a refresh token.

**Authentication**: Requires the refresh token to be revoked in Authorization header

**Response**: `204 No Content`

### Webhooks

#### POST /api/polka/webhooks
Endpoint for Polka service webhooks.

**Response**: `204 No Content` on success

## Error Handling

The API uses consistent error responses:

- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required or invalid
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error responses include a JSON body with error details:
```json
{
  "error": "Error message"
}
```

## Static Content

The API also serves static content from the `/app` directory at the `/app` path.
