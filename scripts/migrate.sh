#!/bin/sh
# Run versioned SQL files against the Compose database. Never source credentials.
set -eu
cd "$(dirname "$0")/.."
for migration in migrations/[0-9][0-9][0-9][0-9]_*.sql; do
    [ -f "$migration" ] || continue
    version=$(basename "$migration" .sql)
    case "$version" in *[!a-zA-Z0-9_]*) echo "Invalid migration name" >&2; exit 1 ;; esac
    echo "Checking migration $version"
    {
        cat <<'SQL'
BEGIN;
SELECT pg_advisory_xact_lock(724091821);
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
SQL
        printf "SELECT NOT EXISTS (SELECT 1 FROM schema_migrations WHERE version = '%s') AS apply_migration \\gset\n" "$version"
        printf '\\if :apply_migration\n'
        cat "$migration"
        printf "\nINSERT INTO schema_migrations (version) VALUES ('%s');\n" "$version"
        printf '\\endif\nCOMMIT;\n'
    } | docker compose exec -T db sh -c 'exec psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
done
