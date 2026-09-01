# Quickstart: OpenAPI Specification Management for REST APIs

**Feature**: `001-openapi-spec-management` | **Date**: 2026-09-01

How the feature behaves once implemented, and how to verify each requirement locally.

## Prerequisites

- A running Platform API with a REST API already created in a project you can access (this feature never creates an API — FR-015).
- A bearer token whose scopes include `ap:rest_api:openapi:read` and `ap:rest_api:openapi:update` (or `ap:rest_api:manage`).
- `PLATFORM_API` set to the base URL, e.g. `http://localhost:9443/api/v0.9`.

## The four operations

```bash
API_ID=my-rest-api        # the API's handle or UUID
TOKEN=...                 # bearer token

# 1. Attach a specification (FR-001). Idempotent — the same call replaces it (FR-006).
curl -sS -X PUT "$PLATFORM_API/rest-apis/$API_ID/openapi" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/yaml" \
  --data-binary @petstore.yaml

# 2. Retrieve it (FR-005), in either representation (FR-002).
curl -sS "$PLATFORM_API/rest-apis/$API_ID/openapi" \
  -H "Authorization: Bearer $TOKEN" -H "Accept: application/json"

curl -sS "$PLATFORM_API/rest-apis/$API_ID/openapi" \
  -H "Authorization: Bearer $TOKEN" -H "Accept: application/yaml"

# 3. See where the specification and the API's declared resources disagree (FR-014).
curl -sS "$PLATFORM_API/rest-apis/$API_ID/openapi/drift" \
  -H "Authorization: Bearer $TOKEN"

# 4. Remove it — the API itself is unaffected (FR-007).
curl -sS -X DELETE "$PLATFORM_API/rest-apis/$API_ID/openapi" \
  -H "Authorization: Bearer $TOKEN"
```

Presence is visible without fetching the document (FR-020):

```bash
curl -sS "$PLATFORM_API/rest-apis/$API_ID" -H "Authorization: Bearer $TOKEN" \
  | jq '{hasOpenApiSpec, openApiSpecUpdatedAt}'
```

## Verifying each requirement

| Check | How | Expected |
|---|---|---|
| FR-002 round trip (SC-002) | `PUT` YAML, `GET` as JSON, `GET` as YAML | Both semantically equal to the upload; the YAML form byte-identical |
| FR-003 validation | `PUT` a file that is not an OpenAPI document | `422`, generic message, no parser internals or file paths |
| FR-004 size cap | `PUT` a file larger than `openapi_spec_max_bytes` | `413`, and the limit value is **not** echoed back |
| FR-006 atomicity | Attach a valid spec, then `PUT` an invalid one, then `GET` | The **first** spec comes back unchanged |
| FR-008 distinguishability | `GET .../openapi` on an API with no spec, vs. on an unknown API id | Two different, distinguishable `404` bodies |
| FR-009 org isolation (SC-006) | Every operation with a token from another organization | Refused; nothing disclosed about whether a spec exists |
| FR-010 split permissions | Token with `:read` only, attempt `PUT` | `403` |
| FR-011 no fetching | `PUT` a document whose `$ref` points at an external URL | Stored or rejected, but **no outbound request** — confirm with a listener on the referenced host that never receives a connection |
| FR-013 / SC-007 routing untouched | Capture the API's runtime artifact, run upload → update → delete, capture again | Byte-identical; declared resources unchanged |
| FR-014 drift | Attach a spec declaring a path the API does not, then `GET .../openapi/drift` | That path listed under `inSpecOnly`; `inSync: false` |
| FR-019 normalization | Upload a spec whose methods are lowercase and paths contain `//` | Drift computed the same as for the normalized form |
| FR-015 out of scope | Look for any create-by-import or import-by-URL surface | None exists |

## Test commands

```bash
cd platform-api

# Unit + integration tests
go test ./internal/handler/... ./internal/service/... ./internal/repository/... ./internal/utils/...

# The existing route/scope coverage gate — fails if a new route has no scope
# declared in resources/openapi.yaml
go test ./internal/server/ -run TestScopeRouteCoverage

# Dependency gates owed by adding kin-openapi (G9)
govulncheck ./...
go-licenses check ./... --disallowed_types=forbidden,restricted
```

## Notes for implementers

- **Regenerate, don't hand-write models.** `api/generated.go` comes from `resources/openapi.yaml` via oapi-codegen (`oapi-codegen.yaml`). Edit the spec, regenerate.
- **The spec file is also the authorization config.** `ScopeRegistry` (`internal/middleware/openapi_scope_registry.go`) parses `resources/openapi.yaml` at startup, so the `security` block *is* the enforcement. A route with no `security` block is unprotected and fails `scope_route_coverage_test.go`.
- **Run the `designing-db-schemas` skill** on `rest_api_specs` before writing the DDL — `data-model.md` proposes a shape consistent with sibling tables, but R1–R10 have not been formally applied.
- **Run the `sync-cli-with-openapi` skill** after the spec change, so the `ap` CLI gains the matching commands.
- **Reuse `normalizePathParams`** from the scope registry for path comparison rather than writing a second normalizer (research R7).
