All paths below are relative to `platform-api/` unless stated otherwise.

## 1. Dependency

- [ ] 1.1 Add `github.com/getkin/kin-openapi v0.133.0` to `go.mod` (pinned to the version already used by `gateway/gateway-controller`), run `go mod tidy`, and verify `go build ./...` succeeds and `go list -m github.com/getkin/kin-openapi` reports v0.133.0
- [ ] 1.2 Clear the dependency-management gate: run `govulncheck ./...` and `go-licenses check ./...` in `platform-api/`, and diff `go mod graph` before/after to list newly introduced indirect modules — verify no high/critical findings and no disallowed license, and record the graph diff for the PR description

## 2. Database schema

- [ ] 2.1 Add `CREATE TABLE IF NOT EXISTS rest_api_openapi_specs` to `internal/database/schema.sql` (SQLite/default) per design D1 — `rest_api_uuid` as PRIMARY KEY, spec bytes, `content_type`, `byte_size`, `openapi_version`, `organization_uuid`, and the `created_by`/`created_at`/`updated_by`/`updated_at` audit set, with `ON DELETE CASCADE` FKs to `rest_apis(uuid)` and `organizations(uuid)`. Verify the file's DDL applies cleanly to a fresh SQLite database
- [ ] 2.2 Mirror the same table into `internal/database/schema.sqlite.sql`, `schema.postgres.sql`, and `schema.sqlserver.sql` with dialect-correct blob types (`BLOB` / `BYTEA` / `VARBINARY(MAX)`) and the `OBJECT_ID`-guarded form for SQL Server. Verify by starting each engine via `it/docker-compose.postgres.yaml` and `it/docker-compose.sqlserver.yaml` and confirming the table is created at startup
- [ ] 2.3 Confirm no `ALTER TABLE` path is needed (no shipped table is modified) and that deleting a `rest_apis` row cascades the spec row away — verify with a repository-level test that deletes an API and asserts zero spec rows remain

## 3. OpenAPI validation package

- [ ] 3.1 Create `internal/openapispec` with a `Validate(ctx, raw []byte) (Metadata, error)` entry point that returns content type, byte size, declared OpenAPI version, and path/operation counts. Verify unit tests cover a valid 3.0 and a valid 3.1 document
- [ ] 3.2 Implement the version gate ahead of full validation: decode only the top-level `openapi`/`swagger` keys with `gopkg.in/yaml.v3`, reject `swagger: "2.0"` with the "OpenAPI 3.x only" message, and reject a missing or non-3.x `openapi` key. Verify with unit tests asserting the distinct Swagger-2.0 message
- [ ] 3.3 Implement JSON-before-YAML content sniffing (design D3) and verify unit tests assert a JSON document is typed `application/json`, a YAML document `application/yaml`, and that a declared part content type or `.yaml` filename never changes the outcome
- [ ] 3.4 Seal the `kin-openapi` loader: `IsExternalRefsAllowed = false` and a `ReadFromURIFunc` that always errors. Verify with a unit test that uploads a document containing an external `http://` `$ref` and asserts no network read is attempted and the document is rejected
- [ ] 3.5 Run validation under a `context.WithTimeout` and verify a unit test with a deeply nested/pathological document returns a timeout error rather than running unbounded
- [ ] 3.6 Sanitize validator error text before it leaves the package (no file paths, stack traces, or internal module names) and verify with a unit test asserting the returned message on a structurally invalid document contains none of those

## 4. Model, repository, service

- [ ] 4.1 Add the spec model type to `internal/model` (bytes plus the metadata and audit fields) and verify `go build ./...` succeeds
- [ ] 4.2 Implement `internal/repository` upsert / get-document / get-metadata-only / delete, every query filtered by `organization_uuid` and keyed by `rest_api_uuid`. Verify with repository tests that a get-metadata query does not select the blob column and that a query for another organization's UUID returns no row
- [ ] 4.3 Implement the service layer in `internal/service`: resolve the REST API by handle scoped to the organization from verified claims **before** any spec lookup, reject a read-only (`origin = gateway_api`) API with a conflict error, and enforce the byte ceiling as a second check after read. Verify with service tests covering first attach, replace, read-only rejection, and cross-organization resolution
- [ ] 4.4 Verify that a rejected upload leaves a previously stored document byte-identical — service test that stores a valid document, uploads an invalid one, and re-reads

## 5. Configuration

- [ ] 5.1 Add `max_upload_bytes` under a new `openapi_spec` section in `config/config.go` with a 2 MiB default in `config/default_config.go`, and document it in `config/config-template.toml`. Verify `config/config_test.go`-style tests assert the default applies when the key is absent and that an explicit value overrides it

## 6. API contract and generated code

- [ ] 6.1 Add the four operations from design D5 to `resources/openapi.yaml` under `/rest-apis/{restApiId}/openapi` (+ `/metadata`), each declaring its `ap:rest_api:*` scopes, with the `multipart/form-data` request body and the `type: string, format: binary` document response. Verify the file parses and that `middleware.LoadScopeRegistry` picks up all four operations at startup
- [ ] 6.2 Add `PayloadTooLarge` (413) and `UnsupportedMediaType` (415) entries under `components/responses` and reference them from the new operations. Verify no existing operation's responses changed
- [ ] 6.3 Map the new operations' scopes in `resources/role-to-scope-mapping.yaml` and verify startup validation of the role-to-scope map passes (no wildcard segment, every scope declared)
- [ ] 6.4 Run `make generate` and verify `api/generated.go` regenerates cleanly with no unrelated diff; move any type unreachable from a path into `api/manual_types.go` per the existing convention

## 7. Handler and routing

- [ ] 7.1 Create `internal/handler/api_openapi_spec.go` with the four handlers, registered via a `RegisterRoutes(mux router.Router)` on `constants.APIBasePath + "/rest-apis/{restApiId}/openapi"` and wrapped in `middleware.MapErrors`, and wire it up where the other REST API handlers are constructed. Verify the routes appear in the router recorder via a `internal/router/router_test.go`-style assertion
- [ ] 7.2 Implement the upload handler: reject a non-multipart body with `415`, wrap the body in `http.MaxBytesReader` with the configured limit **before** `ParseMultipartForm`, reject zero or more-than-one file parts with `400`, and return `413` with a generic message that never states the limit. Verify with handler tests for each of those four rejections
- [ ] 7.3 Implement the document retrieval handler to write stored bytes directly with the stored content type — no re-serialization. Verify a handler test uploads a YAML document with comments and irregular key order and asserts the response body is byte-identical
- [ ] 7.4 Implement the metadata and delete handlers, returning `200` with JSON metadata and `204` respectively, and `404` when no document is attached. Verify with handler tests
- [ ] 7.5 Verify status-code correctness for attach vs replace: a handler test asserting `201` on first attach and `200` on replace

## 8. Security verification

- [ ] 8.1 Verify organization isolation with a handler test: a caller from organization A requesting each of the four operations on an API belonging to organization B receives `404`, identical to a nonexistent identifier
- [ ] 8.2 Verify scope enforcement with handler tests: a caller holding only `ap:rest_api:read` gets `403` on upload and delete; a caller with no valid token gets an identical `401` status and body for a missing, an expired, and a malformed token
- [ ] 8.3 Verify no error response from the new endpoints contains a stack trace, server file path, internal module name, or the configured byte ceiling — assert across the `400`/`413`/`415`/`403`/`404` paths

## 9. Integration and regression

- [ ] 9.1 Add an integration test in the `internal/handler/*_integration_test.go` style covering the full lifecycle against a real database: create API → upload → retrieve verbatim → retrieve metadata → replace → delete → delete API. Verify it passes under `make it-sqlite`, `make it-postgres`, and `make it-sqlserver`
- [ ] 9.2 Verify existing REST API behavior is unchanged: an integration assertion that `GET /rest-apis/{id}` returns the same field set as before (no embedded document) and that an API with a document attached still deploys, and that `make test` passes with no modification to existing REST API tests

## 10. Documentation

- [ ] 10.1 Document the new endpoints and the `openapi_spec.max_upload_bytes` setting wherever REST API endpoints are already documented, and note in the PR description that the stored document is documentation-only and does not drive routing
