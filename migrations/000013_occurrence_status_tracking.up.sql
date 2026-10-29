ALTER TABLE occurrences
ADD COLUMN status_updated_by_user_id UUID
    REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN status_updated_at TIMESTAMPTZ;
CREATE TABLE occurrence_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    occurrence_id UUID NOT NULL
        REFERENCES occurrences(id) ON DELETE CASCADE,
    changed_by_user_id UUID
        REFERENCES users(id) ON DELETE SET NULL,
    previous_status VARCHAR(30) NOT NULL,
    new_status VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_occurrence_status_history_occurrence
ON occurrence_status_history (occurrence_id, created_at DESC);