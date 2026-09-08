ALTER TABLE users
ADD COLUMN IF NOT EXISTS organization_id UUID;
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
