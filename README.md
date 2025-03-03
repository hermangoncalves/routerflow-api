# Router Management System Backend

A simple Golang backend for managing MikroTik routers using the RouterOS API. This project provides RESTful APIs for user authentication and router management, following a feature-based structure.

## Features
- **User Authentication:** Secure login and registration with bcrypt password hashing.
- **Router Management:** Register routers and fetch their status via the RouterOS API.

## Prerequisites
- **Go**: Version 1.21 or higher.
- **MikroTik Router (Optional):** To test `/routers/:id/status`, a router should be accessible (e.g., at `192.168.1.1` with default credentials `admin:password`).

## Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/hermangoncalves/routerflow-api.git
   cd router-management
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run cmd/main.go
   ```
   - This starts the server on `http://localhost:8080`.
   - A SQLite database (`data.db`) and log file (`app.log`) are created in the project root.

## API Endpoints

### Authentication

#### **Login**
Authenticate a user and receive a token (currently a fake JWT).
```bash
curl -X POST -d '{"username":"admin","password":"password"}' http://localhost:8080/login
```
**Response:**
```json
{"token":"fake-jwt-token"}
```

#### **Register User**
Create a new user in the system.
```bash
curl -X POST -d '{"username":"newuser","password":"newpass"}' http://localhost:8080/register
```
**Response:**
```json
{"message":"user registered"}
```

### Router Management

#### **Register Router**
Add a new router to the database.
```bash
curl -X POST -d '{"ip":"192.168.1.2","username":"admin","password":"pass123"}' http://localhost:8080/routers/register
```
**Response:**
```json
{"message":"router registered"}
```

#### **Get Router Status**
Fetch router status (e.g., uptime) from a MikroTik router via the RouterOS API.
```bash
curl http://localhost:8080/routers/1/status
```
**Example Response:**
```json
[{"uptime":"10d2h","version":"6.47.10","cpu":"MIPS 74Kc","cpu-load":"5"}]
```

## Notes
- **Initial Database Setup:**
  - Default user: `admin:password` (bcrypt hashed password).
  - Default router: ID `1`, IP `192.168.1.1`, credentials `admin:password`.
- **Logging:** Logs are written to `app.log` for debugging.
- **Security:** Authentication for router endpoints is not yet implemented (planned for future updates).

## Testing

1. **Start the Server:**
   ```bash
   go run cmd/main.go
   ```

2. **Verify API Functionality:**
   - **Login:**
     ```bash
     curl -X POST -d '{"username":"admin","password":"password"}' http://localhost:8080/login
     ```
     **Expected Output:** `{"token":"fake-jwt-token"}`

   - **Register User:**
     ```bash
     curl -X POST -d '{"username":"newuser","password":"newpass"}' http://localhost:8080/register
     ```
     **Expected Output:** `{"message":"user registered"}`

   - **Register Router:**
     ```bash
     curl -X POST -d '{"ip":"192.168.1.2","username":"admin","password":"pass123"}' http://localhost:8080/routers/register
     ```
     **Expected Output:** `{"message":"router registered"}`

   - **Router Status:**
     ```bash
     curl http://localhost:8080/routers/1/status
     ```
     **Expected Output:**
     - If router is available: `[{"uptime":"..."}]`
     - If unavailable: `{"error":"connection failed"}`
