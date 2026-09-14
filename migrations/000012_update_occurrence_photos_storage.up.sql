-- =========================================================
-- VigiPeçonha
-- Migration 12
-- Adequa occurrence_photos para armazenamento no Cloudflare R2
-- =========================================================

ALTER TABLE occurrence_photos
ADD COLUMN IF NOT EXISTS storage_key TEXT;
ALTER TABLE occurrence_photos
ADD COLUMN IF NOT EXISTS content_type VARCHAR(100);
ALTER TABLE occurrence_photos
ADD COLUMN IF NOT EXISTS file_size BIGINT;
-- O bucket será privado.
-- Portanto não teremos necessariamente uma URL pública permanente.
ALTER TABLE occurrence_photos
ALTER COLUMN photo_url DROP NOT NULL;
-- Registros antigos, caso existam, podem usar a própria URL
-- temporariamente como referência de storage.
UPDATE occurrence_photos
SET storage_key = photo_url
WHERE storage_key IS NULL
  AND photo_url IS NOT NULL;
-- Para dados antigos que não tenham URL, usamos uma chave legacy.
UPDATE occurrence_photos
SET storage_key = 'legacy/' || id::text
WHERE storage_key IS NULL;
ALTER TABLE occurrence_photos
ALTER COLUMN storage_key SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_occurrence_photos_occurrence
ON occurrence_photos(occurrence_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_occurrence_photos_storage_key
ON occurrence_photos(storage_key);