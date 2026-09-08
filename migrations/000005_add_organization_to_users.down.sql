ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_user_organization;
ALTER TABLE users
DROP COLUMN IF EXISTS organization_id;