-- 总后台管理员谷歌验证器(TOTP)
-- +goose Up
ALTER TABLE manage_admin
    ADD COLUMN totp_secret varchar(64) NOT NULL DEFAULT '',
    ADD COLUMN totp_bound_at timestamptz;

COMMENT ON COLUMN manage_admin.totp_secret IS '谷歌验证器 TOTP 密钥(已绑定才写入)';
COMMENT ON COLUMN manage_admin.totp_bound_at IS '谷歌验证器绑定时间';

INSERT INTO manage_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, active_path, sort)
VALUES
 (46, 9, '解绑验证器', '/manage/admins/{id}/totp', '', 'DELETE', '', 0, 0, 0, '', 5)
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
DELETE FROM manage_permission WHERE id = 46;
ALTER TABLE manage_admin
    DROP COLUMN IF EXISTS totp_bound_at,
    DROP COLUMN IF EXISTS totp_secret;
