## Context

See `proposal.md` — Why. The design-relevant state of the code today:

- REST APIs live in `platform-api`. The contract is `resources/openapi.yaml`; `make generate` runs `oapi-codegen` v2.5.1 over it into `api/generated.go`, so request/response types are generated, not hand-written. Types unreachable from a path are added by hand in `api/manual_types.go`.
- Scope enforcement is driven from that same spec: `middleware.LoadScopeRegistry(cfg.OpenAPISpecPath)` parses `openapi.yaml` at startup (with `gopkg.in/yaml.v3`) and `ScopeEnforcer` matches a held scope against an operation's accepted-scope list by **exact string equality** — so every new operation must declare its scopes in the spec, and there is no wildcard form (`.claude/rules/authentication_authorization.md`, GO-AUTH-020).
- Handlers register `net/http` patterns on a `router.Router` recorder, e.g. `constants.APIBasePath + "/rest-apis/{restApiId}/api-keys"`, wrapped in `middleware.MapErrors` (`internal/handler/api_key.go`). Route matching is structural, satisfying GO-AUTH-004/017.
- `rest_apis` rows carry the API's structured fields plus a `configuration BLOB`; the row's `uuid` is also its `artifacts.uuid` (FK with `ON DELETE CASCADE`). The table has shipped, so it is frozen to additive changes only (`.claude/rules/db-schema-changes.md`, R0-FROZEN).
- `internal/handler/secret.go` already handles `multipart/form-data`, via `r.ParseMultipartForm(32 << 20)` with a hardcoded ceiling.
- `platform-api` has **no** OpenAPI parsing dependency. `github.com/getkin/kin-openapi v0.133.0` is already in `gateway/gateway-controller` and `event-gateway/gateway-controller`.
- The API Portal's `POST /apis` takes a `definition` file part (`portals/api-portal/docs/api-portal-openapi-spec-v0.9.yaml`), which is what the stored document ultimately feeds.

## Goals / Non-Goals

**Goals:**

- Store the uploaded document as opaque bytes plus a small amount of derived metadata, so retrieval can be byte-exact.
- Keep validation strictly a gate: parse to decide accept/reject, then discard the parsed model. Nothing downstream reads the parsed form.
- Add exactly one new table, additively, across all four SQL dialects, with the document removed by the existing artifact cascade.
- Introduce the parser dependency in a single wrapper package so a future swap (or a second artifact kind) touches one file.

**Non-Goals:**

- Any coupling between the document and `operations`, deployment, or the gateway translator. Nothing outside the new handler/service/repository path reads the stored bytes in this change.
- Serving the document to the developer portal. That is a separate change; this one only makes the document retrievable.
- Reference resolution. `$ref`s in the uploaded document are validated as far as the parser does so internally, and are never fetched over the network.

## Decisions

### D1: A separate `rest_api_openapi_specs` table, not a column on `rest_apis`

One row per REST API, keyed by `rest_api_uuid` as the primary key (so "at most one document" is enforced by the schema, not by application code), with `FOREIGN KEY (rest_api_uuid) REFERENCES rest_apis(uuid) ON DELETE CASCADE`. Columns: the document bytes (`BLOB`/`BYTEA`/`VARBINARY(MAX)` per dialect), `content_type`, `byte_size`, `openapi_version`, `organization_uuid`, and the standard `created_by`/`created_at`/`updated_by`/`updated_at` audit set.

*Why:* `rest_apis` is a shipped, frozen table — adding a column there is permitted (nullable/defaulted) but would need a per-dialect `ALTER TABLE` for already-provisioned databases (R0-UPGRADE-PATH). A **new** table needs only its guarded `CREATE TABLE IF NOT EXISTS`, which is correct against fresh and existing databases alike, with no `ALTER` path to get wrong. It also keeps a potentially large blob out of every `SELECT` on the hot `rest_apis` row.

*Alternatives:* (a) A column on `rest_apis` — rejected for the `ALTER` burden and the read amplification. (b) Reusing the existing `configuration BLOB` — rejected: that field has a defined meaning for gateway translation, and stuffing an unrelated document into it invites exactly the operations/document coupling the spec forbids. (c) Filesystem or object storage — rejected: introduces a second durability domain and a path-traversal surface (`.claude/rules/file-access.md`) for no benefit at these sizes.

`organization_uuid` is denormalized onto the row so every query filters on the token-derived organization directly (GO-AUTH-005) rather than trusting a join through a path-supplied identifier.

### D2: `github.com/getkin/kin-openapi` for validation, behind a thin internal wrapper

Add `kin-openapi v0.133.0` to `platform-api/go.mod` and confine every call to it to one new internal package (e.g. `internal/openapispec`) exposing roughly `Validate(ctx, raw []byte) (Metadata, error)`.

*Why:* It is already the repo's OpenAPI parser at the same version in two other modules, so it is a known quantity for `.claude/rules/dependency-management.md` (MIT, maintained) rather than a fresh evaluation. The wrapper keeps the parsed `*openapi3.T` from leaking into service or handler code, which is what makes "validation is a gate, not a transform" enforceable rather than merely intended.

*Alternatives:* `pb33f/libopenapi` — capable and 3.x-native, but new to the repo and buys nothing here. Hand-rolled structural checks against `gopkg.in/yaml.v3` — rejected; re-implementing OpenAPI validation is exactly the kind of drift `.claude/rules/ssrf-prevention.md` directive 6 warns about in a different guise.

**Loader must be sealed.** `kin-openapi`'s loader can fetch external `$ref`s over the network — an SSRF primitive driven by uploaded content (`.claude/rules/ssrf-prevention.md`). Configure the loader with `IsExternalRefsAllowed = false` and a `ReadFromURIFunc` that returns an error unconditionally, and assert that in a unit test so a dependency bump cannot silently re-enable it. Validation runs under a `context.WithTimeout` so a pathological document fails fast instead of holding a worker (`.claude/rules/go-network-service-hardening.md`).

**Version gate before validation.** Read the top-level `openapi` / `swagger` key first (a cheap `yaml.v3` decode into a two-field struct). `swagger: "2.0"` → reject with the "OpenAPI 3.x only" message. A missing or non-`3.x` `openapi` key → reject. Only a `3.x` document reaches the full validator. This makes the Swagger-2.0 rejection message deterministic rather than dependent on how the 3.x validator happens to fail.

### D3: Content type is sniffed from the bytes, never taken from the client

Decide JSON vs YAML by attempting a strict JSON decode of the (whitespace-trimmed) bytes first; success ⇒ `application/json`, failure ⇒ attempt YAML ⇒ `application/yaml`, both failing ⇒ `400`. The multipart part's declared `Content-Type` and its filename extension are ignored entirely.

*Why:* `.claude/rules/file-access.md` directive 6 and `js-file-access.md` both require sniffing over declared type. Because JSON is a subset of YAML, order matters: JSON-first means a JSON document is never mislabeled `application/yaml`. The stored `content_type` is what retrieval echoes, so getting this wrong at ingest is what a consumer sees forever.

Note there is no magic-byte test for JSON vs YAML — a parse attempt *is* the sniff here, which is why the parse result, not the declared type, is authoritative.

### D4: Size ceiling from config, enforced before the body is buffered

Add `openapi_spec.max_upload_bytes` to `config.Server` with a safe default (2 MiB) and wire it into the handler as `http.MaxBytesReader(w, r.Body, limit)` **before** `ParseMultipartForm`, with `ParseMultipartForm` given a memory bound no larger than the same limit. On overflow return `413` with a generic message that never states the configured limit.

*Why:* `.claude/rules/file-access.md` directive 5 and `go-network-service-hardening.md` directive 1 both require a config-sourced ceiling and a non-leaky rejection. The existing `secret.go` hardcodes `32 << 20`; this change does **not** retrofit that (out of scope), but the new path must not copy the pattern.

The ceiling is enforced twice on purpose: `MaxBytesReader` caps what can be read at all, and a post-read `len(bytes) > limit` check guards the case where the multipart framing lets a part through under a different accounting.

### D5: Four operations under `/rest-apis/{restApiId}/openapi`

| Operation | Method + path | Scopes | Success |
|---|---|---|---|
| Upload/replace | `PUT /rest-apis/{restApiId}/openapi` | `ap:rest_api:update`, `ap:rest_api:manage` | `201` first attach, `200` replace |
| Retrieve document | `GET /rest-apis/{restApiId}/openapi` | `ap:rest_api:read`, `ap:rest_api:manage` | `200`, raw body |
| Retrieve metadata | `GET /rest-apis/{restApiId}/openapi/metadata` | `ap:rest_api:read`, `ap:rest_api:manage` | `200`, JSON |
| Remove | `DELETE /rest-apis/{restApiId}/openapi` | `ap:rest_api:update`, `ap:rest_api:manage` | `204` |

*Why `PUT`, not `POST`:* the sub-resource is a singleton at a known address, and upload is idempotent-by-replacement. `POST` would imply a collection that does not exist.

*Why a separate metadata path rather than content negotiation:* the spec requires both a byte-exact document response and a JSON metadata response. Selecting between them on `Accept` makes the byte-exactness contingent on a header a proxy may rewrite, and makes the scope registry's per-operation mapping ambiguous. A distinct path keeps each operation single-valued. `HEAD` was considered for metadata and rejected — headers alone cannot carry the audit fields.

Scopes reuse the REST API's own scopes rather than introducing `ap:rest_api_spec:*`: the document is part of the API, not a separately-administered resource, and GO-AUTH-020 warns against minting scopes that are not bound to a distinct resource type. All four must be declared in `resources/openapi.yaml` and mapped in `resources/role-to-scope-mapping.yaml`, since `ScopeEnforcer` matches by exact equality with no wildcard form.

### D6: `404` for cross-organization and for missing, uniformly

The service resolves the REST API by handle **scoped to the organization from verified claims**, and returns the same `404` whether the API does not exist, exists in another organization, or exists with no document attached where a document was required. Ordering matters: authorization and organization scoping resolve *before* the document lookup, so a caller can never distinguish "API exists elsewhere" from "API does not exist" by timing or status.

*Why:* GO-AUTH-005 plus `.claude/rules/error-handling.md` directive 4/5 — an existence-sensitive lookup must not become an enumeration primitive.

### D7: Codegen keeps working for a raw-bytes response

`oapi-codegen` generates types from the spec; a `PUT` with a `multipart/form-data` body and a `GET` returning `application/json` **or** `application/yaml` as an opaque string are both expressible, but the generated types are of limited use for a raw body. Declare the request body as `multipart/form-data` with a single binary `file` property, and the document response as `application/json`/`application/yaml` with `type: string, format: binary`; the handler reads `r.MultipartForm` and writes bytes directly rather than round-tripping through a generated struct. Only the metadata response and the error payloads use generated types. If a needed type ends up unreachable from a path, it goes in `api/manual_types.go`, per the existing convention.

The spec must also gain a `PayloadTooLarge` (`413`) and an `UnsupportedMediaType` (`415`) entry under `components/responses` — neither exists today.

## Risks / Trade-offs

- **A stored document can contradict the API's `operations`.** This is a deliberate consequence of the decision that the document is documentation only, and the spec says so explicitly — but a portal consumer will read a contract that gateway routing does not enforce. → Mitigate by returning the upload's derived metadata (path/operation count) in the upload response so the discrepancy is visible at upload time, and by leaving the door open for an opt-in reconciliation endpoint in a later change. Not silently reconciling is the point; not surfacing the gap would be the defect.
- **`kin-openapi` rejects documents other tools accept.** Strict 3.x validation will turn away real-world specs that Swagger UI renders fine. → Mitigate by returning the validator's own message (sanitized of paths/stack traces per `.claude/rules/error-handling.md`) so the uploader can see what failed, and by pinning the version so acceptance does not shift under a routine dependency bump. Loosening the gate is a future decision with the spec in hand, not a silent one.
- **A document-fetching `$ref` becomes SSRF.** → Mitigated by D2's sealed loader plus a test asserting the loader refuses external reads; the check is structural (the read function errors) rather than a denylist.
- **Blob growth in the primary database.** 2 MiB × every REST API is bounded but not free, and a `SELECT *` on the new table is expensive. → Mitigate by never joining the document table into list queries — retrieval reads it by primary key only — and by keeping the metadata columns in the same row so `GET .../metadata` can project without touching the blob.
- **A second artifact kind will want this.** MCP proxies are an explicit non-goal, but the shape is identical. → Mitigate by keeping the parser wrapper and validation entirely free of REST-API concepts, so extending it later is a new table and a new route, not a refactor.

## Migration Plan

1. Ship the guarded `CREATE TABLE IF NOT EXISTS` in all four dialect files. No `ALTER` is required anywhere, since no shipped table is touched.
2. Deploy. Existing REST APIs have no row in the new table and behave exactly as before; every endpoint added here is new, so there is no version skew for existing clients.
3. Rollback is a code rollback. The orphaned table is inert — nothing outside the removed code reads or writes it — so it can be left in place, and its `ON DELETE CASCADE` continues to keep it consistent with `rest_apis` even while unused.

## Open Questions

- What byte ceiling is right in practice? 2 MiB is a defensible default, but the number should be revisited once a few real customer specs have been uploaded. Deferrable: it is a config value, and the spec requires only that a configured ceiling exists and that rejection does not disclose it.
