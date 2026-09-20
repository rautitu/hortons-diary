-- Run with psql -X -v ON_ERROR_STOP=1. All test data is rolled back.
BEGIN;
DO $$
DECLARE
    uid UUID := gen_random_uuid();
    rid UUID := gen_random_uuid();
    original_updated TIMESTAMPTZ;
BEGIN
    INSERT INTO users(user_id, email, password_hash)
    VALUES (uid, uid::text || '@example.test', 'test-only-not-a-real-hash');
    BEGIN
        INSERT INTO users(user_id, email, password_hash)
        VALUES (gen_random_uuid(), upper(uid::text || '@example.test'), 'test');
        RAISE EXCEPTION 'case-insensitive uniqueness missing';
    EXCEPTION WHEN unique_violation THEN NULL;
    END;
    INSERT INTO pain_record(user_id, pain_record_id, start_time, duration_minutes, severity)
    VALUES (uid, rid, now(), 30, 7) RETURNING row_updated INTO original_updated;
    UPDATE pain_record SET severity = 8 WHERE pain_record_id = rid;
    IF NOT EXISTS (SELECT 1 FROM pain_record WHERE pain_record_id = rid AND row_updated > original_updated) THEN
        RAISE EXCEPTION 'update timestamp did not advance';
    END IF;
    BEGIN
        UPDATE pain_record SET severity = 11 WHERE pain_record_id = rid;
        RAISE EXCEPTION 'severity check missing';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE pain_record SET duration_minutes = 0 WHERE pain_record_id = rid;
        RAISE EXCEPTION 'duration check missing';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE pain_record SET user_id = gen_random_uuid() WHERE pain_record_id = rid;
        RAISE EXCEPTION 'user foreign key missing';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    INSERT INTO sessions(session_token_hash, user_id, expires_at)
    VALUES (decode(replace(gen_random_uuid()::text, '-', '') || replace(gen_random_uuid()::text, '-', ''), 'hex'), uid, now() + interval '1 day');
    BEGIN
        DELETE FROM users WHERE user_id = uid;
        RAISE EXCEPTION 'record owner deletion should be restricted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    DELETE FROM pain_record WHERE pain_record_id = rid;
    DELETE FROM users WHERE user_id = uid;
    IF EXISTS (SELECT 1 FROM sessions WHERE user_id = uid) THEN
        RAISE EXCEPTION 'session cascade missing';
    END IF;
END;
$$;
ROLLBACK;
