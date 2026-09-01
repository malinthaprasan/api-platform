# Feature Specification: OpenAPI Specification Management for REST APIs

**Feature Branch**: `001-openapi-spec-management`

**Created**: 2026-09-01

**Status**: Draft

**Input**: User description: "We currently support creating REST APIs from Platform API by providing individual REST resources etc, but it doesn't allow providing OpenAPI spec to them. Users should be able to upload OpenAPI specs to APIs, retrieve, update them. This is essential when we need to publish the API to the developer portal"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Attach an OpenAPI specification to an existing REST API (Priority: P1)

An API producer has created a REST API in the platform (today, by declaring its resources one at a time). They already maintain an OpenAPI document for that API in their own source repository. They upload that document against the API so the platform holds the API's full, authored contract — descriptions, schemas, examples, security definitions — not just the resource list.

**Why this priority**: Without the ability to attach a specification at all, nothing else in this feature exists. This single capability, on its own, already removes the need to hand-retype an API contract that the producer has already written.

**Independent Test**: Create a REST API, upload a valid OpenAPI document to it, and confirm the platform accepts it and reports the API as having a specification attached. Delivers value on its own: the authored contract is now stored with the API.

**Acceptance Scenarios**:

1. **Given** a REST API exists in a project the caller can access, **When** the caller uploads a valid OpenAPI document in JSON form, **Then** the document is stored against that API and the response confirms success.
2. **Given** a REST API exists, **When** the caller uploads a valid OpenAPI document in YAML form, **Then** it is accepted and stored equivalently to the JSON form.
3. **Given** a REST API exists, **When** the caller uploads a document that is not a parseable OpenAPI document, **Then** the upload is rejected with a message identifying that the document failed validation, and no specification is stored against the API.
4. **Given** a REST API exists, **When** the caller uploads a document exceeding the platform's configured maximum specification size, **Then** the upload is rejected and no partial content is stored.
5. **Given** a REST API in another organization, **When** a caller from a different organization attempts to upload a specification to it, **Then** the request is refused and no change occurs.

---

### User Story 2 - Retrieve the OpenAPI specification of a REST API (Priority: P1)

A consumer-facing surface (the developer portal), a CLI user, or a downstream tool asks the platform for an API's OpenAPI document so it can render documentation, generate a client, or drive a "try it" console.

**Why this priority**: Storing a specification with no way to read it back delivers nothing. Retrieval is what makes the stored contract usable — and is the direct enabler of the stated goal of publishing an API to the developer portal.

**Independent Test**: Upload a specification to an API, then retrieve it and confirm the returned document is semantically equivalent to what was uploaded. Delivers value on its own: the portal and tooling now have a contract to consume.

**Acceptance Scenarios**:

1. **Given** an API with a stored specification, **When** the caller requests that API's specification, **Then** the stored document is returned in full.
2. **Given** an API with a stored specification, **When** the caller requests it in a specific representation (JSON or YAML), **Then** the document is returned in the requested representation with equivalent content.
3. **Given** an API that has no specification attached, **When** the caller requests its specification, **Then** the platform responds that no specification exists for that API, distinguishably from the API itself not existing.
4. **Given** an API in another organization, **When** a caller from a different organization requests its specification, **Then** the request is refused without disclosing whether a specification exists.

---

### User Story 3 - Update or replace an API's OpenAPI specification (Priority: P2)

The API producer changes their contract — adds an operation, corrects a schema, revises descriptions — and re-uploads the revised document so the platform, and everything downstream of it, reflects the current contract.

**Why this priority**: Real APIs change continuously. Without update, producers would have to delete and recreate an API to correct its contract. It ranks below P1 because a single upload plus retrieval is already a usable slice.

**Independent Test**: Upload a specification, upload a revised version of it, retrieve, and confirm the revised content is returned and the prior content is no longer served. Delivers value on its own: the contract stays current without recreating the API.

**Acceptance Scenarios**:

1. **Given** an API with a stored specification, **When** the caller uploads a revised valid specification, **Then** the revised document replaces the previous one and subsequent retrievals return the revised content.
2. **Given** an API with a stored specification, **When** the caller uploads a revised document that fails validation, **Then** the request is rejected and the previously stored specification remains intact and retrievable.
3. **Given** an API with a stored specification, **When** the caller removes the specification, **Then** subsequent retrievals report that no specification exists, and the API itself remains unaffected.

---

### User Story 4 - Publish an API to the developer portal using its specification (Priority: P2)

An API producer publishes an API to the developer portal. The portal renders the API's documentation, operations, schemas, and examples from the specification the producer uploaded, rather than from a minimal resource list.

**Why this priority**: This is the stated business driver for the feature, but it depends on P1 upload and retrieval being in place first.

**Independent Test**: Upload a specification to an API, publish that API to the developer portal, and confirm the portal presents the operations, descriptions, and schemas from the uploaded document.

**Acceptance Scenarios**:

1. **Given** an API with a stored specification, **When** it is published to the developer portal, **Then** the portal presents documentation derived from that specification.
2. **Given** an API with no stored specification, **When** a producer attempts to publish it to the developer portal, **Then** the platform clearly indicates that a specification is required (or presents the API using only its declared resources, per the platform's existing publishing behaviour) rather than failing opaquely.
3. **Given** a published API whose specification is subsequently updated, **When** the update completes, **Then** the portal reflects the updated documentation without the API needing to be republished from scratch.

---

### Edge Cases

- A specification whose declared operations disagree with the API's already-declared resources (paths or methods present in one but not the other).
- A specification that declares server URLs, host names, or backend endpoints that differ from the API's configured backend.
- A specification containing external references (`$ref` to another file or a remote URL) that the platform cannot resolve locally.
- A specification containing an embedded document type declaration or other content that could cause the platform to fetch a remote resource while parsing.
- Two callers uploading a specification to the same API concurrently — the last write must be deterministic and must not leave a half-written document stored.
- A specification that is syntactically valid but semantically empty (no paths declared).
- A specification whose HTTP methods are written in mixed or lower case.
- An extremely large or deeply nested specification submitted to exhaust platform resources.
- A specification uploaded to an API that is currently deployed to one or more gateways — the effect on the running deployment must be defined and must not be an unannounced route change.
- Retrieval of a specification for an API that was deleted between the caller's list and get calls.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Users MUST be able to upload an OpenAPI specification document and associate it with an existing REST API they are authorized to modify.
- **FR-002**: The platform MUST accept specification documents in both JSON and YAML representations and treat them as equivalent.
- **FR-003**: The platform MUST validate an uploaded document as a well-formed OpenAPI document before storing it, and MUST reject a document that fails validation with a message that identifies the failure as a specification-validation failure without disclosing platform internals.
- **FR-004**: The platform MUST reject an uploaded document that exceeds a configurable maximum size, and MUST NOT store partial content on rejection.
- **FR-005**: Users MUST be able to retrieve the stored specification for an API they are authorized to read, in either JSON or YAML representation.
- **FR-006**: Users MUST be able to replace an API's stored specification with a revised document; on validation failure, the previously stored specification MUST remain intact and retrievable.
- **FR-007**: Users MUST be able to remove an API's stored specification without affecting the API itself.
- **FR-008**: The platform MUST distinguish, in its responses, between an API that does not exist, an API that exists but has no specification, and a caller not being permitted to see the API.
- **FR-009**: Every specification operation MUST be authorized against the caller's verified organization and project access for the target API — a caller MUST NOT be able to read, upload, or remove a specification for an API outside their organization.
- **FR-010**: Specification read, upload/update, and removal MUST each be governed by an explicitly named permission consistent with the platform's existing REST API permission naming, so that read access can be granted without write access.
- **FR-011**: The platform MUST NOT resolve external or remote references contained in an uploaded specification while validating or storing it; a document depending on unresolvable external references MUST either be rejected or stored with those references left unresolved, and MUST never cause the platform to issue an outbound request to a location named in the document.
- **FR-012**: The platform MUST record, for each stored specification, when it was last changed and by whom, and MUST surface the last-changed time to callers who can read the specification.
- **FR-013**: The relationship between a stored specification and the API's declared REST resources MUST be defined and applied consistently: [NEEDS CLARIFICATION: does uploading a specification become the source of truth for the API's resources/operations — replacing or synchronising them, and therefore changing gateway routing on the next deployment — or is the specification stored as an independent document used for documentation and portal publishing only, leaving the declared resources untouched?]
- **FR-014**: The platform MUST define whether a REST API can be created directly from a specification: [NEEDS CLARIFICATION: should users be able to create a new REST API by importing an OpenAPI document (and/or by supplying a URL the platform fetches), or is this feature limited to attaching a specification to an API that already exists?]
- **FR-015**: The stored specification MUST be available to the developer portal publishing flow so that a published API's documentation, operations, and schemas are derived from it.
- **FR-016**: An update to a stored specification MUST be reflected to consumers of that specification (including a published portal listing) without requiring the API to be recreated.
- **FR-017**: Concurrent uploads to the same API MUST resolve to exactly one complete stored document; a partially written or interleaved document MUST never be readable.
- **FR-018**: Uploading, updating, or removing a specification MUST NOT silently change the routing of an already-deployed API; any effect on existing deployments MUST require the platform's normal deployment action to take effect.
- **FR-019**: HTTP methods and paths read from an uploaded specification MUST be normalized consistently (case for methods, redundant separators and escaped characters for paths) before any comparison, storage, or matching, so that equivalent documents produce equivalent results.
- **FR-020**: The platform MUST make a specification's presence visible when listing or fetching an API, so a caller can tell which APIs have a specification without fetching each document.

### Key Entities

- **REST API**: The existing platform entity representing a managed REST API within a project and organization. Gains an associated specification and an indicator of whether one is present.
- **OpenAPI Specification Document**: The producer-authored contract for a REST API — its operations, paths, request/response schemas, descriptions, examples, and security declarations. Exactly one current document per REST API. Carries the representation it was supplied in, its size, and when and by whom it was last changed.
- **Project / Organization**: The existing ownership boundary. Every specification is reachable only through the API it belongs to, and therefore only within that API's organization.
- **Developer Portal Listing**: The consumer-facing presentation of a published API, which draws its documentation from the API's specification document.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An API producer can take an OpenAPI document they already maintain and have it attached to an existing API in under 2 minutes, without retyping any part of the contract.
- **SC-002**: A specification retrieved from the platform is semantically equivalent to the document that was uploaded — 100% of operations, paths, schemas, and descriptions are preserved across an upload-then-retrieve round trip, in both JSON and YAML.
- **SC-003**: 95% of specification retrievals return within 1 second for documents up to the platform's supported maximum size.
- **SC-004**: 100% of documents that fail validation are rejected with the API's previously stored specification left unchanged and still retrievable.
- **SC-005**: An API with an attached specification can be published to the developer portal and its documentation rendered from that specification, with no manual re-entry of operations or schemas.
- **SC-006**: Zero cross-organization accesses succeed: no caller can read, replace, or remove a specification belonging to an API outside their own organization, verified by test coverage over every specification operation.
- **SC-007**: Support requests and manual steps relating to "the portal listing doesn't match my API contract" are eliminated for APIs that have an attached specification, because the portal and the contract have a single shared source.

## Assumptions

- The feature targets REST APIs managed by the Platform API. Other API kinds (MCP proxies, LLM proxies, event APIs) are out of scope for this feature.
- OpenAPI 3.0 and 3.1 documents are in scope. Swagger/OpenAPI 2.0 documents and WSDL/SOAP definitions are out of scope for this feature; if support is later needed it is a separate feature.
- Exactly one current specification is stored per REST API. Version history / rollback of specifications is out of scope for this feature; the platform retains only the current document plus its last-changed metadata.
- Callers are already authenticated by the platform's existing mechanism, and organization and project membership are already established. This feature adds authorization rules for the new operations but introduces no new authentication mechanism.
- The maximum accepted specification size is operator-configurable with a safe default, following the platform's existing convention for size-bounded inputs.
- The developer portal already has a publishing flow for platform APIs; this feature supplies the specification that flow consumes rather than building a new publishing pipeline.
- Specification content is treated as untrusted producer-supplied input at every stage — parsing, storage, retrieval, and portal rendering — regardless of the uploading caller's privilege level.
- The `ap` CLI and any management surface that already exposes REST API operations will expose the new specification operations through the platform's normal spec-driven update path; that propagation is not separately enumerated here.
