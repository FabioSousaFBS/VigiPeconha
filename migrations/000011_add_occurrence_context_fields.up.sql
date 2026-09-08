-- =========================================================
-- VigiPeçonha
-- Migration 11
-- Atualiza occurrences para o novo modelo público/web
-- =========================================================


-- ---------------------------------------------------------
-- Organização responsável pela ocorrência
-- NULL para ocorrências públicas do aplicativo
-- ---------------------------------------------------------

ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS organization_id UUID;
-- ---------------------------------------------------------
-- Usuário autenticado que criou a ocorrência
-- NULL para ocorrências públicas
-- ---------------------------------------------------------
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS created_by_user_id UUID;
-- ---------------------------------------------------------
-- Origem da ocorrência
--
-- mobile_public
-- web
-- ---------------------------------------------------------
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS source VARCHAR(30);
-- ---------------------------------------------------------
-- Tipo geral do animal
-- Ex:
-- scorpion
-- snake
-- spider
-- ---------------------------------------------------------
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS animal_type VARCHAR(50);
-- ---------------------------------------------------------
-- Dados opcionais de quem realizou o registro
-- ---------------------------------------------------------
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS reporter_name VARCHAR(150);
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS reporter_phone VARCHAR(30);
-- ---------------------------------------------------------
-- Informações textuais de endereço
-- ---------------------------------------------------------
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS address VARCHAR(255);
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS neighborhood VARCHAR(120);
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS city VARCHAR(120);
ALTER TABLE occurrences
ADD COLUMN IF NOT EXISTS state VARCHAR(2);
-- =========================================================
-- Compatibilidade com modelo antigo
-- =========================================================
-- Se a ocorrência antiga possui user_id,
-- aproveitamos esse valor no novo campo.
UPDATE occurrences
SET created_by_user_id = user_id
WHERE created_by_user_id IS NULL
  AND user_id IS NOT NULL;
-- ---------------------------------------------------------
-- Ocorrências antigas não possuem source.
-- Como foram registradas antes da separação web/mobile,
-- usaremos "legacy" para não classificá-las incorretamente.
-- ---------------------------------------------------------
UPDATE occurrences
SET source = 'legacy'
WHERE source IS NULL;
-- source passa a ser obrigatório depois do backfill.
ALTER TABLE occurrences
ALTER COLUMN source SET NOT NULL;
-- ---------------------------------------------------------
-- animal_type
--
-- Não devemos inventar o animal das ocorrências antigas.
-- Por isso usamos "unknown" somente para registros legados.
-- Novos registros serão validados pelo backend.
-- ---------------------------------------------------------
UPDATE occurrences
SET animal_type = 'unknown'
WHERE animal_type IS NULL;
ALTER TABLE occurrences
ALTER COLUMN animal_type SET NOT NULL;
-- =========================================================
-- Foreign Keys
-- =========================================================
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_occurrence_organization'
    ) THEN

        ALTER TABLE occurrences
        ADD CONSTRAINT fk_occurrence_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE RESTRICT;

    END IF;
END $$;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_occurrence_created_by_user'
    ) THEN
        ALTER TABLE occurrences
        ADD CONSTRAINT fk_occurrence_created_by_user
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL;
    END IF;
END $$;
-- =========================================================
-- Índices
-- =========================================================
CREATE INDEX IF NOT EXISTS idx_occurrences_organization
ON occurrences(organization_id);
CREATE INDEX IF NOT EXISTS idx_occurrences_created_by_user
ON occurrences(created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_occurrences_source
ON occurrences(source);
CREATE INDEX IF NOT EXISTS idx_occurrences_animal_type
ON occurrences(animal_type);