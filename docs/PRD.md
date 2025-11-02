# api-platform Product Requirements Document (PRD)

**Author:** Malintha
**Date:** 2025-11-02
**Project Level:** 2
**Target Scale:** Developer Portal Integration

---

## Goals and Background Context

### Goals

- Enable seamless integration between api-platform and developer-portal to synchronize organization lifecycle and API publishing
- Automate developer-portal provisioning when organizations are created in api-platform, including default subscription policies
- Provide API publishers with a simple publish operation to make APIs available in developer-portal for consumption

### Background Context

The api-platform currently manages API lifecycle and gateway deployments but lacks integration with the developer-portal where APIs are discovered and consumed by external developers. This creates a manual gap where organizations created in api-platform must be manually replicated in developer-portal, and APIs must be manually published.

This integration project addresses this gap by establishing automated synchronization between the two systems. When organizations are created in api-platform, they will automatically be provisioned in developer-portal with default subscription policies (Bronze, Silver, Gold, Unlimited). Additionally, API publishers will have a dedicated endpoint (`POST /apis/{apiId}/publish-to-devportal`) to publish APIs directly to developer-portal, streamlining the API publishing workflow.

---

## Requirements

### Functional Requirements

**FR001:** When a new organization is created in api-platform, the system shall automatically create the corresponding organization in developer-portal using the developer-portal REST API

**FR002:** Organization creation in developer-portal shall use the same organization ID from api-platform as the orgId field to maintain ID consistency across both systems

**FR003:** When an organization is created in developer-portal, the system shall automatically create a default "Unlimited" subscription policy (100000 req/min, FREE billing plan)

**FR004:** The system shall provide a new REST endpoint `POST /apis/{apiId}/publish-to-devportal` that accepts a query parameter `devPortalId` to identify the target developer-portal instance

**FR005:** API publish operation shall validate that the API exists in api-platform before attempting to publish to developer-portal

**FR006:** API publish operation shall retrieve API metadata from api-platform including API name, version, description, endpoints, and OpenAPI definition

**FR007:** API publish operation shall transform api-platform API metadata to developer-portal API metadata format, mapping all required fields (apiName, apiHandle, referenceID, apiVersion, apiType, endpoints)

**FR008:** API publish operation shall submit API metadata and OpenAPI definition to developer-portal using multipart/form-data format with `apiMetadata` and `apiDefinition` fields

**FR009:** The system shall validate that the OpenAPI definition file is named `apiDefinition.json` before submitting to developer-portal

**FR010:** The system shall store developer-portal configuration (base URL, API key, organization mapping) in platform-api configuration

**FR011:** Organization sync shall handle failures gracefully, returning appropriate error responses if developer-portal is unavailable or returns errors

**FR012:** API publish operation shall return success response including the developer-portal API ID when publish is successful

**FR013:** API publish operation shall validate that the target organization exists in both api-platform and developer-portal before publishing

**FR014:** The system shall provide a REST endpoint `DELETE /apis/{apiId}/unpublish-from-devportal` to remove APIs from developer-portal

**FR015:** API unpublish operation shall validate that the API exists in developer-portal before attempting to unpublish

**FR016:** API unpublish operation shall call the developer-portal DELETE API endpoint to completely remove the API from developer-portal

**FR017:** API unpublish operation shall return appropriate error if the API has active subscriptions in developer-portal (cannot delete)

### Non-Functional Requirements

**NFR001:** The developer-portal integration endpoints shall respond within 5 seconds for organization creation and API publish operations under normal conditions

**NFR002:** The system shall implement retry logic with exponential backoff for developer-portal API calls, with a maximum of 3 retry attempts for transient failures

**NFR003:** The system shall log all integration operations (organization sync, API publish) with request/response details to support troubleshooting and audit requirements

---

## User Journeys

**User Journey: Platform Admin Publishes API to Developer Portal**

**Actor:** Platform Admin

**Goal:** Create a new organization and publish an API to the developer portal for external consumption

**Journey Steps:**

1. **Create Organization**
   - Admin calls `POST /organizations` in api-platform with organization details (orgId, orgName, orgHandle, etc.)
   - api-platform creates the organization locally
   - api-platform automatically syncs organization to developer-portal using the same orgId
   - Developer-portal creates organization and provisions default "Unlimited" subscription policy
   - Admin receives success response confirming organization created in both systems

2. **Create/Import API**
   - Admin creates or imports an API in api-platform (existing functionality)
   - API metadata and OpenAPI definition are stored in api-platform

3. **Publish API to Developer Portal**
   - Admin calls `POST /apis/{apiId}/publish-to-devportal?devPortalId={id}`
   - api-platform validates API exists and organization is synced
   - api-platform retrieves API metadata and OpenAPI definition
   - api-platform transforms metadata to developer-portal format
   - api-platform submits API to developer-portal (multipart with apiMetadata + apiDefinition.json)
   - Developer-portal stores API metadata and definition
   - Admin receives success response with developer-portal API ID

4. **Verification**
   - API is now visible in developer-portal for the organization
   - External developers can discover and subscribe to the API using the "Unlimited" policy

**Success Criteria:**
- Organization exists in both systems with matching IDs
- API is published and discoverable in developer-portal
- Default subscription policy is available for API subscriptions

---

## UX Design Principles

Not applicable - backend API integration project with no UI components.

---

## User Interface Design Goals

Not applicable - backend API integration project with no UI components.

---

## Epic List

**Epic 1: Developer Portal Integration Foundation**
- Goal: Establish developer-portal client infrastructure and organization synchronization capability
- Estimated stories: 6-8 stories

**Epic 2: API Publishing to Developer Portal**
- Goal: Enable API publishing from api-platform to developer-portal with metadata transformation
- Estimated stories: 5-7 stories

> **Note:** Detailed epic breakdown with full story specifications is available in [epics.md](./epics.md)

---

## Out of Scope

**Out of Scope:**

- **Bi-directional sync**: Organization updates or deletions in api-platform will not be automatically reflected in developer-portal (one-way sync on creation only)

- **Application/Subscription management**: Developer-portal application subscriptions and lifecycle are managed independently in developer-portal

- **Custom subscription policies**: Only default "Unlimited" policy is created; custom policies must be created manually in developer-portal

- **Multi-tenant developer-portal**: Support for multiple developer-portal instances is limited to configuration; no instance discovery or management

- **API versioning sync**: Updates to existing APIs in api-platform do not automatically update published versions in developer-portal

- **Analytics and monitoring**: No usage analytics or monitoring dashboards for published APIs

- **Developer-portal authentication**: User authentication and authorization for developer-portal are handled independently