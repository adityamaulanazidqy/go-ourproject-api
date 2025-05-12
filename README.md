# Go OurProject API

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-2.50.0-00ADD8)](https://gofiber.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/adityamaulanazidqy/go-ourproject-api/pulls)

🚧 **Project Status: Active Development** 🚧  
*Features and structure may change during development*

Go OurProject API is a high-performance RESTful API built with Fiber framework for Golang. A modern backend solution with modular architecture supporting secure authentication, Redis caching, and email integration.

## Table of Contents
- [Features](#-features)
- [Technologies](#-technologies)
- [Project Structure](#-project-structure)
- [Installation](#-installation)
- [API Endpoints](#-api-endpoints)

## ✨ Features

### Core Architecture
- Clean MVC pattern implementation
- Modular design with separation of concerns
- Centralized error handling
- Environment-based configuration

### Security
- JWT authentication
- Email OTP verification
- Protected routes with middleware
- Secure password hashing

### Performance
- Redis caching layer
- Optimized database queries
- Efficient routing with Fiber
- Structured logging with Logrus

## 🛠 Technologies

| Component       | Technology                  |
|----------------|----------------------------|
| Language       | Go 1.21+                   |
| Framework      | Fiber v2                   |
| Database       | PostgreSQL                 |
| Cache          | Redis                      |
| Logging        | Logrus                     |
| Email          | (Configure your provider)  |
| Testing        | Go test                    |

## 📂 Project Structure

```bash
.
├── config/               # Application configurations
│   ├── database.go       # DB connection setup
│   ├── redis.go          # Redis client
│   └── logger.go         # Logging configuration
├── controllers/          # Business logic
│   ├── auth.go           # Authentication handlers
│   └── user.go           # User management
├── helpers/              # Utility functions
│   ├── otp.go            # OTP generation
│   └── response.go       # API response formatting
├── middlewares/          # Fiber middlewares
│   ├── auth.go           # Authentication
│   ├── logger.go         # Request logging
│   └── recovery.go       # Error recovery
├── models/               # Data structures
│   └── user.go           # User model
├── repositories/         # Database operations
│   └── user_repo.go      # User repository
├── routes/               # API endpoints
│   ├── api.go            # Route definitions
│   └── middleware.go     # Route middlewares
├── test/                 # Test cases
├── go.mod                # Dependency management
├── main.go               # Application entry point
└── .env.example          # Environment template
```

## 🚀 Installation

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Redis 6+
- SMTP credentials (for OTP emails)

### Quick Start
1. Clone the repository:
```bash
git clone https://github.com/adityamaulanazidqy/go-ourproject-api.git
cd go-ourproject-api
```
2. Set up environment variables:
```bash
cp .env.example .env
```
Edit the ``.env`` file with your configuration:
```bash
# Application
APP_PORT=8673
APP_ENV=development
JWT_SECRET=your_secure_jwt_secret

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_SSL_MODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Email
EMAIL_PROVIDER=sendgrid  # or your provider
EMAIL_API_KEY=your_api_key
EMAIL_FROM=noreply@yourdomain.com
```
3. Install Dependencies
```bash
go mod tidy
```
4. go run main.go
```bash
go run main.go
```
The server will be available at ``http://localhost:8673``

## 📌 API Endpoints

### Authentication
| Method    | Endpoint               | Description |
|-----------|------------------------|-------------|
| POST      | /auth/register         | registration |
| POST      | /auth/register         | User login  |
| POST      | /logout                | User logout |
| POST      | /otp/send-otp          | Send OTP to email |
| POST      | /otp/verify-otp        | Verify OTP code |

### User Management
| Method    | Endpoint               | Description |
|-----------|------------------------|-------------|
| POST      | /update_password       | User update pass |  

_Note: Protected endpoints require JWT in Authorization header_  

## 💻 Happy Coding!
For support, please open an issue or contact maintainers.
