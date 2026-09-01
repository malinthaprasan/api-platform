---
description: "Task list for OpenAPI Specification Management for REST APIs"
---

# Tasks: OpenAPI Specification Management for REST APIs

**Input**: Design documents from `/specs/001-openapi-spec-management/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/rest-api-openapi.yaml, quickstart.md

**Tests**: Test tasks ARE included. They are not optional here — the specification itself requires them (SC-006 "verified by test coverage over every specification operation", SC-007 "verified by comparing the API's resources and any active deployment before and after each operation"), and three project rules mandate specific assertions in tests: `xxe-xml-processing.md` D1 (parser hardening flags asserted so an upgrade cannot re-enable them), `go-cors-validation.md`-style regression coverage, and `ssrf-prevention.md` D6 (a rejection test per input kind).

**Organization**: Grouped by user story so each can be implemented, tested and demoed independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependency on an incomplete task)
- **[Story]**: US1–US4, mapping to the user stories in spec.md
- All paths are repo-relative

## Path Conventions

Single Go service — all changes land in `platform-api/`, per plan.md's Structure Decision. Four layers: `internal/handler` → `internal/service` → `internal/repository` → `internal/model`, with shared helpers in `internal/utils`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Resolve blocking decisions and get the dependency and configuration in place

- [ ] T001 Resolve open decision: OpenAPI 3.1 validation depth (research.md R2) — choose (a) accept 3.1 with lighter structural validation, (b) add `pb33f/libopenapi`, or (c) reject 3.1; record the answer in `specs/001-openapi-spec-management/research.md` and update the Assumptions section of `specs/001-openapi-spec-management/spec.md` to match. BLOCKS T012, T013.
- [ ] T002 Resolve open decision: developer-portal publish orchestration (research.md R8) — determine whether the `ap` CLI devportal client (`cli/src/internal/devportal/client.go`), an operator, or a Platform API flow drives the push, and record it in `specs/001-openapi-spec-management/research.md`. BLOCKS T041, T042.
- [ ] T003 Add `github.com/getkin/kin-openapi v0.133.0` to `platform-api/go.mod` (matching the pin in `gateway/gateway-controller/go.mod:8`) and run `go mod tidy` in `platform-api/`
- [ ] T004 Run the dependency gates owed by T003 from `platform-api/`: `govulncheck ./...` and `go-licenses check ./... --disallowed_types=forbidden,restricted`; record results in the PR description per `dependency-management.md`
- [ ] T005 [P] Add `openapi_spec_max_bytes` (default 5242880) to the config struct in `platform-api/config/config.go`, following the existing `OpenAPISpecMaxFetchBytes` field at line 101
- [ ] T006 [P] Document `openapi_spec_max_bytes` in `platform-api/config/config-template.toml` beside the existing `openapi_spec_max_fetch_bytes` at line 66

**Checkpoint**: Dependency vetted, config key available, both blocking decisions answered

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Storage, contract, scopes and the hardened parser — shared by every user story

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T007 Run the `designing-db-schemas` skill against the `rest_api_specs` shape proposed in `specs/001-openapi-spec-management/data-model.md` to apply rules R1–R10; reconcile any differences into data-model.md before writing DDL. BLOCKS T008, T009, T010.
- [ ] T008 Add `CREATE TABLE IF NOT EXISTS rest_api_specs` plus `idx_rest_api_specs_org` to `platform-api/internal/database/schema.sqlite.sql`, per data-model.md
- [ ] T009 [P] Add the PostgreSQL form of the same table and index to `platform-api/internal/database/schema.postgres.sql`
- [ ] T010 [P] Add the SQL Server form to `platform-api/internal/database/schema.sqlserver.sql`, using the `OBJECT_ID`-guarded pattern already used in that file (not `IF NOT EXISTS`)
- [ ] T011 [P] Create `RestAPISpec` and `SpecDrift` structs in `platform-api/internal/model/api_spec.go`, per the data-model.md column and derived-value tables
- [ ] T012 Create the hardened OpenAPI loader in `platform-api/internal/utils/openapi_validator.go` — `openapi3.NewLoader()` with `IsExternalRefsAllowed = false` and a `ReadFromURIFunc` that always returns an error, plus a `context.WithTimeout` parse deadline (research.md R2, FR-011)
- [ ] T013 Add unit tests in `platform-api/internal/utils/openapi_validator_test.go` asserting that a document with an external `$ref` is never fetched and that both hardening settings are in force — so a kin-openapi upgrade cannot silently re-enable them (`xxe-xml-processing.md` D1)
- [ ] T014 Implement `RestAPISpecRepo` in `platform-api/internal/repository/api_spec.go` with `Upsert`, `Get`, `GetMetadata`, `Delete` — every method takes and filters on `organization_uuid`, all queries parameterized (GO-AUTH-005, GO-AUTH-008)
- [ ] T015 [P] Add repository unit tests in `platform-api/internal/repository/api_spec_test.go` covering upsert-replaces, delete, and that a query for another organization's row returns nothing
- [ ] T016 Merge the paths and schemas from `specs/001-openapi-spec-management/contracts/rest-api-openapi.yaml` into `platform-api/resources/openapi.yaml`, including the two `readOnly` additions to the existing `RESTAPI` schema
- [ ] T017 Regenerate `platform-api/api/generated.go` via oapi-codegen using `platform-api/oapi-codegen.yaml` (do not hand-edit the generated file)
- [ ] T018 Grant `ap:rest_api:openapi:read`, `ap:rest_api:openapi:update` and `ap:rest_api:openapi:delete` to the appropriate roles in `platform-api/resources/role-to-scope-mapping.yaml`, following the existing `ap:rest_api:*` entries
- [ ] T019 Create `RestAPISpecService` in `platform-api/internal/service/api_spec.go` with the shared org-scope-and-existence predicate every entry point funnels through (GO-AUTH-015) — service layer, not handler-only

**Checkpoint**: Table exists on all three dialects, contract and scopes are declared, parser is hardened, storage and service skeletons are in place

---

## Phase 3: User Story 1 - Attach an OpenAPI specification (Priority: P1) 🎯 MVP

**Goal**: A producer can upload a JSON or YAML OpenAPI document against an existing REST API, with invalid or oversized documents rejected and nothing stored.

**Independent Test**: Create a REST API, `PUT` a valid document, confirm success; `PUT` garbage and an oversized file, confirm both are rejected and nothing is stored.

### Tests for User Story 1

- [ ] T020 [P] [US1] Integration test in `platform-api/internal/handler/api_spec_integration_test.go` covering acceptance scenarios 1–2: valid JSON upload and valid YAML upload both stored (FR-001, FR-002)
- [ ] T021 [P] [US1] Integration test in `platform-api/internal/handler/api_spec_integration_test.go` for scenario 3: a non-OpenAPI document returns `422` with no parser internals in the body, and no row is written (FR-003, `error-handling.md` D1)
- [ ] T022 [P] [US1] Integration test in `platform-api/internal/handler/api_spec_integration_test.go` for scenario 4: a body above `openapi_spec_max_bytes` returns `413` without echoing the configured limit, and no partial content is stored (FR-004)
- [ ] T023 [P] [US1] Integration test in `platform-api/internal/handler/api_spec_integration_test.go` for scenario 5: upload with a token from another organization is refused and nothing changes (FR-009, SC-006)

### Implementation for User Story 1

- [ ] T024 [US1] Implement `PutSpec` in `platform-api/internal/service/api_spec.go` — validate fully before any write, so a failure leaves any existing row untouched (FR-006); return sterile, typed errors via `internal/apperror`
- [ ] T025 [US1] Implement format detection in `platform-api/internal/service/api_spec.go` from `Content-Type` (`application/json`, `application/yaml`, `application/x-yaml`) and store the original bytes plus the `format` discriminator (research.md R5, FR-002)
- [ ] T026 [US1] Implement the `PUT /rest-apis/{restApiId}/openapi` handler in `platform-api/internal/handler/api_spec.go`, applying `http.MaxBytesReader` sized from `openapi_spec_max_bytes` **before** any body read (`go-network-service-hardening.md` D1, FR-004)
- [ ] T027 [US1] Map service errors to responses in `platform-api/internal/handler/api_spec.go`: `413` oversized, `422` invalid document, `404` unknown API, `403` wrong organization — all generic, with the specific reason logged internally only (`error-handling.md` D1/D4)
- [ ] T028 [US1] Add `RegisterRoutes(mux router.Router)` to `platform-api/internal/handler/api_spec.go` and wire the handler into `platform-api/internal/server/server.go` alongside the other `RegisterRoutes` calls
- [ ] T029 [US1] Add the new handler to `platform-api/internal/server/scope_route_coverage_test.go` so the existing gate verifies every new route has a scope declared in the spec (GO-AUTH-017)

**Checkpoint**: A producer can attach a specification to an API; invalid and oversized documents are rejected cleanly. MVP is demoable.

---

## Phase 4: User Story 2 - Retrieve the specification (Priority: P1)

**Goal**: The stored document can be read back in either representation, and an API's spec presence is visible without fetching the document.

**Independent Test**: Upload a specification, retrieve it as JSON and as YAML, and confirm both are semantically equivalent to the upload.

### Tests for User Story 2

- [ ] T030 [P] [US2] Round-trip test in `platform-api/internal/handler/api_spec_integration_test.go`: `PUT` YAML then `GET` as both JSON and YAML; assert operations, paths, schemas and descriptions all survive (FR-005, SC-002)
- [ ] T031 [P] [US2] Test in `platform-api/internal/handler/api_spec_integration_test.go` that an API with no specification, and an unknown API id, produce two distinguishable `404` responses (FR-008)
- [ ] T032 [P] [US2] Test in `platform-api/internal/handler/api_spec_integration_test.go` that a cross-organization `GET` is refused without disclosing whether a specification exists (FR-009)

### Implementation for User Story 2

- [ ] T033 [US2] Implement `GetSpec` in `platform-api/internal/service/api_spec.go`, returning the stored bytes plus format and last-changed metadata (FR-005, FR-012)
- [ ] T034 [US2] Implement the `GET /rest-apis/{restApiId}/openapi` handler in `platform-api/internal/handler/api_spec.go` with `Accept`-based content negotiation, defaulting to the stored format, and a `Last-Modified` header (research.md R5)
- [ ] T035 [US2] Populate the `hasOpenApiSpec` and `openApiSpecUpdatedAt` response fields in `platform-api/internal/service/api.go` on both get and list, sourcing them from `GetMetadata` so the document body is never loaded (FR-020, FR-012)
- [ ] T036 [P] [US2] Test in `platform-api/internal/service/api_spec_test.go` that the two new fields are response-only — rejected if supplied on create or update, and absent from `model.RestAPIConfig`

**Checkpoint**: Upload and retrieval both work; the portal and tooling now have a contract to consume.

---

## Phase 5: User Story 3 - Update, remove, and see drift (Priority: P2)

**Goal**: A producer can replace or remove the specification, with a failed replacement leaving the previous document intact, and can see where the specification and the API's declared resources disagree.

**Independent Test**: Upload a spec, upload a revised one, confirm the revision is served; upload an invalid one, confirm the prior spec survives; remove it, confirm it is gone and the API is unaffected.

### Tests for User Story 3

- [ ] T037 [P] [US3] Test in `platform-api/internal/handler/api_spec_integration_test.go` covering scenarios 1–2: a valid revision replaces the previous document; an invalid revision is rejected and the previous document is still retrievable (FR-006)
- [ ] T038 [P] [US3] Test in `platform-api/internal/handler/api_spec_integration_test.go` for scenario 3: `DELETE` removes the specification and leaves the API itself intact (FR-007)
- [ ] T039 [P] [US3] Test in `platform-api/internal/handler/api_spec_integration_test.go` for scenario 4: a specification declaring a path the API's resources do not (and vice versa) uploads successfully, leaves resources untouched, and is reported as drift (FR-014)
- [ ] T040 [P] [US3] Test in `platform-api/internal/service/api_spec_drift_test.go` that drift is computed identically for lowercase methods, duplicate path separators, and differing path-parameter names (`/pets/{id}` vs `/pets/{petId}`) (FR-019)
- [ ] T041 [US3] Test in `platform-api/internal/service/api_spec_runtime_test.go` that `BuildAPIDeploymentYAML` output is byte-identical before and after upload, update and delete, and that the API's declared operations are unchanged (FR-013, SC-007, research.md R6)

### Implementation for User Story 3

- [ ] T042 [US3] Implement `DeleteSpec` in `platform-api/internal/service/api_spec.go` and the `DELETE /rest-apis/{restApiId}/openapi` handler in `platform-api/internal/handler/api_spec.go`, returning `204` (FR-007)
- [ ] T043 [US3] Implement drift computation in `platform-api/internal/service/api_spec.go` — compare the specification's `(METHOD, path)` set against `model.RestAPIConfig.Operations[].Request`, reporting `inSpecOnly` and `inResourcesOnly` (FR-014, research.md R7)
- [ ] T044 [US3] Reuse `normalizePathParams` from `platform-api/internal/middleware/openapi_scope_registry.go:41` for path comparison and `strings.ToUpper` for methods — do not write a second normalizer (FR-019, GO-AUTH-006)
- [ ] T045 [US3] Implement the `GET /rest-apis/{restApiId}/openapi/drift` handler in `platform-api/internal/handler/api_spec.go` returning the `RESTAPIOpenAPISpecDrift` schema

**Checkpoint**: The full lifecycle works, drift is visible, and routing is provably untouched.

---

## Phase 6: User Story 4 - Publish to the developer portal (Priority: P2)

**Goal**: A published API's portal documentation is derived from its stored specification.

**Independent Test**: Upload a specification, publish the API to the portal, and confirm the portal renders operations, descriptions and schemas from that document.

> **Note**: research.md R8 established that the portal needs **no** change — `POST /apis` already accepts a definition file and `api_contents` already stores it. T046–T047 depend on the T002 decision about who drives the push; T048 does not.

- [ ] T046 [US4] BLOCKED ON T002 — wire the publish path identified in T002 to read the specification from `GET /rest-apis/{restApiId}/openapi` and supply it as the portal's API definition file (FR-016)
- [ ] T047 [US4] BLOCKED ON T002 — ensure a specification update propagates to an already-published portal listing without republishing the API from scratch (FR-017, acceptance scenario 3)
- [ ] T048 [P] [US4] Confirm and test in `platform-api/internal/handler/api_spec_integration_test.go` that publishing an API with no specification still succeeds, presenting only its declared resources as it does today (acceptance scenario 2)

**Checkpoint**: The stated business driver is delivered — portal documentation comes from the producer's own contract.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T049 [P] Test in `platform-api/internal/repository/api_spec_test.go` that concurrent `PUT`s to the same API resolve to exactly one complete document, with no partially-written row ever readable (FR-018)
- [ ] T050 [P] Run the `sync-cli-with-openapi` skill so the `ap` CLI gains commands for the four new endpoints, per research.md R9
- [ ] T051 [P] Add the four new endpoints to the Platform API docs under `docs/`, alongside the existing REST API operations
- [ ] T052 Run every check in the `specs/001-openapi-spec-management/quickstart.md` verification table against a running instance
- [ ] T053 Run `go test ./...` in `platform-api/` and confirm `scope_route_coverage_test.go` passes with the new routes
- [ ] T054 Confirm no `// TODO`/`FIXME` was used to defer any rule finding in the delivered code (GO-AUTH-019 and the equivalent directive in every rule cited by plan.md)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies. T001 and T002 are decision tasks — start them first, since T001 blocks the validator and T002 blocks two US4 tasks.
- **Foundational (Phase 2)**: Depends on Setup. T007 (schema skill) blocks all DDL. BLOCKS every user story.
- **US1 (Phase 3)**: Depends on Foundational.
- **US2 (Phase 4)**: Depends on Foundational. Independently testable, but only demoable end-to-end once US1 can put a document in place.
- **US3 (Phase 5)**: Depends on Foundational. T037's "replace" scenarios exercise the US1 handler.
- **US4 (Phase 6)**: Depends on Foundational and on T002; T046–T047 also need US2's `GET` endpoint as their read source.
- **Polish (Phase 7)**: After all desired stories.

### Within Each User Story

Tests → service → handler → route registration. Service before handler throughout, since org-scoping and existence checks live in the service layer (GO-AUTH-015).

### Parallel Opportunities

- **Phase 1**: T005 and T006 in parallel; T001 and T002 are independent decisions that can be pursued at once.
- **Phase 2**: T009 and T010 in parallel once T008 establishes the shape; T011 and T015 in parallel with the DDL work.
- **Phase 3–5**: All test tasks within a story are `[P]` — different scenarios, and after T020 creates it, the same test file can take parallel additions if contributors coordinate on placement.
- **Across stories**: US1, US2 and US3 can be staffed in parallel once Phase 2 completes, though they share `internal/handler/api_spec.go` and `internal/service/api_spec.go` — split by function, not by file, to avoid conflicts.

---

## Parallel Example: User Story 1

```bash
# Launch all US1 tests together:
Task: "Valid JSON and YAML upload test in platform-api/internal/handler/api_spec_integration_test.go"
Task: "Invalid-document 422 test in platform-api/internal/handler/api_spec_integration_test.go"
Task: "Oversized-body 413 test in platform-api/internal/handler/api_spec_integration_test.go"
Task: "Cross-organization refusal test in platform-api/internal/handler/api_spec_integration_test.go"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

Both are P1, and neither is useful alone — storing a document nobody can read delivers nothing, and retrieval needs something to retrieve. Treat US1 + US2 together as the MVP:

1. Phase 1: Setup (answer T001 first — it blocks the validator)
2. Phase 2: Foundational
3. Phase 3: US1 → **validate**: a document can be attached and bad ones are rejected
4. Phase 4: US2 → **validate**: the document round-trips in both representations
5. Demo: a producer's real contract is stored in, and served from, Platform API

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. US1 + US2 → MVP, demoable
3. US3 → full lifecycle plus drift visibility; SC-007 proven by T041
4. US4 → the portal publishing payoff (needs T002 answered)
5. Polish → CLI, docs, quickstart validation

---

## Notes

- **US1's `PUT` is inherently create-or-replace**, so US3 does not add a new endpoint for update — it adds the *guarantee* (T037: a failed replacement preserves the prior document) plus `DELETE` and drift. This is why the two stories share a handler file.
- `platform-api/api/generated.go` is generated — edit `resources/openapi.yaml` and regenerate (T017), never hand-edit.
- `resources/openapi.yaml` is also the authorization config: the `security` block *is* the enforcement, via `ScopeRegistry`.
- Commit after each task or logical group; stop at any checkpoint to validate a story independently.
