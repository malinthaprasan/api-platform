# api-platform - Epic Breakdown

**Author:** Malintha
**Date:** 2025-11-02
**Project Level:** 2
**Target Scale:** Developer Portal Integration

---

## Epic 1: Developer Portal Integration Foundation

**Expanded Goal:**

Establish the foundational infrastructure for developer-portal integration by implementing configuration management, HTTP client capabilities, and organization synchronization. This epic delivers the core integration mechanism that enables api-platform to communicate with developer-portal and automatically provision organizations with default subscription policies.

**Story Breakdown:**

---

**Story 1.1: Developer Portal Configuration**

As a platform administrator,
I want to configure developer-portal connection settings,
So that api-platform can connect to the developer-portal REST API.

**Acceptance Criteria:**
1. Configuration file or environment variables support the following settings: baseURL, apiKey, timeout, maxRetries, retryBackoffMs
2. Configuration is loaded at application startup
3. Configuration validation ensures required fields (baseURL, apiKey) are present
4. Invalid configuration logs error and prevents application startup
5. Configuration supports multiple developer-portal instances identified by devPortalId

**Prerequisites:** None

---

**Story 1.2: Developer Portal HTTP Client Service**

As a developer,
I want a reusable HTTP client service for developer-portal API calls,
So that I can make REST requests with proper authentication and error handling.

**Acceptance Criteria:**
1. HTTP client service provides methods: GET, POST, PUT, DELETE
2. All requests include `x-wso2-api-key` header from configuration
3. Client sets appropriate Content-Type headers (application/json, multipart/form-data)
4. Client returns standardized response structure with status code, body, and error details
5. Client handles network errors and returns descriptive error messages
6. Client supports timeout configuration from settings
7. Unit tests verify request construction and error handling

**Prerequisites:** Story 1.1 (Configuration)

---

**Story 1.3: Developer Portal Client**

As a developer,
I want a unified client module for all developer-portal operations,
So that I can interact with developer-portal REST APIs for organizations, subscription policies, and APIs.

**Acceptance Criteria:**
1. Client provides `createOrganization(orgData)` method that constructs proper request body with required fields: orgId, orgName, orgHandle, organizationIdentifier
2. `createOrganization` maps optional fields: businessOwner, businessOwnerContact, businessOwnerEmail
3. `createOrganization` calls `POST /devportal/organizations` endpoint and returns created organization response
4. `createOrganization` handles 409 conflict errors (org already exists) gracefully
5. Client provides `createSubscriptionPolicy(orgId, policyData)` method with required fields: policyName, displayName, type, requestCount
6. `createSubscriptionPolicy` supports optional fields: description, billingPlan, timeUnit, unitTime
7. `createSubscriptionPolicy` calls `POST /devportal/organizations/{orgId}/subscription-policies` endpoint
8. Client provides `publishAPI(orgId, apiMetadata, apiDefinition)` method for multipart API publishing
9. Client provides `unpublishAPI(apiId, orgId)` method that calls `DELETE /devportal/apis/{apiId}` endpoint
10. All methods use the HTTP client service for request execution
11. Unit tests verify request mapping, error scenarios, and multipart handling

**Prerequisites:** Story 1.2 (HTTP Client)

---

**Story 1.4: Organization Sync Service with Retry Logic**

As a platform administrator,
I want organization creation to sync to developer-portal with retry capability,
So that organizations are reliably provisioned even during transient failures.

**Acceptance Criteria:**
1. Service provides `syncOrganizationToDevPortal(orgData, devPortalId)` method
2. Method calls developer-portal client to create organization in developer-portal
3. Method implements exponential backoff retry (max 3 attempts) for transient errors (5xx, network timeout)
4. Method does not retry for client errors (4xx except 409)
5. Retry backoff delays: 1s, 2s, 4s
6. Method logs each retry attempt with error details
7. Method returns success with orgId or failure with error details
8. Unit tests verify retry logic and backoff timing

**Prerequisites:** Story 1.3 (Developer Portal Client)

---

**Story 1.5: Subscription Policy Sync Service**

As a platform administrator,
I want a default "Unlimited" subscription policy created automatically after organization sync,
So that APIs can be subscribed to immediately after publishing.

**Acceptance Criteria:**
1. Service provides `syncSubscriptionPolicy(orgId, devPortalId)` method
2. Method creates "Unlimited" policy with: policyName="Unlimited", requestCount=100000, billingPlan="FREE", type="requestCount"
3. Method calls developer-portal client's `createSubscriptionPolicy` method
4. Method implements same retry logic as organization sync (exponential backoff, max 3 attempts)
5. Method logs retry attempts and final status
6. Method returns success with policyID or failure with error details
7. Unit tests verify policy creation and retry behavior

**Prerequisites:** Story 1.4 (Organization Sync Service)

---

**Story 1.6: Organization Creation Hook Integration**

As a platform administrator,
I want developer-portal sync to trigger automatically when I create an organization in api-platform,
So that I don't have to manually provision organizations in developer-portal.

**Acceptance Criteria:**
1. Identify organization creation endpoint/handler in platform-api codebase
2. Add post-creation hook that calls `syncOrganizationToDevPortal` service
3. Hook executes after successful organization creation in api-platform database
4. After successful org sync, hook calls `syncSubscriptionPolicy` to create default policy
5. Hook passes organization data (orgId, orgName, orgHandle, etc.) to sync services
6. Hook logs sync success or failure with correlation ID for tracing
7. If sync fails, organization still exists in api-platform (non-blocking)
8. Hook includes organization creation correlation ID in all logs
9. Integration test verifies org and policy created in both systems

**Prerequisites:** Story 1.5 (Subscription Policy Sync Service)

---

**Story 1.7: Organization Sync Error Handling and Observability**

As a platform operator,
I want detailed logging and error responses for organization sync operations,
So that I can troubleshoot integration issues effectively.

**Acceptance Criteria:**
1. All sync operations log: timestamp, correlationId, orgId, devPortalId, operation, status, duration
2. Failed syncs log full error details: HTTP status, response body, retry attempts
3. Successful syncs log developer-portal orgId confirmation
4. API response includes sync status in organization creation response
5. If sync fails, response includes warning with error summary (org created locally but sync failed)
6. Logs use structured format (JSON) for easy parsing
7. Critical failures (config missing, client creation failure) log at ERROR level
8. Transient failures log at WARN level

**Prerequisites:** Story 1.6 (Hook Integration)

---
