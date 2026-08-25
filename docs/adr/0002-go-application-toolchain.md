# ADR 0002: Initial Go application toolchain

- Status: accepted
- Date: 2026-08-25

## Context

Phase 2 needs a narrow, maintainable toolchain for HTTP routing, templates,
PostgreSQL access, schema changes and styling. Exact dependency versions should
be pinned only when implementation starts, so they can be checked against
supported releases at that time.

## Proposed decision

Use Go 1.27 (latest patch), chi v5, `html/template`, pgx v5 with sqlc, Goose v3
SQL migrations, and hand-written semantic CSS. Serve a pinned local copy of
HTMX. Keep normal HTTP navigation and form posts as the functional baseline.

## Rationale and consequences

This selection stays close to Go's standard library, keeps SQL visible and
type-checked, and avoids a separate frontend build toolchain. Code generation
from sqlc becomes a required reproducible development step, and generated files
must be checked consistently in CI. Goose and sqlc configuration add tooling,
but make schema history and queries explicit.

`html/template` is less component-oriented than Templ, and plain CSS offers
fewer utility conventions than Tailwind. Those costs are acceptable for a
small first UI; both choices can be revisited with concrete evidence.

Authentication, registration and deployment are recorded separately in
`docs/decisions.md`.
