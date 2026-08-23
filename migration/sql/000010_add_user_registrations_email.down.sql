DROP INDEX IF EXISTS idx_user_registrations_email;
ALTER TABLE user_registrations DROP COLUMN IF EXISTS email;
