# Phase 1 Data Model: OpenAPI Specification Management

**Feature**: `001-openapi-spec-management` | **Date**: 2026-09-01

## Entity: REST API OpenAPI Specification

One current specification per REST API. New table — under `db-schema-changes.md` Directive 3 a new table needs only its guarded `CREATE TABLE IF NOT EXISTS`, which covers fresh and already-provisioned databases alike, so **no per-dialect `ALTER TABLE` is required**.

### Table `rest_api_specs`

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `uuid` | `VARCHAR(40)` | PRIMARY KEY | Row identity, matching every sibling table's PK shape |
| `rest_api_uuid` | `VARCHAR(40)` | NOT NULL, FK → `rest_apis(uuid)` ON DELETE CASCADE, UNIQUE | One current spec per API; cascade satisfies "API deleted ⇒ spec gone" |
| `organization_uuid` | `VARCHAR(40)` | NOT NULL, FK → `organizations(uuid)` ON DELETE CASCADE | Org scoping at the data layer (GO-AUTH-005), not only in the handler |
| `format` | `VARCHAR(10)` | NOT NULL | `JSON` or `YAML` — the representation as uploaded (R5) |
| `openapi_version` | `VARCHAR(10)` | NOT NULL | e.g. `3.0.3`, `3.1.0`, read from the document's `openapi` field |
| `content` | `BLOB` | NOT NULL | Original uploaded bytes, unmodified |
| `content_size` | `INTEGER` | NOT NULL | Byte length, so listing/metadata reads never load `content` |
| `created_by` | `VARCHAR(200)` | | Audit, matching `rest_apis` |
| `created_at` | `DATETIME` | DEFAULT CURRENT_TIMESTAMP | |
| `updated_by` | `VARCHAR(200)` | | FR-012 "by whom" |
| `updated_at` | `DATETIME` | DEFAULT CURRENT_TIMESTAMP | FR-012 "when last changed" |

### Index

```sql
CREATE INDEX IF NOT EXISTS idx_rest_api_specs_org ON rest_api_specs(organization_uuid);
```
`rest_api_uuid` needs no separate index — its UNIQUE constraint provides one.

### Dialect files to update

All three, kept aligned (`db-schema-changes.md` multi-engine alignment):
- `platform-api/internal/database/schema.sqlite.sql`
- `platform-api/internal/database/schema.postgres.sql`
- `platform-api/internal/database/schema.sqlserver.sql`

SQL Server needs the `OBJECT_ID`-guarded form rather than `IF NOT EXISTS`, per the convention already used in that file.

> **Required before implementation**: run the `designing-db-schemas` skill against this table to apply rules R1–R10 (column types and widths, PK/FK shape, org-scoping, audit columns, indexing, naming). The shape above follows the conventions observed in `rest_apis` and its siblings, but R1–R10 have **not** been formally applied — that is a task, not a completed step.

---

## Validation rules

| Rule | Source | Enforced at |
|---|---|---|
| Body ≤ `openapi_spec_max_bytes` (default 5 MiB) | FR-004 | Handler, via `http.MaxBytesReader` before any read |
| Parses as JSON or YAML per `Content-Type` | FR-002 | Service |
| Valid OpenAPI document; `openapi` field present; `paths` present | FR-003 | Service, via kin-openapi loader |
| No external `$ref` resolution, ever | FR-011 | Loader configured with `IsExternalRefsAllowed = false` and a failing `ReadFromURIFunc` (R2) |
| Target API exists, is in the caller's org | FR-009, FR-015 | Service, before any write |
| Methods uppercased, paths normalized before comparison | FR-019 | Drift computation (R7) |
| Rejection leaves the prior spec intact | FR-006 | Validate fully **before** the write; single-statement upsert |

## State transitions

```
(no spec) --PUT valid--> (spec present)
(no spec) --PUT invalid--> (no spec)          [rejected, nothing stored]
(spec present) --PUT valid--> (spec present, replaced, updated_at bumped)
(spec present) --PUT invalid--> (spec present, unchanged)   [FR-006]
(spec present) --DELETE--> (no spec)          [API unaffected, FR-007]
(any) --API deleted--> (row gone via FK cascade)
```

There is no draft/published state on the specification itself — publication state belongs to the API, not to this row.

---

## Derived, non-persisted values

**Drift report** (FR-014) — computed on read, never stored (R7):

| Field | Meaning |
|---|---|
| `inSpecOnly[]` | `{method, path}` pairs declared in the specification but not in the API's declared resources |
| `inResourcesOnly[]` | `{method, path}` pairs declared as API resources but absent from the specification |
| `inSync` | `true` when both lists are empty |

**API response additions** (FR-020) — two response-only fields on the existing `RESTAPI` schema:

| Field | Type | Notes |
|---|---|---|
| `hasOpenApiSpec` | boolean | Present on both get and list, so presence is visible without fetching each document |
| `openApiSpecUpdatedAt` | date-time, nullable | FR-012's last-changed time |

> These are **response-only**. They must not enter `model.RestAPIConfig`, must not be accepted on create/update, and must not reach `BuildAPIDeploymentYAML` — see plan.md gate G2.

---

## Relationships

```
organizations 1 ──── * rest_apis 1 ──── 0..1 rest_api_specs
                          │
                          └── configuration (BLOB) ── Operations[]  ← routing source of truth
                                                          ▲
                                            drift compared against, never written from
```

The specification and the declared operations are compared but never synchronised (FR-013). Routing reads only the left-hand path.
