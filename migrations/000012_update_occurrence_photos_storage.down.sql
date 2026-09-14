DROP INDEX IF EXISTS idx_occurrence_photos_storage_key;
-- Só podemos voltar photo_url para NOT NULL se todos os registros
-- possuírem valor.
UPDATE occurrence_photos
SET photo_url = storage_key
WHERE photo_url IS NULL;
ALTER TABLE occurrence_photos
ALTER COLUMN photo_url SET NOT NULL;
ALTER TABLE occurrence_photos
DROP COLUMN IF EXISTS file_size;
ALTER TABLE occurrence_photos
DROP COLUMN IF EXISTS content_type;
ALTER TABLE occurrence_photos
DROP COLUMN IF EXISTS storage_key;