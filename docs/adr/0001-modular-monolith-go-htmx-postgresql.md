# ADR 0001: Modular monolith with Go, HTMX and PostgreSQL

- Status: accepted as project direction
- Date: 2026-08-25

## Context

Hortons Diary starts with one cohesive domain: authenticated users recording
and reviewing their own headache attacks. It needs reliable relational storage
and a mobile-first server-driven interface. The team also wants to explore Go
and HTMX without adopting a JavaScript SPA.

## Decision

Use a modular Go monolith that renders HTML pages and HTMX fragments, backed by
PostgreSQL. Run the application and database as separate deployable containers
or services, but keep one application codebase and process boundary.

## Consequences

The application has simple in-process transactions, authorization checks and
deployment. HTML remains progressively enhanced and core forms can work without
JavaScript. Internal modules must still have clear boundaries so reporting and
imports do not become coupled to HTTP handlers.

This rejects microservices for the initial product. They would add network
contracts, distributed authentication and operational overhead without an
independent scaling or ownership need. Services can be extracted later if a
measured boundary emerges.
