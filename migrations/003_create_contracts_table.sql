CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    organization_id UUID NOT NULL,

    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    user_limit INTEGER NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'active',
    payment_status VARCHAR(30) NOT NULL DEFAULT 'paid',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_contract_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_contract_user_limit
        CHECK (user_limit > 0),

    CONSTRAINT chk_contract_dates
        CHECK (end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_contracts_organization
ON contracts(organization_id);

CREATE INDEX IF NOT EXISTS idx_contracts_status
ON contracts(status);