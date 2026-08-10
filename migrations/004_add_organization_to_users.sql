ALTER TABLE users
ADD COLUMN IF NOT EXISTS organization_id UUID;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS status VARCHAR(30)
NOT NULL DEFAULT 'active';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_user_organization'
    ) THEN
        ALTER TABLE users
        ADD CONSTRAINT fk_user_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE RESTRICT;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_organization
ON users(organization_id);