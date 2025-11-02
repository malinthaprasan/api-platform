# Development Guide

**Generated:** 2025-11-02

## Prerequisites

### General Requirements
- Git
- Docker and Docker Compose
- Code editor (VS Code recommended)

### Part-Specific Requirements

#### Platform API & Gateway Controller (Go Services)
- **Go:** 1.24+ (Platform API), 1.25.1+ (Gateway Controller)
- **Make:** For build automation (Gateway Controller)
- **SQLite3:** Database (included in Go dependencies)

#### Management Portal (React SPA)
- **Node.js:** v20+
- **pnpm:** Package manager (recommended) or npm
- **TypeScript:** 4.9.5+

#### STS (Next.js Application)
- **Node.js:** v20+
- **pnpm:** Package manager (workspace support)
- **OpenSSL:** For HTTPS certificate generation
- **TypeScript:** 5+

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/wso2/api-platform.git
cd api-platform
```

### 2. Part-Specific Setup

#### Platform API

```bash
cd platform-api/src

# Install Go dependencies
go mod download

# Build the application
go build -o main cmd/main.go

# Run the application
./main
```

**Configuration:**
- Configuration files in `config/`
- Database initialization handled automatically (SQLite)

#### Gateway Controller

```bash
cd gateway/gateway-controller

# Install dependencies
go mod download

# Build using Makefile
make build

# Or build manually
go build -o gateway-controller cmd/main.go

# Run the application
./gateway-controller
```

**Configuration:**
- Configuration files in `config/`
- Uses Koanf for configuration management
- Supports environment variables and file-based config

#### Management Portal

```bash
cd portals/management-portal

# Install dependencies
pnpm install
# or: npm install

# Run development server
pnpm dev
# or: npm run dev

# Build for production
pnpm build
# or: npm run build

# Preview production build
pnpm preview
# or: npm run preview
```

**Development Server:** http://localhost:5173 (Vite default)

**Configuration:**
- `vite.config.ts` - Build configuration
- `tsconfig.json` - TypeScript configuration
- `.env*` files - Environment variables (if present)

#### STS

```bash
cd sts

# Install dependencies for all workspace packages
pnpm install

# Run gate-app in development mode
cd gate-app
pnpm dev

# Build gate-app
pnpm build

# Run in production mode
pnpm start
```

**Setup Script:**
```bash
cd sts
chmod +x kickstart.sh
./kickstart.sh
```

**HTTPS Certificates:**
```bash
cd sts/gate-app
pnpm ensure-certificates
# Generates server.key and server.cert for localhost
```

## Development Workflow

### Running Individual Parts

#### Terminal 1: Platform API
```bash
cd platform-api/src
go run cmd/main.go
```

#### Terminal 2: Gateway Controller
```bash
cd gateway/gateway-controller
go run cmd/main.go
```

#### Terminal 3: Management Portal
```bash
cd portals/management-portal
pnpm dev
```

#### Terminal 4: STS
```bash
cd sts/gate-app
pnpm dev
```

### Running with Docker

Each part has a Dockerfile:

```bash
# Platform API
cd platform-api
docker build -t platform-api .
docker run -p 8080:8080 platform-api

# Gateway Controller
cd gateway/gateway-controller
docker build -t gateway-controller .
docker run -p 8081:8081 gateway-controller

# STS
cd sts
docker build -t sts .
docker run -p 3000:3000 sts
```

## Testing

### Platform API (Go)
```bash
cd platform-api/src
go test ./...
go test -v ./internal/...
```

### Gateway Controller (Go)
```bash
cd gateway/gateway-controller
make test
# or
go test ./...
```

### Management Portal (React)
```bash
cd portals/management-portal
pnpm test
# or: npm test
```

**Test Files:**
- Jest configuration included
- Tests co-located with components: `*.test.tsx`, `*.spec.tsx`
- Testing Library for React components

## Code Style & Linting

### Go Services (Platform API, Gateway Controller)
```bash
# Format code
go fmt ./...

# Run linter (if golangci-lint installed)
golangci-lint run
```

### Management Portal
```bash
cd portals/management-portal

# Run linter
pnpm lint
# or: npm run lint

# Auto-fix linting issues
pnpm lint --fix
```

### STS
```bash
cd sts/gate-app

# Format code
pnpm format
# or: npm run format

# Run linter
pnpm lint
# or: npm run lint
```

## Environment Variables

### Platform API
- Check `config/` directory for configuration files
- Environment-specific settings can be overridden

### Gateway Controller
- Uses Koanf for configuration
- Supports `.env` files and environment variables
- Configuration hierarchy: env vars → files → defaults

### Management Portal
Create `.env.local` (gitignored):
```env
VITE_API_URL=http://localhost:8080
```

### STS
Create `.env.local` in `sts/gate-app/`:
```env
NODE_ENV=development
PORT=3000
```

## Build Commands

### Development Builds

| Part | Command | Output |
|------|---------|--------|
| Platform API | `go build` | `main` binary |
| Gateway Controller | `make build` | `gateway-controller` binary |
| Management Portal | `pnpm dev` | Dev server (HMR) |
| STS | `pnpm dev` | Next.js dev server |

### Production Builds

| Part | Command | Output |
|------|---------|--------|
| Platform API | `go build -ldflags="-s -w"` | Optimized binary |
| Gateway Controller | `make build-prod` | Optimized binary |
| Management Portal | `pnpm build` | `dist/` folder |
| STS | `pnpm build` | `.next/` folder |

## Common Development Tasks

### Adding a New API Endpoint (Platform API)

1. Define handler in `internal/handler/`
2. Add business logic in `internal/service/`
3. Add data access in `internal/repository/`
4. Register route in server initialization
5. Add tests

### Adding a UI Component (Management Portal)

1. Create component in `src/components/`
2. Use Material-UI components
3. Add TypeScript types
4. Export from component index
5. Add tests with Testing Library

### Modifying xDS Configuration (Gateway Controller)

1. Update models in `pkg/models/`
2. Modify xDS logic in `pkg/xds/`
3. Update control plane in `pkg/controlplane/`
4. Test with Envoy proxy

## Troubleshooting

### Go Services Won't Build
- Verify Go version: `go version`
- Clean module cache: `go clean -modcache`
- Re-download dependencies: `go mod download`

### Node.js Build Failures
- Clear node_modules: `rm -rf node_modules && pnpm install`
- Clear build cache: `rm -rf .next dist`
- Verify Node.js version: `node --version`

### Port Already in Use
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

### HTTPS Certificate Issues (STS)
```bash
cd sts/gate-app
rm server.key server.cert
pnpm ensure-certificates
```

## IDE Setup

### VS Code (Recommended)

**Recommended Extensions:**
- Go (for Go services)
- ESLint (for TypeScript/React)
- Prettier (code formatting)
- TypeScript and JavaScript Language Features
- Vite

**Workspace Settings:**
Check `.vscode/` folder for project-specific settings.

## Additional Resources

- **API Specification:** See `concepts/api-yaml-specification.md`
- **Documentation Guidelines:** See `guidelines/DOCUMENTATION.md`
- **Part-Specific Docs:** Each part has a `README.md` and `spec/` folder
- **Main README:** See root `README.md` for quick start guide

## Getting Help

- Check existing documentation in `/spec` folders
- Review PRDs for feature context
- Consult architecture documents for design decisions
- Check `.serena/memories/` for AI-curated knowledge base