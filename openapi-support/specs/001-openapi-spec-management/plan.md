# Implementation Plan: OpenAPI Specification Management for REST APIs

**Branch**: `001-openapi-spec-management` | **Date**: 2026-09-01 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/001-openapi-spec-management/spec.md`

## Summary

Platform API can today describe a REST API only as a list of individually-declared resources. This feature adds a single current OpenAPI document per REST API — uploaded, retrieved, replaced and removed through an `openapi` sub-resource on the existing REST API resource — so the producer's authored contract lives with the API and can feed developer-portal publishing.

The specification is stored in a **new** `rest_api_specs` table (one row per API), deliberately kept out of `rest_apis.configuration` so it is never an input to the gateway runtime artifact. Where the specification's operations and the API's declared resources disagree, the platform reports the drift on read; it never synchronises them. Validation uses `getkin/kin-openapi`, already pinned in-repo, with external reference resolution hard-disabled so no upload can cause an outbound fetch.

## Technical Context

**Language/Version**: Go 1.26.5 (`platform-api/go.mod`)

**Primary Dependencies**: `github.com/getkin/kin-openapi v0.133.0` (new to `platform-api`; already pinned at this version in `gateway/gateway-controller` and `event-gateway/gateway-controller`). Existing: `oapi-codegen` for model generation, `gopkg.in/yaml.v3`, `sqlx`/`database/sql`.

**Storage**: The platform's existing relational store, across all three supported dialects — SQLite, PostgreSQL, SQL Server (`platform-api/internal/database/schema.*.sql`).

**Testing**: `go test` — unit tests beside each package, plus the existing integration-test pattern (`*_integration_test.go` in `internal/handler` and `internal/service`) and the route/scope coverage test at `internal/server/scope_route_coverage_test.go`.

**Target Platform**: Linux server (containerized), same deployment unit as the rest of Platform API.

**Project Type**: Backend service — additive change to an existing Go web service. No new deployable.

**Performance Goals**: SC-003 — 95% of retrievals under 1s for documents up to the configured maximum. Retrieval is a single indexed row read plus, only when the caller asks for the other representation, one format conversion.

**Constraints**: Upload bodies bounded by a new `openapi_spec_max_bytes` config key (default 5 MiB) — the management API installs no global body cap today (research R3). No outbound request may originate from document content (FR-011). No specification operation may alter routing (FR-013).

**Scale/Scope**: 4 new endpoints, 1 new table, 1 new config key, 3 new scopes, 2 additive response fields on the existing `RESTAPI` schema. Documents typically 10 KB–2 MB.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Constitution status**: `.specify/memory/constitution.md` is the **unfilled Spec Kit template** — every principle is still a `[PRINCIPLE_N_NAME]` placeholder. No ratified principles exist, so no gate result can honestly be reported from it. Running `/speckit.constitution` is recommended, but is not a blocker for this feature.

**Substitute gates**: the repo's checked-in `.claude/rules/*.md`, which carry MUST force and govern this feature's surface directly.

| Gate | Rule | Status after Phase 1 design |
|---|---|---|
| **G1 — Shipped tables stay frozen** | `db-schema-changes.md` D1/D3 | **PASS.** New table only; `rest_apis` untouched. A new table needs only its guarded `CREATE`, so no per-dialect `ALTER` path is owed. |
| **G2 — Specification never reaches routing** | FR-013, SC-007 | **PASS by design.** Spec lives outside `rest_apis.configuration`; the two new `RESTAPI` fields are `readOnly`. Enforced by a test asserting `BuildAPIDeploymentYAML` output is byte-identical across spec upload/update/delete (research R6). |
| **G3 — Deny-by-default authorization** | GO-AUTH-005, GO-AUTH-007 | **PASS.** Every operation carries an explicit `security` block with named scopes; org scoping comes from verified JWT claims and is also a column on the new table. Read and write scopes are separate (FR-010). |
| **G4 — Scope enforcement is structural** | GO-AUTH-017 | **PASS.** Scopes derive from `resources/openapi.yaml` via `ScopeRegistry`; `scope_route_coverage_test.go` already fails any route lacking a spec entry. No hand-written path matching. |
| **G5 — Bounded request bodies** | `go-network-service-hardening.md` D1, `file-access.md` D5 | **PASS.** `http.MaxBytesReader` sized from config, applied before any read; generic `413` that does not echo the limit. |
| **G6 — No document-driven outbound requests** | `ssrf-prevention.md`, FR-011 | **PASS.** Loader configured with `IsExternalRefsAllowed = false` and a `ReadFromURIFunc` that always errors, asserted in a unit test. Import-by-URL is out of scope (FR-015), so no fetch path is introduced at all. |
| **G7 — Untrusted structured-document parsing** | `xxe-xml-processing.md` (pattern), `file-access.md` D6 | **PASS.** Size ceiling before parse; parse bounded by a context deadline; documents are only ever stored and returned, never compiled or evaluated. Hardening flags asserted in tests so a dependency upgrade cannot silently re-enable them. |
| **G8 — Sterile error responses** | `error-handling.md` D1/D4 | **PASS.** Validation failures return a generic `422` with no parser internals, file paths, or stack detail; the specific reason is logged internally only. |
| **G9 — Dependency addition** | `dependency-management.md` | **PASS with a task.** `kin-openapi v0.133.0` is already in-repo at this exact version, but adding it to `platform-api/go.mod` still owes `govulncheck ./...` and `go-licenses check ./...` runs before the PR. |
| **G10 — No deferral behind comments** | GO-AUTH-019 and the equivalent directive in every rule above | **PASS.** No gap in this plan is parked behind a `// TODO`. The two genuinely open items (below) are surfaced as decisions, not annotations. |

**Open items carried into task breakdown** — neither blocks Phase 2, both need an answer before the affected tasks are written:
1. **OpenAPI 3.1 depth** (research R2). kin-openapi targets 3.0.x; the spec's Assumptions claim 3.0 *and* 3.1. Recommendation: accept 3.1 with lighter structural validation and say so, rather than adding a second parser or rejecting 3.1.
2. **Publish orchestration** (research R8). The portal already accepts a definition file and already stores it in `api_contents` — no portal change is needed. What was not determined is *who* drives the push (the `ap` CLI's devportal client, an operator, or a Platform API flow), which is what FR-017's "reflects an update without republishing" depends on.

## Project Structure

### Documentation (this feature)

```text
specs/001-openapi-spec-management/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── rest-api-openapi.yaml   # Phase 1 output — merge target: platform-api/resources/openapi.yaml
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 output (/speckit.tasks — NOT created by /speckit.plan)
```

### Source Code (repository root)

All changes land in the existing `platform-api` service. No new module, no new deployable, no portal-side change.

```text
platform-api/
├── resources/
│   ├── openapi.yaml                      # MODIFY — merge contracts/rest-api-openapi.yaml; source of truth for
│   │                                     #          both generated models and scope enforcement
│   └── role-to-scope-mapping.yaml        # MODIFY — grant the three new scopes to the appropriate roles
├── api/
│   └── generated.go                      # REGENERATE via oapi-codegen (do not hand-edit)
├── config/
│   ├── config.go                         # MODIFY — add openapi_spec_max_bytes
│   └── config-template.toml              # MODIFY — document it, default 5 MiB
├── internal/
│   ├── database/
│   │   ├── schema.sqlite.sql             # MODIFY — CREATE TABLE IF NOT EXISTS rest_api_specs + index
│   │   ├── schema.postgres.sql           # MODIFY — same, postgres form
│   │   └── schema.sqlserver.sql          # MODIFY — same, OBJECT_ID-guarded form
│   ├── model/
│   │   └── api_spec.go                   # NEW — RestAPISpec, SpecDrift
│   ├── repository/
│   │   └── api_spec.go                   # NEW — upsert / get / get-metadata / delete, all org-scoped
│   ├── service/
│   │   ├── api_spec.go                   # NEW — validation, org authorization, drift computation
│   │   └── api.go                        # MODIFY — populate hasOpenApiSpec / openApiSpecUpdatedAt on read
│   ├── handler/
│   │   └── api_spec.go                   # NEW — 4 handlers + RegisterRoutes, body cap, content negotiation
│   └── utils/
│       └── openapi_validator.go          # NEW — hardened loader (no external refs), parse deadline
└── go.mod / go.sum                       # MODIFY — add getkin/kin-openapi v0.133.0
```

**Structure Decision**: Follow the service's existing four-layer split — `handler` (HTTP concerns: body cap, content negotiation, status mapping) → `service` (validation, authorization, drift) → `repository` (SQL) → `model`. New files are named `api_spec.go` in each layer, matching the `api.go` / `api_key.go` / `api_deployment.go` naming already in use. Handlers are registered through the existing `RegisterRoutes(mux router.Router)` convention so the route/scope coverage test picks them up automatically.

Per GO-AUTH-015, the org-scoping and existence checks live in the **service** layer, not only in the handler, so any future caller (event handler, internal RPC) passes through them.

## Complexity Tracking

No Constitution Check gate is violated, so no justification is owed. Two design choices worth recording as deliberate, since each could look like extra machinery:

| Choice | Why | Simpler alternative rejected because |
|---|---|---|
| Separate `rest_api_specs` table rather than a column on `rest_apis` | Keeps a multi-hundred-KB blob off every API read and list, and makes G2 structural rather than a convention | A column would need a three-dialect `ALTER TABLE` upgrade path and would load spec bytes into every list query |
| Drift as its own endpoint, computed on read | The `GET .../openapi` response is the raw document under content negotiation, so it has no envelope to carry a drift object | Computing drift on write would go stale whenever the API's operations change through the existing update path |
