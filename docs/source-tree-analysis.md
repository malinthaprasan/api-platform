# Source Tree Analysis

**Generated:** 2025-11-02
**Repository Type:** Monorepo with 4 parts

## Complete Directory Structure

```
api-platform/                          # Root of monorepo
│
├── platform-api/                      # Part 1: Backend Control Plane API
│   ├── src/                          # Go source code
│   │   ├── cmd/                      # Application entry points
│   │   ├── internal/                 # Private application code
│   │   │   ├── handler/              # → HTTP request handlers
│   │   │   ├── service/              # → Business logic layer
│   │   │   ├── repository/           # → Data access layer
│   │   │   ├── model/                # → Data models
│   │   │   ├── dto/                  # → Data transfer objects
│   │   │   ├── database/             # → Database initialization
│   │   │   ├── middleware/           # → HTTP middleware
│   │   │   ├── websocket/            # → WebSocket handlers
│   │   │   ├── utils/                # → Utility functions
│   │   │   └── constants/            # → Application constants
│   │   ├── config/                   # Configuration files
│   │   ├── data/                     # Data files
│   │   ├── resources/                # Static resources
│   │   ├── go.mod                    # Go module definition
│   │   └── main                      # Compiled binary
│   ├── spec/                         # Specifications & planning docs
│   │   ├── architecture.md           # Architecture documentation
│   │   ├── prd.md                    # Product requirements
│   │   ├── design.md                 # Design documentation
│   │   ├── impl.md                   # Implementation tracking
│   │   ├── prds/                     # Individual PRDs
│   │   └── impls/                    # Implementation docs
│   ├── Dockerfile                    # Docker containerization
│   └── README.md                     # Part-specific README
│
├── gateway/                           # Part 2: Gateway Components
│   ├── gateway-controller/           # Envoy xDS control plane
│   │   ├── cmd/                      # Application entry points
│   │   ├── pkg/                      # Public packages
│   │   │   ├── api/                  # → API definitions
│   │   │   │   └── handlers/         # → API handlers
│   │   │   ├── controlplane/         # → Envoy control plane logic
│   │   │   ├── xds/                  # → xDS protocol implementation
│   │   │   ├── models/               # → Data models
│   │   │   ├── storage/              # → Persistence layer (SQLite)
│   │   │   ├── config/               # → Configuration management
│   │   │   ├── logger/               # → Logging utilities
│   │   │   └── utils/                # → Utility functions
│   │   ├── api/                      # OpenAPI specifications
│   │   ├── config/                   # Configuration files
│   │   ├── data/                     # Data files
│   │   ├── resources/                # Static resources
│   │   ├── tests/                    # Test files
│   │   ├── go.mod                    # Go module definition
│   │   ├── Makefile                  # Build automation
│   │   ├── oapi-codegen.yaml         # OpenAPI code generation config
│   │   ├── Dockerfile                # Docker containerization
│   │   └── README.md                 # Component README
│   ├── router/                       # Envoy proxy (separate component)
│   │   └── README.md
│   └── spec/                         # Gateway specifications
│       ├── architecture/
│       ├── design/
│       ├── prds/
│       └── impls/
│
├── portals/                           # Part 3: Web Portals
│   └── management-portal/            # Management UI
│       ├── src/                      # React source code
│       │   ├── pages/                # → Route-based pages (~33 pages)
│       │   ├── components/           # → Reusable UI components (~1,352 files)
│       │   │   └── src/
│       │   │       ├── layouts/      # → Layout components
│       │   │       ├── components/   # → Shared components
│       │   │       └── Icons/        # → Icon components
│       │   └── ...                   # → Additional source folders
│       ├── public/                   # Static assets
│       ├── node_modules/             # Dependencies
│       ├── package.json              # NPM package definition
│       ├── vite.config.ts            # Vite build configuration
│       ├── tsconfig.json             # TypeScript configuration
│       ├── eslint.config.js          # ESLint configuration
│       └── README.md                 # Component README
│
├── sts/                               # Part 4: Security Token Service
│   ├── gate-app/                     # OAuth authentication UI
│   │   ├── app/                      # Next.js app directory
│   │   ├── public/                   # Static assets
│   │   ├── scripts/                  # Build/package scripts
│   │   ├── package.json              # NPM package definition
│   │   ├── next.config.ts            # Next.js configuration
│   │   ├── tsconfig.json             # TypeScript configuration
│   │   ├── server.js                 # Custom server with HTTPS
│   │   └── README.md                 # App README
│   ├── sample-app/                   # Sample OAuth client app
│   │   ├── package.json
│   │   └── README.md
│   ├── packages/                     # Shared packages
│   │   └── oxygen-ui/                # Custom UI component library
│   │       ├── src/
│   │       ├── package.json
│   │       ├── tsconfig.json
│   │       └── README.md
│   ├── spec/                         # STS specifications
│   │   ├── architecture.md
│   │   ├── design.md
│   │   ├── prd.md
│   │   ├── prds/
│   │   └── impls/
│   ├── conf/                         # Configuration files
│   ├── scripts/                      # Setup/deployment scripts
│   ├── kickstart.sh                  # → Automated setup script
│   ├── Dockerfile                    # Docker containerization
│   ├── inputs.yaml                   # Configuration inputs
│   ├── registration.yaml             # Service registration
│   └── README.md                     # STS README
│
├── concepts/                          # Shared concepts & specifications
│   └── api-yaml-specification.md     # Core API YAML format spec
│
├── guidelines/                        # Development guidelines
│   └── DOCUMENTATION.md              # Documentation standards
│
├── docs/                              # Generated documentation (this folder)
│   ├── index.md                      # Master index
│   ├── project-overview.md           # This file
│   ├── source-tree-analysis.md       # Source tree analysis
│   └── project-scan-report.json      # Workflow state
│
├── .github/                           # GitHub configuration
│   └── pull_request_template.md      # PR template
│
├── bmad/                              # BMad workflow system
├── .claude/                           # Claude Code configuration
├── .serena/                           # Serena AI memories
├── .vscode/                           # VS Code settings
├── .git/                              # Git repository
├── .gitignore                         # Git ignore rules
├── LICENSE                            # Apache 2.0 license
├── README.md                          # → Main project README
└── CLAUDE.md                          # Claude-specific instructions
```

## Critical Directories by Part

### Platform API
- **Entry Point:** `platform-api/src/cmd/`
- **API Handlers:** `platform-api/src/internal/handler/`
- **Business Logic:** `platform-api/src/internal/service/`
- **Data Access:** `platform-api/src/internal/repository/`
- **Models:** `platform-api/src/internal/model/`

### Gateway Controller
- **Entry Point:** `gateway/gateway-controller/cmd/`
- **API Handlers:** `gateway/gateway-controller/pkg/api/handlers/`
- **Control Plane:** `gateway/gateway-controller/pkg/controlplane/`
- **xDS Logic:** `gateway/gateway-controller/pkg/xds/`
- **Storage:** `gateway/gateway-controller/pkg/storage/`

### Management Portal
- **Pages:** `portals/management-portal/src/pages/` (33 pages)
- **Components:** `portals/management-portal/src/components/` (1,352 components)
- **Build Config:** `portals/management-portal/vite.config.ts`

### STS
- **Auth UI:** `sts/gate-app/`
- **Sample App:** `sts/sample-app/`
- **UI Library:** `sts/packages/oxygen-ui/`
- **Setup Script:** `sts/kickstart.sh`

## Integration Points

### Platform API ↔ Gateway Controller
- **Protocol:** HTTP REST + WebSocket
- **Direction:** Bidirectional
- **Purpose:** Gateway registration, configuration push, status updates

### Management Portal → Platform API
- **Protocol:** HTTP REST
- **Direction:** Frontend to Backend
- **Purpose:** UI operations, data fetching, CRUD operations

### Platform API ↔ STS
- **Protocol:** OAuth 2.0 / OpenID Connect
- **Direction:** Authentication delegation
- **Purpose:** User authentication, token validation

### Gateway Controller → Envoy Router
- **Protocol:** gRPC (xDS)
- **Direction:** Controller to Proxy
- **Purpose:** Dynamic configuration, route updates, policy enforcement

## Notes

- **Monorepo Structure:** Each part is independently deployable
- **Docker Ready:** All backend services have Dockerfiles
- **Specifications:** Extensive `/spec` folders contain PRDs, architecture, and implementation docs
- **Testing:** Test files located alongside source code
- **Configuration:** Each part manages its own configuration

## Entry Points

| Part | Entry Point | Command |
|------|-------------|---------|
| Platform API | `platform-api/src/cmd/` | Compiled to `main` binary |
| Gateway Controller | `gateway/gateway-controller/cmd/` | Build via Makefile |
| Management Portal | `portals/management-portal/src/main.tsx` | `npm run dev` |
| STS (gate-app) | `sts/gate-app/server.js` | `npm run dev` |

## Build Artifacts

- **Go Binaries:** `main` executables in service roots
- **Web Builds:** Static files in `dist/` or `build/` directories (gitignored)
- **Docker Images:** Built from Dockerfiles in each part root
- **Node Modules:** Dependencies in `node_modules/` (gitignored)