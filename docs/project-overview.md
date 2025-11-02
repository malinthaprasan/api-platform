# WSO2 API Platform - Project Overview

**Generated:** 2025-11-02
**Repository Type:** Monorepo
**Total Parts:** 4

## Executive Summary

The WSO2 API Platform is a comprehensive, AI-ready API management platform that provides full lifecycle management across cloud, hybrid, and on-premises deployments. The platform is built as a monorepo containing multiple independent services and applications working together to deliver a complete API management solution.

## Project Structure

This is a **monorepo** containing 4 distinct parts:

### 1. Platform API (Backend Service)
- **Type:** Backend REST API
- **Technology:** Go 1.24 + Gin Framework
- **Location:** `platform-api/src/`
- **Purpose:** Central control plane API for managing gateways, APIs, projects, and organizations
- **Database:** SQLite3
- **Key Features:** JWT auth, WebSocket support, RESTful endpoints

### 2. Gateway Controller (Backend Service)
- **Type:** Backend Gateway Controller
- **Technology:** Go 1.25.1 + Envoy Control Plane
- **Location:** `gateway/gateway-controller/`
- **Purpose:** Envoy-based API gateway controller implementing xDS protocol
- **Database:** SQLite3
- **Key Features:** Envoy proxy management, gRPC, OpenAPI support

### 3. Management Portal (Web Frontend)
- **Type:** Web Application (SPA)
- **Technology:** React 19 + TypeScript + Vite
- **Location:** `portals/management-portal/`
- **Purpose:** Web UI for managing the entire API platform
- **UI Framework:** Material-UI (MUI) 7.3.4
- **Key Features:** ~1,352 components, i18n support, responsive design

### 4. STS - Security Token Service (Web/Backend)
- **Type:** OAuth/Authentication Service
- **Technology:** Next.js 15.5.2 + TypeScript
- **Location:** `sts/`
- **Purpose:** OAuth authentication and token management
- **Components:** gate-app (auth UI), sample-app (demo), oxygen-ui (component library)
- **Key Features:** SSR, custom UI components, HTTPS support

## Quick Reference

### Technology Stack Summary

| Part | Language | Framework | Database | Deployment |
|------|----------|-----------|----------|------------|
| Platform API | Go 1.24 | Gin | SQLite | Docker |
| Gateway Controller | Go 1.25.1 | Gin + Envoy | SQLite | Docker |
| Management Portal | TypeScript | React 19 + Vite | N/A | Static Build |
| STS | TypeScript | Next.js 15 | N/A | Docker |

### Architecture Patterns

- **Platform API:** Layered architecture (Handler → Service → Repository)
- **Gateway Controller:** Gateway pattern with Envoy xDS control plane
- **Management Portal:** Component-based SPA with Material-UI
- **STS:** SSR application with custom OAuth flows

## Repository Organization

```
api-platform/
├── platform-api/          # Backend: Central control plane API
├── gateway/               # Backend: Envoy gateway controller
├── portals/               # Frontend: Management portal
├── sts/                   # OAuth/Auth: Security token service
├── concepts/              # Core specifications and formats
├── guidelines/            # Development guidelines
└── docs/                  # Generated documentation (this folder)
```

## Key Features

- **GitOps Ready:** Configuration as code for APIs and gateway configs
- **AI-Ready by Design:** MCP-enabled servers for AI agent integration
- **Multi-Deployment:** Supports cloud, on-premises, and hybrid modes
- **Docker-Based:** All components containerized for easy deployment
- **Microservices:** Independent, loosely-coupled services

## Documentation Structure

This documentation is organized to support AI-assisted development:

- **[Architecture Documentation](./architecture-*.md)** - Detailed architecture for each part
- **[Source Tree Analysis](./source-tree-analysis.md)** - Annotated directory structure
- **[Development Guide](./development-guide.md)** - Setup and development instructions
- **[Integration Architecture](./integration-architecture.md)** - How parts communicate

## Getting Started

For development setup, refer to:
- [Development Guide](./development-guide.md)
- Root [README.md](../README.md)

## Related Resources

- **Existing Documentation:** See `spec/` folders in each part for PRDs, architecture docs, and implementation plans
- **API Specification:** `concepts/api-yaml-specification.md`
- **Guidelines:** `guidelines/DOCUMENTATION.md`