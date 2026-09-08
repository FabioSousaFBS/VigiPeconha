CREATE TABLE IF NOT EXISTS occurrence_validations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurrence_id UUID NOT NULL,
    validator_user_id UUID NOT NULL,
    status VARCHAR(30) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_validation_occurrence
        FOREIGN KEY (occurrence_id)
        REFERENCES occurrences(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_validation_user
        FOREIGN KEY (validator_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);