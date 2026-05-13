# Компьютерийн хараанд суурилсан эд хөрөнгийн бүртгэлийн систем

Computer Vision-Based Intelligent Asset Management System

**Бакалаврын төгсөлтийн ажил** — Б.Оргил, ШУТИС МХТС 2026

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend API | Go 1.23 + Gin Framework |
| CV Service | Python 3.11 + FastAPI + YOLO26/Ultralytics + OpenCV |
| Database | PostgreSQL 15 |
| Cache | Redis 7 |
| Frontend | React 18 + TypeScript + Ant Design |
| DevOps | Docker + Docker Compose |

## Project Structure

```
├── backend/                 # Go API server
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Environment configuration
│   │   ├── database/        # DB connection, migration, seeding
│   │   ├── dto/             # Request/response types
│   │   ├── handler/         # HTTP controllers
│   │   ├── middleware/       # JWT auth + RBAC
│   │   ├── model/           # GORM entity models (12 tables)
│   │   ├── repository/      # Data access layer
│   │   ├── router/          # Route definitions
│   │   └── service/         # Business logic
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── cv-service/              # Python CV microservice
│   ├── main.py              # FastAPI + YOLO detection
│   ├── models/              # Place .pt weights here
│   ├── uploads/             # Processed images
│   ├── requirements.txt
│   └── Dockerfile
├── frontend/                # React app (TODO)
├── docker-compose.yml       # 5-container orchestration
└── .env.example             # Environment template
```

## Quick Start

### 1. Clone & configure
```bash
git clone <repo-url>
cd Computer-Vision-Based-Asset_Management-System
cp .env.example .env
# Edit .env as needed
```

### 2. Start all services
```bash
docker compose up --build -d
```

This starts:
- **PostgreSQL** on port 5432
- **Redis** on port 6379
- **Go Backend** on port 8080
- **CV Service** on port 8000

### 3. Default admin login
```
Username: admin
Password: Admin@123
```

## API Endpoints

### Authentication
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/login` | Login, returns JWT |
| GET | `/api/v1/auth/me` | Current user profile |
| PUT | `/api/v1/auth/password` | Change password |

### Assets (Custodian, Admin)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/assets` | Create asset + auto-barcode |
| GET | `/api/v1/assets` | List with pagination & search |
| GET | `/api/v1/assets/:id` | Get asset detail |
| PUT | `/api/v1/assets/:id` | Update asset |
| DELETE | `/api/v1/assets/:id` | Dispose asset |
| POST | `/api/v1/assets/:id/assign` | Assign to employee |
| POST | `/api/v1/assets/:id/transfer` | Transfer between users |
| POST | `/api/v1/assets/:id/dispose` | Dispose with reason |
| GET | `/api/v1/assets/:id/history` | Asset change history |

### CV Audit (Custodian, Admin)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/audits` | Start audit session |
| GET | `/api/v1/audits` | List audit sessions |
| GET | `/api/v1/audits/:id` | Get audit details + findings |
| POST | `/api/v1/audits/:id/cv` | Upload image → run YOLO detection |

### Depreciation (Accountant, Admin)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/depreciation/calculate` | Calculate monthly depreciation |
| POST | `/api/v1/depreciation/revalue` | Revalue asset |

### Reports (Custodian, Accountant, Admin)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/dashboard` | Dashboard statistics |
| POST | `/api/v1/reports/generate` | Generate CSV/JSON report |
| GET | `/api/v1/reports/download/:file` | Download report file |

### Users (Admin only)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/users` | Create user |
| GET | `/api/v1/users` | List users |
| PUT | `/api/v1/users/:id` | Update user |
| PUT | `/api/v1/users/:id/deactivate` | Deactivate user |

### Locations, Departments, Asset Models, Notifications
Full CRUD endpoints available — see `router/router.go` for complete list.

## User Roles (RBAC)

| Role | Permissions |
|------|------------|
| **ADMIN** | Full access: users, departments, locations, assets, audits, reports |
| **ASSET_CUSTODIAN** (Нярав) | Assets, audits, locations, reports |
| **ACCOUNTANT** (Нягтлан бодогч) | Depreciation, revaluation, reports, asset view |
| **EMPLOYEE** (Ажилтан) | View own assets, notifications |

## CV Service API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/detect` | Detect objects in image |
| POST | `/detect/batch` | Batch detection |
| GET | `/model-info` | Current model information |
| GET | `/health` | Service health check |

## Development (without Docker)

### Backend
```bash
cd backend
# Set env vars or create .env
export DB_HOST=localhost DB_PORT=5432 DB_USER=asset_admin \
       DB_PASSWORD=changeme_secret DB_NAME=asset_management
go run cmd/server/main.go
```

### CV Service
```bash
cd cv-service
pip install -r requirements.txt
python main.py
```
