-- 站点详情改为顶级隐藏菜单, 并设置 active_path, 避免「站点管理」被 redirect 到无 id 的详情页。
-- +goose Up
UPDATE manage_permission
SET parent_id = 0,
    active_path = '/site',
    hide_in_menu = 1,
    updated_at = now()
WHERE id = 5
  AND route_url = '/site/detail';

-- +goose Down
UPDATE manage_permission
SET parent_id = 4,
    active_path = '',
    updated_at = now()
WHERE id = 5
  AND route_url = '/site/detail';
