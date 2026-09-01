## Why

A REST API in the Platform API is defined today only by structured fields — `displayName`, `context`, `version`, `upstream`, and a list of `operations` (`platform-api/resources/openapi.yaml`, `/rest-apis`). There is no way to attach the API's OpenAPI document. The API Portal, in contrast, requires a `definition` file when an API is created or updated (`portals/api-portal/docs/api-portal-openapi-spec-v0.9.yaml`, `POST /apis`), so an API authored in the Platform API cannot be published to the developer portal with the contract consumers need to read, try out, or generate clients from.

## What Changes

- Add an OpenAPI spec sub-resource to a REST API: upload/replace, retrieve, and delete the document attached to `/rest-apis/{restApiId}`.
- Uploads are `multipart/form-data` with a single spec file (JSON or YAML); retrieval returns the document **verbatim**, byte-for-byte as uploaded, with the content type it was stored under.
- Uploaded documents are parsed and structurally validated as **OpenAPI 3.x**. OpenAPI 2.0 (Swagger) and anything that fails validation are rejected with `400` — no conversion, no partial acceptance.
- The stored spec is **documentation only**. It does not derive, replace, or reconcile the API's `operations` list, and it has no effect on gateway routing or deployment. Operations remain managed exclusively through the existing REST API endpoints.
- The spec is scoped to the API's organization and project, follows the existing `ap:rest_api:*` scope model, and is deleted with the API.
- Upload size is bounded by configuration; oversize documents are rejected with `413` and a generic message.

Not breaking: every existing REST API endpoint keeps its current request and response shape. An API with no spec attached behaves exactly as it does today.

## Capabilities

### New Capabilities
- `rest-api-openapi-spec`: attaching an OpenAPI 3.x document to a REST API in the Platform API — upload/replace, verbatim retrieval, deletion, validation rules, size bounds, and authorization.

### Modified Capabilities
<!-- None. openspec/specs/ has no existing capability specs; no existing documented
     requirement changes, since REST API create/update/list/delete behavior is unchanged. -->

## Impact

- **API surface**: new paths in `platform-api/resources/openapi.yaml` under `/rest-apis/{restApiId}/openapi`, guarded by the existing `ap:rest_api:read`/`ap:rest_api:update`/`ap:rest_api:manage` scopes and declared in `platform-api/resources/role-to-scope-mapping.yaml`.
- **Code**: new handler, service, and repository paths in `platform-api/internal/{handler,service,repository,model,dto}`; router registration in `platform-api/internal/router`.
- **Database**: a new table storing the spec bytes, content type, and audit columns, keyed to the REST API's artifact UUID — additive only, across all four dialects (`schema.sql`, `schema.sqlite.sql`, `schema.postgres.sql`, `schema.sqlserver.sql`) per `.claude/rules/db-schema-changes.md`.
- **Dependencies**: an OpenAPI 3.x parser/validator for Go. `github.com/getkin/kin-openapi v0.133.0` is already used elsewhere in this repo (`gateway/gateway-controller`, `event-gateway/gateway-controller`) and is the natural choice; it is new to `platform-api/go.mod` and must clear `.claude/rules/dependency-management.md`.
- **Config**: a new bounded upload-size setting in `platform-api/config`.
- **Downstream**: unblocks publishing a Platform-API-authored REST API to the API Portal with its definition. Wiring that publish flow is a separate change.
- **Non-goals**: MCP proxies and LLM proxies, deriving `operations` from the spec, import-by-URL, Swagger 2.0 conversion, and spec versioning/history.
