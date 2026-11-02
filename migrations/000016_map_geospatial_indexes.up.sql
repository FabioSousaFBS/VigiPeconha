-- Indexes for viewport searches and map filters.
CREATE INDEX IF NOT EXISTS idx_occurrences_location_geometry ON occurrences USING GIST ((location::geometry));
CREATE INDEX IF NOT EXISTS idx_health_facilities_location_geometry ON health_facilities USING GIST ((location::geometry));
CREATE INDEX IF NOT EXISTS idx_occurrences_map_filters ON occurrences (occurred_at DESC, status, animal_type);
CREATE INDEX IF NOT EXISTS idx_health_facilities_map_filters ON health_facilities (is_active, state, unit_type);
