CREATE TABLE IF NOT EXISTS occurrences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID,
    created_by_user_id UUID,
    source VARCHAR(30) NOT NULL,
    occurrence_type VARCHAR(30) NOT NULL,
    animal_type VARCHAR(50) NOT NULL,
    species VARCHAR(120),
    description TEXT,
    reporter_name VARCHAR(150),
    reporter_phone VARCHAR(30),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    address VARCHAR(255),
    neighborhood VARCHAR(120),
    city VARCHAR(120),
    state VARCHAR(2),
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_occurrence_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_occurrence_user
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);