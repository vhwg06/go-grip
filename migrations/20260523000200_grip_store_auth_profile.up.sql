BEGIN;

CREATE TABLE IF NOT EXISTS login_users (
    id TEXT PRIMARY KEY,
    provider TEXT,
    provider_id TEXT,
    username TEXT,
    email TEXT,
    password_hash TEXT,
    role_id TEXT,
    role TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    points INTEGER NOT NULL DEFAULT 0,
    trust_level INTEGER NOT NULL DEFAULT 0,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    desktop_notifications_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMPTZ,
    last_checkin_at TIMESTAMPTZ,
    consecutive_days INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE login_users
    ADD COLUMN IF NOT EXISTS provider TEXT,
    ADD COLUMN IF NOT EXISTS provider_id TEXT,
    ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS trust_level INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS desktop_notifications_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_checkin_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS consecutive_days INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

CREATE TABLE IF NOT EXISTS provider_identities (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    email TEXT,
    username TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token_id TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS daily_checkins_v2 (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    checkin_date TIMESTAMPTZ NOT NULL,
    reward INTEGER NOT NULL DEFAULT 0,
    streak_after INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_login_users_provider_identity
    ON login_users (provider, provider_id)
    WHERE provider IS NOT NULL AND provider_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_provider_identities_provider_id
    ON provider_identities (provider, provider_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_refresh_sessions_token_id
    ON refresh_sessions (token_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_daily_checkins_user_day
    ON daily_checkins_v2 (user_id, (DATE(checkin_date)));

CREATE INDEX IF NOT EXISTS idx_login_users_email_lower
    ON login_users (LOWER(email));
CREATE INDEX IF NOT EXISTS idx_login_users_username_lower
    ON login_users (LOWER(username));
CREATE INDEX IF NOT EXISTS idx_refresh_sessions_user_id
    ON refresh_sessions (user_id);

COMMIT;
