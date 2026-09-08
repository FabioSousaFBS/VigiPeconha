CREATE TABLE IF NOT EXISTS occurrence_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurrence_id UUID NOT NULL,
    storage_key VARCHAR(500) NOT NULL,
    file_url TEXT,
    content_type VARCHAR(100),
    file_size BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_photo_occurrence
        FOREIGN KEY (occurrence_id)
        REFERENCES occurrences(id)
        ON DELETE CASCADE
);