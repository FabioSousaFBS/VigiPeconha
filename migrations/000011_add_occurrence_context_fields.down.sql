DROP INDEX IF EXISTS idx_occurrences_animal_type;
DROP INDEX IF EXISTS idx_occurrences_source;
DROP INDEX IF EXISTS idx_occurrences_created_by_user;
DROP INDEX IF EXISTS idx_occurrences_organization;

ALTER TABLE occurrences
DROP CONSTRAINT IF EXISTS fk_occurrence_created_by_user;

ALTER TABLE occurrences
DROP CONSTRAINT IF EXISTS fk_occurrence_organization;

ALTER TABLE occurrences
DROP COLUMN IF EXISTS state;

ALTER TABLE occurrences
DROP COLUMN IF EXISTS city;

ALTER TABLE occurrences
DROP COLUMN IF EXISTS neighborhood;

ALTER TABLE occurrences
DROP COLUMN IF EXISTS address;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS reporter_phone;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS reporter_name;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS animal_type;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS source;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS created_by_user_id;
ALTER TABLE occurrences
DROP COLUMN IF EXISTS organization_id;