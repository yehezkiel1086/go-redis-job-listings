# Go REST API Caching and Secure JWT Token Rotation

A job listings REST API built with Go, demonstrating a Redis caching strategy on top of a hexagonal architecture. Features user authentication with JWT, job posting, and job enrollment — with Redis handling both response caching and refresh token storage.

## Tech Stack

- **Go** — core language
- **Gin** — HTTP web framework
- **GORM** — ORM for PostgreSQL
- **PostgreSQL** — primary database
- **Redis** — caching and token store
- **Docker** — containerized infrastructure
- **Task** — task runner
- **Air** — live reload for development

## Features

- **Authentication** — JWT-based login with access and refresh tokens. Refresh tokens are stored in Redis and support true server-side revocation on logout and token rotation on refresh.
- **Users** — registration, profile management, and admin controls with RBAC middleware.
- **Jobs** — full CRUD for job listings with filtering by type, experience level, location, and keyword search. Filter-aware cache keys ensure each unique query is cached independently.
- **Enrollments** — users can apply to job listings. Job owners and admins can accept or reject applications.
- **Redis Caching** — cache-aside pattern on all read operations with targeted invalidation on writes across users, jobs, and enrollments.

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- [Task](https://taskfile.dev)
- [Air](https://github.com/air-verse/air)

### Installation

**1. Clone the repository**
```bash
git clone https://github.com/yehezkiel1086/go-redis-job-listings.git
cd go-redis-job-listings
```

**2. Install dependencies**
```bash
go mod tidy
```

**3. Set up environment variables**
```bash
cp .env.example .env
```

Update `.env` with your configuration:
```dotenv
APP_NAME=go-redis-job-listings
APP_ENV=development

HTTP_HOST=127.0.0.1
HTTP_PORT=8080

DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=job_listings

REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=admin
REDIS_DB=0

ACCESS_TOKEN_SECRET=your_access_token_secret
REFRESH_TOKEN_SECRET=your_refresh_token_secret

# access token in mins, refresh in days
ACCESS_TOKEN_DURATION=15
REFRESH_TOKEN_DURATION=7
```

**4. Start PostgreSQL and Redis**
```bash
task compose:up
```

**5. Run the development server**
```bash
task dev
```

The API will be available at `http://127.0.0.1:8080`.

## Available Tasks

| Command | Description |
|---|---|
| `task compose:up` | Start all Docker containers |
| `task compose:down` | Stop all Docker containers |
| `task dev` | Start the development server with live reload |
| `task db:cli` | Open the PostgreSQL CLI |
| `task redis:cli` | Open the Redis CLI |

## API Endpoints

### Auth
| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/v1/register` | Public | Register a new user |
| POST | `/api/v1/login` | Public | Login and receive tokens |
| POST | `/api/v1/refresh` | Public | Refresh access token |
| POST | `/api/v1/logout` | Authenticated | Logout and revoke refresh token |

### Users
| Method | Endpoint | Access | Description |
|---|---|---|---|
| GET | `/api/v1/users` | Admin | List all users |
| GET | `/api/v1/users/:id` | Owner or Admin | Get a user by ID |
| PUT | `/api/v1/users/:id` | Owner or Admin | Update a user |
| DELETE | `/api/v1/users/:id` | Owner or Admin | Delete a user |

### Jobs
| Method | Endpoint | Access | Description |
|---|---|---|---|
| GET | `/api/v1/jobs` | Authenticated | List all jobs (supports filters) |
| GET | `/api/v1/jobs/:id` | Authenticated | Get a job by ID |
| GET | `/api/v1/jobs/me` | Authenticated | Get jobs posted by the logged-in user |
| GET | `/api/v1/jobs/user/:id` | Owner or Admin | Get jobs posted by a specific user |
| POST | `/api/v1/jobs` | Admin | Create a job listing |
| PUT | `/api/v1/jobs/:id` | Admin | Update a job listing |
| DELETE | `/api/v1/jobs/:id` | Admin | Delete a job listing |

### Enrollments
| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/v1/jobs/:id/enroll` | Authenticated | Apply to a job |
| GET | `/api/v1/enrollments/me` | Authenticated | Get the logged-in user's applications |
| GET | `/api/v1/jobs/:id/enrollments` | Owner or Admin | Get all applicants for a job |
| PUT | `/api/v1/enrollments/:id` | Owner or Admin | Accept or reject an application |
| DELETE | `/api/v1/enrollments/:id` | Owner or Admin | Withdraw an application |

## Caching Strategy

All read operations follow the **cache-aside pattern**:
```
Request → Check Redis → HIT: return cached data
                      → MISS: query PostgreSQL → store in Redis → return data
```

Cache keys follow a consistent naming convention:

| Key Pattern | Description | TTL |
|---|---|---|
| `user:{id}` | Single user | 5 min |
| `user:all` | All users list | 5 min |
| `job:{id}` | Single job | 5 min |
| `job:all:{params}` | Filtered job list | 5 min |
| `job:user:{id}` | Jobs by a user | 5 min |
| `enrollment:{id}` | Single enrollment | 5 min |
| `enrollment:user:{id}` | Enrollments by a user | 5 min |
| `enrollment:job:{id}` | Enrollments for a job | 5 min |
| `refresh_token:{token}` | Active session | Matches token expiry |
