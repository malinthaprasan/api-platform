## Purpose

Lets a REST API in the Platform API carry the OpenAPI 3.x document that describes its contract, so the API can be published to the developer portal with a definition that consumers can read, try out, and generate clients from.

## ADDED Requirements

### Requirement: Attach an OpenAPI document to a REST API

The system SHALL allow a caller to attach an OpenAPI 3.x document to an existing REST API, and SHALL store the uploaded bytes verbatim. A REST API SHALL have at most one attached document; attaching a document to an API that already has one SHALL replace it. Attaching a document SHALL NOT alter the API's `operations`, `upstream`, `lifeCycleStatus`, or any other field of the API itself.

#### Scenario: Attaching a document to an API that has none

- **WHEN** a caller uploads a valid OpenAPI 3.x document to a REST API that has no document attached
- **THEN** the system stores the document and responds `201` with a body describing the stored document (its content type, byte size, declared OpenAPI version, and the timestamp and user identifier of the upload)

#### Scenario: Replacing an existing document

- **WHEN** a caller uploads a valid OpenAPI 3.x document to a REST API that already has one attached
- **THEN** the system replaces the stored document entirely and responds `200` with a body describing the newly stored document

#### Scenario: The API's own fields are untouched

- **WHEN** a caller uploads a document to a REST API whose `operations` list is non-empty
- **THEN** a subsequent read of that REST API returns exactly the `operations`, `upstream`, and `lifeCycleStatus` it had before the upload

#### Scenario: Attaching to a REST API that does not exist

- **WHEN** a caller uploads a document addressing a REST API identifier that does not exist in the caller's organization
- **THEN** the system responds `404` and stores nothing

#### Scenario: Attaching to a read-only REST API

- **WHEN** a caller uploads a document to a REST API that is marked read-only because it originated from a data-plane gateway
- **THEN** the system responds `409` and stores nothing

### Requirement: Upload format

The system SHALL accept the document as a `multipart/form-data` request carrying exactly one file part whose content is the OpenAPI document in JSON or YAML. The system SHALL determine the stored content type by inspecting the uploaded bytes, and SHALL NOT trust a client-declared content type or filename extension.

#### Scenario: JSON document uploaded

- **WHEN** a caller uploads a file part whose bytes parse as a JSON OpenAPI 3.x document
- **THEN** the system stores it with content type `application/json` regardless of the part's declared content type or file extension

#### Scenario: YAML document uploaded

- **WHEN** a caller uploads a file part whose bytes parse as a YAML (non-JSON) OpenAPI 3.x document
- **THEN** the system stores it with content type `application/yaml` regardless of the part's declared content type or file extension

#### Scenario: Request is not multipart

- **WHEN** a caller sends an upload request whose body is not `multipart/form-data`
- **THEN** the system responds `415` and stores nothing

#### Scenario: Wrong number of file parts

- **WHEN** a caller sends a `multipart/form-data` request containing zero file parts, or more than one
- **THEN** the system responds `400` and stores nothing

### Requirement: OpenAPI 3.x validation

The system SHALL parse every uploaded document and SHALL reject it unless it declares an OpenAPI 3.x version and passes structural validation of that version. A document that fails parsing or validation SHALL NOT be stored, and SHALL NOT replace a previously stored document. Rejection responses SHALL identify that the document was invalid without disclosing internal parser paths, stack traces, or server file locations.

#### Scenario: Well-formed OpenAPI 3.x document

- **WHEN** a caller uploads a document declaring `openapi: 3.0.3` (or any `3.x` version) that passes structural validation
- **THEN** the system accepts and stores it

#### Scenario: Swagger 2.0 document

- **WHEN** a caller uploads a document declaring `swagger: "2.0"`
- **THEN** the system responds `400` with a message stating that only OpenAPI 3.x is supported, and stores nothing

#### Scenario: Structurally invalid OpenAPI document

- **WHEN** a caller uploads a document that declares an OpenAPI 3.x version but fails structural validation
- **THEN** the system responds `400` describing that validation failed, and stores nothing

#### Scenario: Unparseable content

- **WHEN** a caller uploads a file part that is neither valid JSON nor valid YAML
- **THEN** the system responds `400` and stores nothing

#### Scenario: A rejected upload leaves the previous document intact

- **WHEN** a caller uploads an invalid document to a REST API that already has a valid document attached
- **THEN** the system responds `400` and a subsequent retrieval returns the previously stored document unchanged

#### Scenario: Validation failure does not leak internals

- **WHEN** any upload is rejected for being invalid
- **THEN** the response body contains no stack trace, no server file path, and no internal component or module name

### Requirement: Bounded upload size

The system SHALL enforce a maximum upload size for the document, sourced from configuration with a safe default, and SHALL reject a larger upload before reading the whole body into memory. The rejection SHALL NOT disclose the configured limit.

#### Scenario: Upload exceeds the configured maximum

- **WHEN** a caller uploads a document larger than the configured maximum size
- **THEN** the system responds `413` with a generic message, stores nothing, and leaves any previously stored document unchanged

#### Scenario: Rejection does not disclose the limit

- **WHEN** an upload is rejected for exceeding the maximum size
- **THEN** the response body does not state the configured byte ceiling

### Requirement: Verbatim retrieval

The system SHALL return the stored document byte-for-byte as it was uploaded, served with the content type it was stored under. The system SHALL NOT reformat, re-serialize, normalize, resolve external references in, or otherwise rewrite the document on retrieval.

#### Scenario: Retrieving a stored YAML document

- **WHEN** a caller retrieves the document of a REST API to which a YAML document was uploaded
- **THEN** the system responds `200` with `Content-Type: application/yaml` and a body byte-identical to the uploaded bytes, including comments, key order, and whitespace

#### Scenario: Retrieving a stored JSON document

- **WHEN** a caller retrieves the document of a REST API to which a JSON document was uploaded
- **THEN** the system responds `200` with `Content-Type: application/json` and a body byte-identical to the uploaded bytes

#### Scenario: Retrieving when no document is attached

- **WHEN** a caller retrieves the document of a REST API that has none attached
- **THEN** the system responds `404`

#### Scenario: Retrieving metadata without the document body

- **WHEN** a caller requests the document's metadata rather than its content
- **THEN** the system responds `200` with the content type, byte size, declared OpenAPI version, and the timestamps and user identifiers of creation and last update, and without the document bytes

### Requirement: Removing an attached document

The system SHALL allow a caller to remove the document attached to a REST API, after which the API behaves as one that never had a document. Deleting a REST API SHALL also remove its attached document.

#### Scenario: Removing an attached document

- **WHEN** a caller removes the document of a REST API that has one attached
- **THEN** the system responds `204` and a subsequent retrieval responds `404`

#### Scenario: Removing when no document is attached

- **WHEN** a caller removes the document of a REST API that has none attached
- **THEN** the system responds `404`

#### Scenario: Deleting the REST API removes its document

- **WHEN** a REST API with an attached document is deleted
- **THEN** the stored document is removed with it, leaving no orphaned document

### Requirement: Authorization and organization scoping

Every document operation SHALL be authorized against the organization taken from the caller's verified token, never from the request path, query, or body. Read operations SHALL require a read-level REST API scope and write operations an update-level REST API scope, consistent with the scopes guarding the REST API itself. A caller SHALL NOT be able to observe, through status code or response body, whether a REST API exists in an organization other than their own.

#### Scenario: Caller lacks the required scope

- **WHEN** an authenticated caller without an update-level REST API scope attempts to upload or remove a document
- **THEN** the system responds `403` without revealing which scope was required, and stores or removes nothing

#### Scenario: Caller has no valid credentials

- **WHEN** a caller with a missing, expired, or invalid token attempts any document operation
- **THEN** the system responds `401` with the same status and body for every one of those failure causes

#### Scenario: REST API belongs to another organization

- **WHEN** a caller requests a document operation on a REST API identifier that exists only in a different organization
- **THEN** the system responds `404`, identically to how it responds for an identifier that exists nowhere

### Requirement: Existing REST API behavior is unchanged

Adding document support SHALL NOT change the request or response shape of any existing REST API endpoint, and SHALL NOT make attaching a document a precondition for creating, updating, deploying, or deleting a REST API.

#### Scenario: Creating a REST API without a document

- **WHEN** a caller creates a REST API without ever attaching a document
- **THEN** creation succeeds and the API can be updated, deployed, and deleted exactly as before

#### Scenario: Reading a REST API does not embed the document

- **WHEN** a caller reads a REST API that has a document attached
- **THEN** the response body has the same fields it had before this capability existed, and does not embed the document content
