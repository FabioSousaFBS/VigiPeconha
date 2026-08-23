CREATE TABLE IF NOT EXISTS hospitals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID,
    name VARCHAR(150) NOT NULL,
    address VARCHAR(255),
    neighborhood VARCHAR(120),
    city VARCHAR(120) NOT NULL,
    state VARCHAR(2) NOT NULL,
    phone VARCHAR(30),
    has_antivenom BOOLEAN NOT NULL DEFAULT FALSE,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_hospital_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE SET NULL
);