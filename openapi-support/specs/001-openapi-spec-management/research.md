# Phase 0 Research: OpenAPI Specification Management for REST APIs

**Feature**: `001-openapi-spec-management` | **Date**: 2026-09-01

All findings below were verified against the codebase at commit `996492ec1` (branch `test-spec-kit`). File references are repo-relative.

---

## R1. Where REST API state already lives

**Decision**: Store the specification in a **new** `rest_api_specs` table, one row per REST API, rather than in `rest_apis.configuration`.

**Rationale**:
- `rest_apis` (`platform-api/internal/database/schema.sqlite.sql:79`) is a shipped, GA table. Under `db-schema-changes.md` Directive 1 it is frozen to additive changes only, and a new table needs only its guarded `CREATE TABLE IF NOT EXISTS` — no per-dialect `ALTER TABLE` upgrade path (Directive 3). A new column on `rest_apis` would need one for all three dialects.
- `rest_apis.configuration` (BLOB) deserializes into `model.RestAPIConfig` (`platform-api/internal/model/api.go:45`) and is the input to the gateway runtime artifact. Putting spec bytes there would couple the specification to routing — directly contradicting FR-013 — and would load a multi-hundred-KB document on every API read and list.

**Alternatives considered**:
- Column on `rest_apis` — rejected: needs a three-dialect `ALTER` path, and loads spec bytes into every list query.
- Reuse the `artifacts` table — rejected: `artifacts` carries only identity/type/org, and the artifact registry (`internal/repository/artifact_tables.go`) treats it as a kind index, not a content store.

---

## R2. OpenAPI parsing and validation library

**Decision**: Use `github.com/getkin/kin-openapi v0.133.0` — already pinned at this exact version in `gateway/gateway-controller/go.mod:8` and `event-gateway/gateway-controller/go.mod:8`.

**Rationale**: `dependency-management.md` is satisfied more cheaply by an in-repo, already-vetted version than by introducing a new module. Matching the existing pin also avoids two versions of the same parser in one workspace (`go.work`).

**Important caveat — OpenAPI 3.1**: kin-openapi targets OpenAPI **3.0.x**. Its 3.1 handling is incomplete (3.1's JSON-Schema-2020-12 keywords, `type` arrays, and webhook objects are not fully modeled). The spec's Assumptions currently claim 3.0 **and** 3.1 are in scope.

> **Open decision for the user** — three viable resolutions, none silently chosen here:
> - **(a)** Narrow the spec's assumption to OpenAPI 3.0.x for structural validation; accept and store 3.1 documents after lighter checks (parseable, declares `openapi: 3.x`, has `paths`), flagging that 3.1 validation is shallower.
> - **(b)** Add `github.com/pb33f/libopenapi` as a new platform-api dependency — genuine 3.1 support, but a new module to vet under `dependency-management.md`.
> - **(c)** Restrict this feature to 3.0.x outright and reject 3.1 with a clear message.
>
> Recommendation: **(a)** — it keeps FR-003's promise honest for the common case without a new dependency, and 3.1 documents remain publishable to the portal.

**No-fetch requirement (FR-011)**: kin-openapi's `openapi3.Loader` resolves external `$ref`s over the network by default. Both of these must be set on every loader instance:
```go
loader := openapi3.NewLoader()
loader.IsExternalRefsAllowed = false
loader.ReadFromURIFunc = func(*openapi3.Loader, *url.URL) ([]byte, error) {
    return nil, errors.New("external references are not resolved")
}
```
This must be asserted by a unit test — the same "assert the hardening flags in a test so an upgrade can't silently re-enable them" requirement `xxe-xml-processing.md` Directive 1 applies to DTD-aware parsers.

---

## R3. Request body size bounding

**Decision**: Apply `http.MaxBytesReader` in the upload handler, sized from a new config key `openapi_spec_max_bytes` (default 5 MiB, matching the existing `openapi_spec_max_fetch_bytes = 5242880` in `platform-api/config/config-template.toml:66`).

**Rationale**: The main HTTP server (`platform-api/internal/server/server.go:881`) sets `ReadTimeout`/`ReadHeaderTimeout` but installs **no** global `MaxBytesReader`. The existing `max_body_size` key is scoped to the webhook receiver only (`internal/webhook/receiver.go:150`), not the management API. So without an explicit cap in this handler, FR-004 is unmet and `go-network-service-hardening.md` Directive 1 is violated.

**Alternatives considered**: A global body cap on the whole management API — correct in principle but a much wider blast radius (it would newly bound every existing endpoint) and out of this feature's scope. Worth raising separately.

---

## R4. Endpoint shape, scopes, and enforcement

**Decision**: A single `openapi` sub-resource under the existing REST API resource:

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/rest-apis/{restApiId}/openapi` | Retrieve (FR-005) |
| `PUT` | `/rest-apis/{restApiId}/openapi` | Create-or-replace (FR-001, FR-006) |
| `DELETE` | `/rest-apis/{restApiId}/openapi` | Remove (FR-007) |
| `GET` | `/rest-apis/{restApiId}/openapi/drift` | Drift report (FR-014) |

**Rationale**: Exactly one current specification exists per API (spec Assumptions), so upload and update are the same idempotent operation — `PUT` expresses that without a POST/PUT split that would need conflict semantics. This mirrors the existing sub-resource pattern (`/rest-apis/{restApiId}/gateways`, `/rest-apis/{restApiId}/api-keys`) in `platform-api/resources/openapi.yaml`.

**Scopes** follow the existing sub-resource convention seen at `resources/openapi.yaml:532` (`ap:rest_api:gateway:read` / `:manage`, with parent `ap:rest_api:manage` always accepted):

| Operation | Accepted scopes |
|---|---|
| `GET .../openapi`, `GET .../openapi/drift` | `ap:rest_api:openapi:read`, `ap:rest_api:openapi:manage`, `ap:rest_api:manage` |
| `PUT .../openapi` | `ap:rest_api:openapi:update`, `ap:rest_api:openapi:manage`, `ap:rest_api:manage` |
| `DELETE .../openapi` | `ap:rest_api:openapi:delete`, `ap:rest_api:openapi:manage`, `ap:rest_api:manage` |

Separate `:read` and `:update` satisfies FR-010's "read access can be granted without write access". No wildcard forms — `ScopeEnforcer` matches by exact string equality (GO-AUTH-020).

**Enforcement is spec-driven, not hand-written**: `internal/middleware/openapi_scope_registry.go` builds the (method, path) → scopes map by parsing `resources/openapi.yaml` at startup. Declaring the `security` block in the spec *is* the enforcement. `internal/server/scope_route_coverage_test.go` already asserts every registered route has a scope entry — a new route without a spec `security` block fails that existing test.

---

## R5. Content negotiation (JSON and YAML)

**Decision**: Store the document **as uploaded** (original bytes plus a `format` discriminator); convert on read only when the caller asks for the other representation. Request format from `Content-Type` (`application/json`, `application/yaml`/`application/x-yaml`); response format from `Accept`, defaulting to the stored format.

**Rationale**: FR-002 requires the two forms to be treated equivalently, and SC-002 requires a round trip to preserve operations, paths, schemas and descriptions. Storing normalized-to-JSON would silently drop YAML comments and key ordering, which producers notice. Storing original bytes makes the round trip exact for the uploaded form and semantically equivalent for the converted form — which is what SC-002 actually asks for.

---

## R6. Keeping the specification out of the runtime artifact (FR-013 / SC-007)

**Decision**: Assert the invariant at `platform-api/internal/utils` `BuildAPIDeploymentYAML`, which is the single function that renders an API into its gateway runtime artifact (called from `ensureRESTRuntimeArtifactUnchanged`, `internal/service/api.go:419`).

**Rationale**: FR-013 says no specification operation may change routing. The provable form of that statement is: the specification is not an input to `BuildAPIDeploymentYAML`. A test that builds the runtime artifact before and after a spec upload/update/delete and asserts byte equality gives SC-007 a direct, non-hand-wavy check. This reuses the codebase's own existing notion of "the runtime artifact" rather than inventing a parallel one.

---

## R7. Drift detection (FR-014)

**Decision**: Compute drift on read, not on write. Compare the set of `(METHOD, path)` pairs from the stored specification's `paths` against `model.RestAPIConfig.Operations[].Request{Method,Path}` (`internal/model/api.go:71`), and report the symmetric difference in both directions.

**Rationale**: Computing on write would need re-computation whenever the API's operations change through the *existing* update path, silently staling the stored answer. On-read is always correct and costs one parse of an already-size-bounded document.

**Normalization (FR-019)**: Methods uppercased via `strings.ToUpper` at extraction (GO-AUTH-006), paths compared after collapsing duplicate separators and resolving `.`/`..` segments. OpenAPI path templating means `/pets/{id}` and `/pets/{petId}` are the same route — compare with parameter names normalized away, exactly as `normalizePathParams` (`internal/middleware/openapi_scope_registry.go:41`) already does for the scope registry. Reuse that helper rather than writing a second one.

---

## R8. Developer portal integration (FR-016)

**Finding**: The API Portal **already** has both halves of what this feature feeds:
- `POST /apis` (`portals/api-portal/docs/api-portal-openapi-spec-v0.9.yaml:286`) accepts an API definition file at creation time — "An API definition file is required unless supplied by the artifact ZIP".
- `api_contents` (`portals/api-portal/database/schema.sqlite.sql:149`) stores spec files, docs and icons keyed by `(api_uuid, type, file_name)`.

So the portal side needs **no schema or endpoint change**. The gap this feature closes is purely on the Platform API side: today there is nowhere to hold the definition that the portal's publish call wants.

> **Open question — publish orchestration**: what actually drives the push into the portal (the `ap` CLI's devportal client at `cli/src/internal/devportal/client.go`, an operator, or a Platform API-side flow) was not determined from the code in this pass. FR-016 says only that the stored specification MUST be *available* to that flow, which the `GET .../openapi` endpoint satisfies regardless of who calls it. Wiring the automatic push — and FR-017's "portal reflects an updated specification without republishing" — depends on the answer and should be confirmed before task breakdown.

---

## R9. Model generation and CLI propagation

**Findings**:
- Request/response models are generated by oapi-codegen from `resources/openapi.yaml` into `api/generated.go` (`platform-api/oapi-codegen.yaml`, `models: true`, `skip-prune: true`). Editing the spec then regenerating is the workflow; hand-writing DTOs would diverge.
- The `sync-cli-with-openapi` skill covers propagating a `platform-api/resources/openapi.yaml` change into the `ap` CLI. This feature adds endpoints, so that skill applies at implementation time.

---

## R10. Constitution status

`.specify/memory/constitution.md` is still the **unfilled Spec Kit template** — every principle is a `[PRINCIPLE_N_NAME]` placeholder. There are therefore no ratified project principles to gate against, and no gate outcome can honestly be reported from it.

**Substitute gates used**: the repo's `.claude/rules/*.md`, which are checked-in project instructions with MUST force and directly govern this feature's surface — `db-schema-changes.md`, `authentication_authorization.md`, `error-handling.md`, `file-access.md`, `go-network-service-hardening.md`, `ssrf-prevention.md`, `xxe-xml-processing.md`, `dependency-management.md`. These are evaluated in `plan.md`'s Constitution Check. Running `/speckit.constitution` to fill the real constitution is recommended but is not a blocker for this feature.
