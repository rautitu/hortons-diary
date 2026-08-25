# Product and deployment decisions

Status: confirmed by Tuomas on 2026-08-25.

## Authentication model

Use local email and password authentication. Hash passwords with Argon2id and
store only the encoded hash and its parameters. Use opaque, server-side
sessions stored in PostgreSQL; the browser receives only a secure session
identifier cookie.

The app therefore owns credential security, session rotation and invalidation,
rate limiting, and eventual password-reset behavior. Password recovery and
email verification are not automatically part of the MVP and need separate
scope before a public launch.

## Registration policy

Registration is closed. Users are invited or created administratively. There
is no public self-registration endpoint in the MVP.

Phase 5 must select a safe initial-user and user-creation workflow. It must not
require inserting password hashes manually into PostgreSQL.

## Initial deployment environment

Deploy on a controlled Linux/Docker host. Run the Go application and
PostgreSQL in separate containers, with PostgreSQL data on a persistent volume.
PostgreSQL is reachable only through the private Compose network and does not
publish a host port.

An existing host-level Caddy installation terminates HTTPS and reverse-proxies
a dedicated application hostname to the app on a private or loopback-only
port. The Caddy configuration and other environment-specific identifiers are
managed outside this repository.

The deployment phase must also define access-controlled off-host backups and
prove that PostgreSQL can be restored from them.

## Decisions intentionally left for later

- password recovery and email verification behavior;
- exact admin/invitation workflow;
- application hostname, network attachment and port;
- backup destination and retention;
- whether MVP needs CSV export;
- account deletion and data-retention semantics;
- hard versus soft deletion for attack records;
- imported Manage My Pain Pro fields beyond the MVP.
