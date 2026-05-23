BEGIN;

DROP INDEX IF EXISTS idx_refresh_sessions_user_id;
DROP INDEX IF EXISTS idx_login_users_username_lower;
DROP INDEX IF EXISTS idx_login_users_email_lower;
DROP INDEX IF EXISTS uq_daily_checkins_user_day;
DROP INDEX IF EXISTS uq_refresh_sessions_token_id;
DROP INDEX IF EXISTS uq_provider_identities_provider_id;
DROP INDEX IF EXISTS uq_login_users_provider_identity;

DROP TABLE IF EXISTS provider_identities;
DROP TABLE IF EXISTS refresh_sessions;
DROP TABLE IF EXISTS daily_checkins_v2;

ALTER TABLE login_users
    DROP COLUMN IF EXISTS provider,
    DROP COLUMN IF EXISTS provider_id,
    DROP COLUMN IF EXISTS points,
    DROP COLUMN IF EXISTS trust_level,
    DROP COLUMN IF EXISTS is_admin,
    DROP COLUMN IF EXISTS desktop_notifications_enabled,
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS last_checkin_at,
    DROP COLUMN IF EXISTS consecutive_days,
    DROP COLUMN IF EXISTS status;

COMMIT;
