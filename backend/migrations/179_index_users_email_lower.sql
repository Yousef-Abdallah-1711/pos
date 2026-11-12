-- The auth repository normalizes email identifiers with lower(email) for
-- case-insensitive login and race-safe registration checks.
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON users (lower(email));
