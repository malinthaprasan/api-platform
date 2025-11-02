# Integration Architecture

**Generated:** 2025-11-02
**Repository Type:** Monorepo

## Overview

This document describes how the four parts of the WSO2 API Platform integrate and communicate with each other.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Control Plane                             │
│                                                                  │
│  ┌──────────────────┐          ┌─────────────────┐             │
│  │  Management      │   HTTP   │  Platform API   │             │
│  │  Portal          │ ────────▶│  (Go Backend)   │             │
│  │  (React SPA)     │  REST    │                 │             │
│  └──────────────────┘          │  • Gateway Mgmt │             │
│          │                      │  • API Mgmt     │             │
│          │                      │  • Project Mgmt │             │
│          │ HTTP/REST            └────────┬────────┘             │
│          │                               │                      │
│          │                               │ WebSocket            │
│          │                               │ (Real-time)          │
│          │                               │                      │
│          └───────────────────────────────┘                      │
│                                                                  │
│                        ┌──────────────┐                         │
│                        │     STS      │                         │
│                        │   (OAuth)    │                         │
│                        └──────┬───────┘                         │
│                               │ OAuth 2.0                       │
│                               │ OIDC                            │
└───────────────────────────────┼─────────────────────────────────┘
                                │
                                │ Authentication
                                │
┌───────────────────────────────┼─────────────────────────────────┐
│                        Data Plane                                │
│                                │                                 │
│  ┌─────────────────────────────▼──────────────────┐             │
│  │         Gateway Controller                      │             │
│  │         (Go + Envoy Control Plane)              │             │
│  │                                                  │             │
│  │  • xDS Server                                   │             │
│  │  • Configuration Management                     │             │
│  │  • Policy Enforcement                           │             │
│  └──────────────────┬───────────────────────────────┘            │
│                     │ gRPC (xDS Protocol)                        │
│                     │                                            │
│  ┌─────────────────▼───────────────────────────┐                │
│  │         Envoy Proxy (Router)                │                │
│  │                                              │                │
│  │  • HTTP/HTTPS Routing                       │                │
│  │  • Policy Execution                         │                │
│  │  • Rate Limiting                            │                │
│  └──────────────────────────────────────────────┘               │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Integration Points

### 1. Management Portal → Platform API

**Protocol:** HTTP REST
**Direction:** Frontend → Backend
**Authentication:** JWT tokens (via STS)

**API Endpoints:**
- Gateway management operations
- API lifecycle operations
- Project/organization management
- User management

**Data Flow:**
1. User interacts with Management Portal UI
2. Portal makes REST API calls to Platform API
3. Platform API processes requests and returns JSON responses
4. Portal updates UI based on responses

### 2. Platform API ↔ Gateway Controller

**Protocol:** HTTP REST + WebSocket
**Direction:** Bidirectional
**Authentication:** Gateway key-based authentication

**Communication Patterns:**

**A. Gateway Registration (HTTP POST)**
```
Gateway Controller → Platform API
POST /api/v1/gateways/register
{
  "gatewayKey": "...",
  "hostname": "...",
  "capabilities": [...]
}
```

**B. Configuration Push (HTTP POST)**
```
Platform API → Gateway Controller
POST /api/v1/gateway/{id}/config
{
  "apis": [...],
  "policies": [...],
  "routes": [...]
}
```

**C. Status Updates (WebSocket)**
```
Gateway Controller → Platform API
ws://platform-api/ws/gateway/{id}
{
  "type": "status",
  "status": "healthy",
  "metrics": {...}
}
```

**D. Real-Time Events (WebSocket)**
- API deployment status
- Health check results
- Metrics updates
- Error notifications

### 3. Platform API ↔ STS

**Protocol:** OAuth 2.0 / OpenID Connect
**Direction:** Authentication delegation
**Flow:** Authorization Code Flow with PKCE

**Integration Points:**

**A. User Login**
1. User accesses Management Portal
2. Portal redirects to STS (gate-app)
3. User authenticates via STS OAuth UI
4. STS issues tokens (access + refresh)
5. Portal receives tokens and stores securely
6. Portal uses access token for Platform API calls

**B. Token Validation**
```
Platform API → STS
POST /oauth/introspect
{
  "token": "access_token"
}

Response:
{
  "active": true,
  "sub": "user_id",
  "scope": "api.read api.write"
}
```

**C. Token Refresh**
```
Management Portal → STS
POST /oauth/token
{
  "grant_type": "refresh_token",
  "refresh_token": "..."
}
```

### 4. Gateway Controller → Envoy Proxy

**Protocol:** gRPC (xDS - Envoy Discovery Service)
**Direction:** Controller → Proxy
**Purpose:** Dynamic configuration management

**xDS APIs Used:**
- **LDS (Listener Discovery Service):** HTTP listeners configuration
- **RDS (Route Discovery Service):** Route configurations
- **CDS (Cluster Discovery Service):** Upstream clusters
- **EDS (Endpoint Discovery Service):** Endpoint/backend configurations

**Configuration Flow:**
1. Platform API pushes API config to Gateway Controller
2. Gateway Controller translates to Envoy xDS format
3. Gateway Controller serves xDS via gRPC
4. Envoy Proxy connects and receives configuration
5. Envoy applies configuration dynamically (zero-downtime)

## Data Flow Examples

### Example 1: Deploying an API

```
1. User (Management Portal)
   ↓ REST: POST /apis
2. Platform API
   - Validates API specification
   - Stores in database (SQLite)
   - Pushes to Gateway Controller
   ↓ REST: POST /gateway/{id}/config
3. Gateway Controller
   - Receives API configuration
   - Translates to Envoy xDS format
   - Stores in local database
   ↓ gRPC: xDS update
4. Envoy Proxy
   - Receives configuration update
   - Applies routes dynamically
   - API is now live
   ↓ WebSocket: Deployment success
5. Platform API
   ↓ WebSocket: Real-time update
6. Management Portal
   - Shows "Deployed" status to user
```

### Example 2: User Authentication

```
1. User → Management Portal
   - Clicks "Login"
2. Management Portal → STS (gate-app)
   - Redirects to /oauth/authorize
3. User → STS
   - Enters credentials
   - Authenticates
4. STS → Management Portal
   - Returns authorization code
5. Management Portal → STS
   - Exchanges code for tokens
   - POST /oauth/token
6. STS → Management Portal
   - Returns access_token + refresh_token
7. Management Portal → Platform API
   - All API calls include: Authorization: Bearer {access_token}
8. Platform API → STS (optional)
   - Validates token via introspection
   - Caches validation results
```

### Example 3: Real-Time Gateway Monitoring

```
1. Management Portal
   ↓ HTTP connection
2. Platform API
   - Polls for gateway events
3. Gateway Controller
   ↓ WebSocket connection
4. Platform API
   - Receives health/metrics from gateway
5. Platform API
   ↓ WebSocket broadcast
6. Management Portal
   - Receives updates via HTTP polling
   - Updates dashboard UI
```

## Shared Data Models

### API Configuration
Shared between Platform API and Gateway Controller:
```yaml
version: api-platform.wso2.com/v1
kind: http/rest
data:
  name: string
  version: string
  context: string
  upstream:
    - url: string
  operations:
    - method: string
      path: string
      requestPolicies: []
      responsePolicies: []
```

### Gateway Metadata
```json
{
  "id": "gateway-uuid",
  "name": "dev-gateway",
  "hostname": "localhost",
  "type": "basic|standard",
  "status": "healthy|degraded|offline",
  "capabilities": ["rate-limiting", "auth", "analytics"]
}
```

## Database Isolation

Each part maintains its own database:

- **Platform API:** SQLite database for:
  - Gateway registrations
  - API definitions
  - Projects and organizations
  - User data (if not delegated)

- **Gateway Controller:** SQLite database for:
  - API configurations (synced from Platform API)
  - xDS snapshots
  - Local cache
  - Metrics

- **No Shared Database:** Parts communicate via APIs only

## Security Considerations

### Authentication Flow
- All user authentication goes through STS (OAuth 2.0)
- JWT tokens used for API authentication
- Tokens validated by Platform API (introspection or JWT verification)

### Gateway Security
- Gateway Controller authenticates with Platform API using pre-shared keys
- Mutual TLS can be enabled for production deployments
- WebSocket connections secured with tokens

### API Security
- HTTPS/TLS for all HTTP communication
- gRPC with TLS for xDS communication
- Secure WebSocket (wss://) for real-time updates

## Deployment Considerations

### Independent Deployment
- Each part can be deployed independently
- Docker containers for all services
- No tight coupling between parts

### Hybrid Deployment
- **Control Plane (Cloud):**
  - Platform API
  - Management Portal
  - STS

- **Data Plane (On-Premises):**
  - Gateway Controller
  - Envoy Proxy

### Scaling
- **Platform API:** Horizontal scaling with load balancer
- **Gateway Controller:** One per gateway instance
- **Management Portal:** Static files on CDN
- **STS:** Horizontal scaling with session affinity

## Communication Protocols Summary

| From | To | Protocol | Purpose |
|------|-----|----------|---------|
| Management Portal | Platform API | HTTP REST | CRUD operations |
| Management Portal | STS | OAuth 2.0/OIDC | Authentication |
| Platform API | Gateway Controller | WebSocket | Configuration push |
| Platform API | Gateway Controller | WebSocket | Status monitoring |
| Platform API | STS | HTTP REST | Token validation |
| Gateway Controller | Envoy Proxy | gRPC (xDS) | Dynamic config |

## Error Handling

### Gateway Offline
- Platform API marks gateway as offline
- WebSocket sends notification to Management Portal
- User sees warning in UI

### Authentication Failure
- STS rejects invalid credentials
- Management Portal shows error
- User prompted to retry

### Configuration Push Failure
- Platform API retries configuration push
- Gateway Controller logs error
- Status reported back to Portal via WebSocket

## Future Enhancements

- Message queue (Redis/NATS) for async communication
- Event streaming for analytics
- Multi-region gateway coordination