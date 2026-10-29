ALTER TABLE occurrences
ADD CONSTRAINT occurrences_status_check
CHECK (
    status IN (
        'pending',
        'acknowledged',
        'addressed',
        'archived'
    )
);