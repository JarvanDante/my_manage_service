-- 站点域名用途: admin / h5 / app
-- +goose Up
ALTER TABLE site_domain ADD COLUMN IF NOT EXISTS kind varchar(16) NOT NULL DEFAULT '';
COMMENT ON COLUMN site_domain.kind IS '域名用途: admin / h5 / app';
CREATE INDEX IF NOT EXISTS idx_site_domain_kind ON site_domain (site_id, kind);

-- +goose Down
DROP INDEX IF EXISTS idx_site_domain_kind;
ALTER TABLE site_domain DROP COLUMN IF EXISTS kind;
