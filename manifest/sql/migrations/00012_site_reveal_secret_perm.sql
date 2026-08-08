-- 查看站点 PaaS APPSECRET(AES 解密回显)
-- +goose Up

INSERT INTO manage_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, active_path, sort)
VALUES
 (45, 4, '查看密钥', '/manage/sites/{id}/app-secret', '', 'GET', '', 0, 0, 0, '', 17)
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
DELETE FROM manage_permission WHERE id = 45;
