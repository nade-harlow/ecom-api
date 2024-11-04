# Starting up the Application

This is a README file for running the Go application.
## 📋 Prerequisites

- Go 1.x
- PostgreSQL

## 🛠️ Installation

1. Clone the repository


2. Copy the example environment file
```bash
cp .env.example .env
```

3. Install dependencies
```bash
go mod tidy
```

## 🔧 Configuration

Update the `.env` file with your configuration:

```env
PORT=8082
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=yourdatabase
DB_DEFAULT_SCHEMA=public
```

## 🚀 Running the Application

### Local Environment

```bash
make run
```

## 📚 API Documentation

Postman documentation is available at:
```
https://documenter.getpostman.com/view/17760778/2sAY4xBhGE
```

---