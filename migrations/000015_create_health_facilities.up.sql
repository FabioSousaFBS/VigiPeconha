CREATE TABLE health_facilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    cnes VARCHAR(20) UNIQUE,
    unit_type VARCHAR(40) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    address TEXT,
    city VARCHAR(120) NOT NULL,
    state CHAR(2) NOT NULL,
    phone VARCHAR(40),
    is_24_hours BOOLEAN NOT NULL DEFAULT FALSE,
    treats_envenomation BOOLEAN,
    has_antivenom BOOLEAN,
    last_verified_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT health_facilities_unit_type_check CHECK (unit_type IN ('hospital','upa','emergency','reference_center','other')),
    CONSTRAINT health_facilities_state_check CHECK (state ~ '^[A-Z]{2}$')
);
CREATE INDEX idx_health_facilities_location ON health_facilities USING GIST (location);
CREATE INDEX idx_health_facilities_active_city ON health_facilities (is_active, state, city);
