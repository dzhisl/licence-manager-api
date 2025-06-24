# License Manager API

A robust RESTful API for managing user licenses, built with Go, Gin, and MongoDB. This service provides endpoints for license verification, user management, device binding, and more, with a focus on security and extensibility.

## Features

- **User Management**: Create, retrieve, and delete users.
- **License Management**: Issue, verify, renew, and update licenses.
- **Device Management**: Add, remove, and reset devices (HWIDs) per user.
- **Third-party Bindings**: Bind Discord and Telegram accounts to users.
- **Status Control**: Change license status (active, frozen, burned).
- **Swagger Documentation**: Interactive API docs available.
- **Admin Authentication**: Secure private endpoints with API key middleware.
- **MongoDB Storage**: Persistent, scalable data backend.

## Tech Stack

- **Language**: Go (1.24+)
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: MongoDB
- **Config**: Viper
- **Logging**: Uber Zap
- **API Docs**: Swagger (OpenAPI)

## Getting Started

### Prerequisites

- Go 1.24 or newer
- MongoDB instance (local or remote)
- (Optional) [swag](https://github.com/swaggo/swag) for generating Swagger docs

### Configuration

Create a `.env` file in the project root with the following variables:

```
STAGE_ENV=production (or dev)
ADMIN_SECRET_KEY=your_password_to_private_endpoints
MONGODB_URI=mongodb://localhost:27017
LICENSE_PREFIX=your_prefix
LICENSE_LENGTH=16
```

### Installation

```sh
git clone https://github.com/dzhisl/license-manager-api.git
cd licence-manager-api
go mod tidy
```

### Running the API

```sh
make run
```

The server will start on `http://localhost:8080`.

### API Documentation

Swagger UI is available at:  
`http://localhost:8080/swagger/index.html`

## Usage

### Public Endpoints

- `POST /api/license/verify` — Verify a license by key and HWID
- `GET /api/ping` — Health check

### Private (Admin) Endpoints

Require `X-API-Key` header for authentication.

- `POST /api/user/create` — Create a new user
- `GET /api/user` — Retrieve user by Telegram ID, Discord ID, or license key
- `POST /api/user/:user_id/device` — Add a device (HWID)
- `DELETE /api/user/:user_id/device` — Remove a device (HWID)
- `POST /api/user/:user_id/devices/reset` — Reset all devices
- `POST /api/user/:user_id/license/status` — Change license status
- `POST /api/user/:user_id/license/hwid_limit` — Update HWID limit
- `POST /api/user/:user_id/license/renew` — Renew license
- `POST /api/user/:user_id/discord` — Bind Discord account
- `POST /api/user/:user_id/telegram` — Bind Telegram account
- `DELETE /api/user/:user_id` — Delete user

See [Swagger docs](http://localhost:8080/swagger/index.html) for full request/response schemas.

## Data Model

### User

```go
type User struct {
    Id         int
    TelegramId int
    DiscordId  int
    License    License
    CreatedAt  int64
}
```

### License

```go
type License struct {
    Key            string
    MaxActivations int
    Devices        []string
    IssuedAt       int64
    ExpiresAt      int64
    Status         string // "active", "frozen", "burned"
}
```

## Developments

### Run Tests

```sh
make test
```

### Run API Router Tests

```sh
make test-api
```

### Run DB Tests

```sh
make test-storage
```

### Generate Swagger Docs

```sh
make swagger
```

## Project Structure

```
cmd/server/           # Main entry point
internal/api/         # API handlers, middleware, router
internal/storage/     # MongoDB storage logic and models
pkg/config/           # Configuration loader
pkg/logger/           # Logging setup
docs/                 # Swagger/OpenAPI docs
Makefile              # Common dev commands
```

## Contributing

Contributions are welcome! Please open issues or submit pull requests.

## Contacts

telegram: [@isdzh](http://t.me/isdzh)
