-- 开站即发 PaaS 调用凭证: app_key(明文可展示) + app_secret(AES 落库, 仅创建/重置时回传明文)
-- +goose Up

ALTER TABLE site ADD COLUMN IF NOT EXISTS app_key varchar(64) NOT NULL DEFAULT '';
ALTER TABLE site ADD COLUMN IF NOT EXISTS app_secret varchar(512) NOT NULL DEFAULT '';
COMMENT ON COLUMN site.app_key IS 'PaaS APPKEY, 开站自动生成';
COMMENT ON COLUMN site.app_secret IS 'PaaS APPSECRET, AES-GCM 加密存储(同 db_pass)';

CREATE UNIQUE INDEX IF NOT EXISTS uk_site_app_key ON site (app_key) WHERE app_key <> '';

-- 接口权限: 重置密钥
INSERT INTO manage_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, active_path, sort)
VALUES
 (44, 4, '重置密钥', '/manage/sites/{id}/reset-secret', '', 'POST', '', 0, 0, 0, '', 16)
ON CONFLICT (id) DO UPDATE SET
    parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    route_url = EXCLUDED.route_url,
    method = EXCLUDED.method,
    is_menu = EXCLUDED.is_menu,
    sort = EXCLUDED.sort,
    updated_at = now();

SELECT setval(pg_get_serial_sequence('manage_permission', 'id'), (SELECT MAX(id) FROM manage_permission));

-- +goose Down
DELETE FROM manage_permission WHERE id = 44;
DROP INDEX IF EXISTS uk_site_app_key;
ALTER TABLE site DROP COLUMN IF EXISTS app_secret;
ALTER TABLE site DROP COLUMN IF EXISTS app_key;
