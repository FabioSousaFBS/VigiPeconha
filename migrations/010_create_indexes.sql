CREATE INDEX IF NOT EXISTS idx_users_status
ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_organization
ON users(organization_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_document
ON organizations(document)
WHERE document IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_organizations_status
ON organizations(status);
CREATE INDEX IF NOT EXISTS idx_contracts_organization
ON contracts(organization_id);
CREATE INDEX IF NOT EXISTS idx_contracts_status
ON contracts(status);
CREATE INDEX IF NOT EXISTS idx_contracts_payment_status
ON contracts(payment_status);
CREATE INDEX IF NOT EXISTS idx_contracts_dates
ON contracts(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_occurrences_location
ON occurrences
USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_occurrences_status
ON occurrences(status);
CREATE INDEX IF NOT EXISTS idx_occurrences_source
ON occurrences(source);
CREATE INDEX IF NOT EXISTS idx_occurrences_type
ON occurrences(occurrence_type);
CREATE INDEX IF NOT EXISTS idx_occurrences_animal_type
ON occurrences(animal_type);
CREATE INDEX IF NOT EXISTS idx_occurrences_occurred_at
ON occurrences(occurred_at);
CREATE INDEX IF NOT EXISTS idx_occurrences_organization
ON occurrences(organization_id);
CREATE INDEX IF NOT EXISTS idx_occurrences_user
ON occurrences(created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_occurrence_photos_occurrence
ON occurrence_photos(occurrence_id);
CREATE INDEX IF NOT EXISTS idx_validations_occurrence
ON occurrence_validations(occurrence_id);
CREATE INDEX IF NOT EXISTS idx_validations_user
ON occurrence_validations(validator_user_id);
CREATE INDEX IF NOT EXISTS idx_hospitals_location
ON hospitals
USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_hospitals_antivenom
ON hospitals(has_antivenom);
CREATE INDEX IF NOT EXISTS idx_hospitals_status
ON hospitals(status);
CREATE INDEX IF NOT EXISTS idx_hospitals_organization
ON hospitals(organization_id);