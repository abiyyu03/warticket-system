ALTER TABLE events DROP CONSTRAINT IF EXISTS ck_events_quota_remaining_le_quota;
ALTER TABLE events DROP CONSTRAINT IF EXISTS ck_events_quota_remaining_non_negative;
ALTER TABLE events DROP COLUMN IF EXISTS quota_remaining;
