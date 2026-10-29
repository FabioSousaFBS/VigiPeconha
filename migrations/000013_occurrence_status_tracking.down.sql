DROP TABLE IF EXISTS occurrence_status_history;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS status_updated_at,
DROP COLUMN IF EXISTS status_updated_by_user_id;