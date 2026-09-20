CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    email TEXT NOT NULL CHECK (btrim(email) <> ''),
    password_hash TEXT NOT NULL CHECK (password_hash <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_unique ON users (lower(email));

CREATE TABLE pain_record (
    user_id UUID NOT NULL REFERENCES users (user_id),
    pain_record_id UUID PRIMARY KEY,
    start_time TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    severity SMALLINT NOT NULL CHECK (severity BETWEEN 1 AND 10),
    row_updated TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX pain_record_user_start_time ON pain_record (user_id, start_time);

CREATE FUNCTION set_pain_record_updated() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.row_updated := clock_timestamp();
    RETURN NEW;
END;
$$;
CREATE TRIGGER pain_record_updated
BEFORE UPDATE ON pain_record
FOR EACH ROW EXECUTE FUNCTION set_pain_record_updated();

CREATE TABLE sessions (
    session_token_hash BYTEA PRIMARY KEY CHECK (octet_length(session_token_hash) = 32),
    user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    CHECK (expires_at > created_at)
);
CREATE INDEX sessions_user_id ON sessions (user_id);
CREATE INDEX sessions_expires_at ON sessions (expires_at);
