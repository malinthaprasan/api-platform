# WSO2 API Platform - Documentation Index

**Generated:** 2025-11-02
**Last Updated:** 2025-11-02
**Scan Level:** Quick (Pattern-based)
**Repository Type:** Monorepo with 4 parts

---

## Project Overview

- **Type:** Monorepo with 4 independent parts
- **Primary Languages:** Go (backend), TypeScript (frontend)
- **Architecture:** Microservices with independent deployment

### Quick Reference

| Part | Type | Tech Stack | Location |
|------|------|------------|----------|
| **Platform API** | Backend | Go 1.24 + Gin + SQLite | `platform-api/src/` |
| **Gateway Controller** | Backend | Go 1.25.1 + Envoy + gRPC | `gateway/gateway-controller/` |
| **Management Portal** | Web | React 19 + Vite + MUI | `portals/management-portal/` |
| **STS** | Web/Auth | Next.js 15 + TypeScript | `sts/` |

---

## Generated Documentation

### Core Documentation

- **[Project Overview](./project-overview.md)** - Executive summary, project structure, and key features
- **[Source Tree Analysis](./source-tree-analysis.md)** - Annotated directory structure with integration points
- **[Development Guide](./development-guide.md)** - Setup instructions, build commands, and development workflow
- **[Integration Architecture](./integration-architecture.md)** - How parts communicate and integrate

### Part-Specific Documentation

#### Platform API (Backend)
- **Architecture:** [platform-api/spec/architecture.md](../platform-api/spec/architecture.md)
- **PRD:** [platform-api/spec/prd.md](../platform-api/spec/prd.md)
- **Design:** [platform-api/spec/design.md](../platform-api/spec/design.md)
- **Implementation:** [platform-api/spec/impl.md](../platform-api/spec/impl.md)
- **README:** [platform-api/README.md](../platform-api/README.md)

**Key PRDs:**
- [Gateway Management](../platform-api/spec/prds/gateway-management.md)
- [API Lifecycle Management](../platform-api/spec/prds/api-lifecycle-management.md)
- [Project Management](../platform-api/spec/prds/project-management.md)
- [Platform Bootstrap](../platform-api/spec/prds/platform-bootstrap.md)
- [Organization Management](../platform-api/spec/prds/organization-management.md)
- [Gateway WebSocket Events](../platform-api/spec/prds/gateway-websocket-events.md)

#### Gateway Controller (Backend)
- **Architecture:** [gateway/spec/architecture/architecture.md](../gateway/spec/architecture/architecture.md)
- **Spec:** [gateway/spec/spec.md](../gateway/spec/spec.md)
- **PRD:** [gateway/spec/prd.md](../gateway/spec/prd.md)
- **Design:** [gateway/spec/design/design.md](../gateway/spec/design/design.md)
- **Implementation:** [gateway/spec/impl.md](../gateway/spec/impl.md)
- **README:** [gateway/gateway-controller/README.md](../gateway/gateway-controller/README.md)

**Key PRDs:**
- [XDS Server](../gateway/spec/prds/xds-server.md)
- [API Configuration Management](../gateway/spec/prds/api-configuration-management.md)
- [SQLite Persistence](../gateway/spec/prds/sqlite-persistence.md)
- [Policy Engine](../gateway/spec/prds/policy-engine.md)
- [Zero Downtime Updates](../gateway/spec/prds/zero-downtime-updates.md)

#### Management Portal (Web)
- **README:** [portals/management-portal/README.md](../portals/management-portal/README.md)
- **Component Count:** ~1,352 components
- **Pages:** ~33 route-based pages

#### STS - Security Token Service (Web/Auth)
- **Architecture:** [sts/spec/architecture.md](../sts/spec/architecture.md)
- **Design:** [sts/spec/design.md](../sts/spec/design.md)
- **PRD:** [sts/spec/prd.md](../sts/spec/prd.md)
- **Implementation:** [sts/spec/impl.md](../sts/spec/impl.md)
- **README:** [sts/README.md](../sts/README.md)

**Key PRDs:**
- [OAuth Support](../sts/spec/prds/oauth-support.md)
- [Authentication UI](../sts/spec/prds/authentication-ui.md)
- [Single Container Deployment](../sts/spec/prds/single-container-deployment.md)
- [Automated Setup](../sts/spec/prds/automated-setup.md)
- [Sample OAuth App](../sts/spec/prds/sample-oauth-app.md)

---

## Existing Documentation

### Project-Level

- **[README.md](../README.md)** - Main project README with quick start guide
- **[API YAML Specification](../concepts/api-yaml-specification.md)** - Core API definition format
- **[Documentation Guidelines](../guidelines/DOCUMENTATION.md)** - Documentation standards
- **[Pull Request Template](../.github/pull_request_template.md)** - PR template

### Component Specifications (Reference)

The following components have comprehensive specification documentation but are not yet implemented in this repository:

- **API Designer:** [api-designer/spec/](../api-designer/spec/)
  - [Specification](../api-designer/spec/spec.md)
  - [Architecture](../api-designer/spec/architecture/architecture.md)
  - [Design](../api-designer/spec/design/design.md)
  - [Use Cases](../api-designer/spec/use-cases/use_cases.md)

- **API Portal:** [api-portal/spec/](../api-portal/spec/)
  - [Specification](../api-portal/spec/spec.md)
  - [Architecture](../api-portal/spec/architecture/architecture.md)
  - [Design](../api-portal/spec/design/design.md)
  - [Use Cases](../api-portal/spec/use-cases/use_cases.md)

- **Enterprise Portal:** [enterprise-portal/spec/](../enterprise-portal/spec/)
  - [Specification](../enterprise-portal/spec/spec.md)
  - [Architecture](../enterprise-portal/spec/architecture/architecture.md)
  - [Design](../enterprise-portal/spec/design/design.md)
  - [Use Cases](../enterprise-portal/spec/use-cases/use_cases.md)

- **CLI:** [cli/spec/](../cli/spec/)
  - [Specification](../cli/spec/spec.md)
  - [Architecture](../cli/spec/architecture/architecture.md)
  - [Design](../cli/spec/design/design.md)
  - [Use Cases](../cli/spec/use-cases/use_cases.md)

- **Management Portal Spec:** [management-portal/spec/](../management-portal/spec/)
  - [Specification](../management-portal/spec/spec.md)
  - [Architecture](../management-portal/spec/architecture/architecture.md)
  - [Design](../management-portal/spec/design/design.md)
  - [Use Cases](../management-portal/spec/use-cases/use_cases.md)

---

## Getting Started

### For New Developers

1. **Read:** [Project Overview](./project-overview.md) - Understand the high-level architecture
2. **Read:** [Source Tree Analysis](./source-tree-analysis.md) - Understand the codebase structure
3. **Follow:** [Development Guide](./development-guide.md) - Set up your development environment
4. **Review:** Part-specific READMEs for detailed component information

### For AI-Assisted Development

When working on features:

1. **Reference this index** for navigation
2. **Consult integration architecture** for cross-component changes
3. **Check existing PRDs** for context on implemented features
4. **Follow documentation guidelines** when creating new docs

### Quick Start Commands

```bash
# Platform API
cd platform-api/src && go run cmd/main.go

# Gateway Controller
cd gateway/gateway-controller && go run cmd/main.go

# Management Portal
cd portals/management-portal && pnpm dev

# STS
cd sts/gate-app && pnpm dev
```

---

## Architecture Highlights

### Clean Separation
- Each part is independently deployable
- No shared databases between parts
- Communication via well-defined APIs

### Technology Choices
- **Go** for high-performance backend services
- **React/Next.js** for modern web applications
- **Envoy** for robust API gateway capabilities
- **SQLite** for embedded persistence
- **Docker** for containerization

### Key Patterns
- **Platform API:** Layered architecture (Handler → Service → Repository)
- **Gateway Controller:** Envoy xDS control plane pattern
- **Management Portal:** Component-based SPA with Material-UI
- **STS:** OAuth 2.0 server with Next.js SSR

---

## Integration Summary

```
Management Portal (React)
    ↓ HTTP REST + WebSocket
Platform API (Go)
    ↓ HTTP REST + WebSocket
Gateway Controller (Go)
    ↓ gRPC (xDS)
Envoy Proxy

           ↓ OAuth 2.0
          STS (Next.js)
```

See [Integration Architecture](./integration-architecture.md) for detailed communication patterns.

---

## File Organization

### Documentation
- **`/docs`** - Generated documentation (this folder)
- **`/concepts`** - Core specifications
- **`/guidelines`** - Development guidelines
- **`/spec`** folders - Part-specific specs and PRDs

### Code
- **`platform-api/src`** - Platform API source code
- **`gateway/gateway-controller`** - Gateway controller source code
- **`portals/management-portal/src`** - Management portal source code
- **`sts/`** - STS applications and packages

### Configuration
- **`.github`** - GitHub workflows and templates
- **`.claude`** - Claude Code configuration and agents
- **`.serena`** - Serena AI memories
- **`bmad`** - BMad workflow system

---

## Development Resources

### Code Quality
- All Go services use `go fmt` and standard Go conventions
- TypeScript services use ESLint and Prettier
- Testing frameworks in place for all parts

### CI/CD
- Docker support for all backend services
- Build automation via Makefiles (Go) and package.json scripts (Node)

### Knowledge Base
- **`.serena/memories/`** - AI-curated knowledge base with:
  - Project overview
  - Repository structure
  - Tech stack details
  - STS overview
  - API YAML format
  - Documentation conventions
  - Task completion checklists
  - Suggested commands

---

## Next Steps

### For Brownfield Development
1. Review this index for navigation
2. Check relevant PRDs for feature context
3. Consult architecture docs for design decisions
4. Follow development guide for setup

### For New Features
1. Create PRD in appropriate `spec/prds/` folder
2. Update architecture docs as needed
3. Implement following established patterns
4. Add tests and documentation

### For Bug Fixes
1. Identify affected part
2. Review related architecture/design docs
3. Check existing tests
4. Fix and add regression tests

---

## Contact & Support

- **License:** Apache 2.0 (see [LICENSE](../LICENSE))
- **Copyright:** WSO2 Inc. (2012-2025)

---

**This documentation was generated by the BMad document-project workflow.**
**For questions about documentation structure, refer to the workflow system in `/bmad`.**