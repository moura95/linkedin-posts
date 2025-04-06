-- Remover constraints
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS fk_ticket_severity;
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS fk_ticket_category;
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS fk_ticket_subcategory;

-- Remover índices
DROP INDEX IF EXISTS idx_ticket_status;
DROP INDEX IF EXISTS idx_ticket_severity;
DROP INDEX IF EXISTS idx_ticket_category;

-- Remover tabelas
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS severities;