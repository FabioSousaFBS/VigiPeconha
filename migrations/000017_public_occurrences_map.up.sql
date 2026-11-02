CREATE OR REPLACE VIEW public_occurrences_map AS
SELECT
    id,
    ST_Transform(
        ST_SnapToGrid(
            ST_Transform(location::geometry, 3857),
            1000
        ),
        4326
    ) AS public_location,
    animal_type,
    occurrence_type,
    status,
    occurred_at
FROM occurrences
WHERE status <> 'archived';
