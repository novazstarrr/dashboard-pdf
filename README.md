

## Running the application

1. Run `docker compose up` to start the application.
2. The frontend will be available at http://localhost:3000.
3. The backend will be available at http://localhost:8080.

# File Management System

A full-stack application for managing files and users with secure authentication, file uploads, and user management capabilities.

## Features

- 🔐 User Authentication (Register/Login)
- 📁 File Management
  - Upload PDF files
  - View files
  - Download files
  - Share files via links
  - Search functionality
  - Pagination
- 👥 User Management
  - CRUD operations for users
  - Search users
  - Pagination
- 🔒 Security
  - JWT Authentication
  - Rate limiting
  - CORS protection
  - Secure headers

## Tech Stack

### Backend
- Go
- Gorilla Mux (Router)
- GORM (ORM)
- SQLite (Database)
- Zap (Logging)
- Swagger (API Documentation)

### Frontend
- React
- React Bootstrap
- Axios
- React Icons
- React Toastify

## Security Features

- JWT-based authentication
- Rate limiting for API endpoints
- CORS protection with whitelisted origins
- Secure HTTP headers
- Password hashing
- File type validation
- Maximum upload size limits
