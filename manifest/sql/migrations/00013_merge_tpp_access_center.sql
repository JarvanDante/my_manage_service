-- 合并「系统接入」「接入文档」为「接入中心」；旧 docs 菜单隐藏并指向中心。
-- +goose Up

UPDATE manage_permission
SET name = '接入中心',
    route_url = '/tpp/center',
    component = 'tpp/center',
    icon = 'lucide:key-round',
    hide_in_menu = 0,
    sort = 2,
    updated_at = now()
WHERE id = 102;

UPDATE manage_permission
SET name = '接入文档(已合并)',
    route_url = '/tpp/center',
    component = 'tpp/center',
    hide_in_menu = 1,
    active_path = '/tpp/center',
    updated_at = now()
WHERE id = 103;

-- +goose Down

UPDATE manage_permission
SET name = '系统接入',
    route_url = '/tpp/access',
    component = 'tpp/access',
    icon = 'lucide:key-round',
    hide_in_menu = 0,
    sort = 2,
    updated_at = now()
WHERE id = 102;

UPDATE manage_permission
SET name = '接入文档',
    route_url = '/tpp/docs',
    component = 'tpp/docs',
    icon = 'lucide:book-open-text',
    hide_in_menu = 0,
    active_path = '',
    sort = 3,
    updated_at = now()
WHERE id = 103;
