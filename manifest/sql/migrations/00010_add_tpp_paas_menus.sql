-- TPP 门户 + 平台服务菜单骨架(总后台中台入口)。
-- 能力服务仍独立部署;此处仅挂目录/文档/后续配置页。超管自动可见全量菜单。
-- +goose Up

INSERT INTO manage_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, active_path, sort)
VALUES
 -- TPP 门户
 (100, 0,   'TPP门户',   '/tpp',            '',                 '', 'lucide:layout-grid',     1, 0, 0, '', 30),
 (101, 100, '能力概览',   '/tpp/overview',   'tpp/overview',     '', 'lucide:panels-top-left', 1, 0, 0, '', 1),
 (102, 100, '系统接入',   '/tpp/access',     'tpp/access',       '', 'lucide:key-round',       1, 0, 0, '', 2),
 (103, 100, '接入文档',   '/tpp/docs',       'tpp/docs',         '', 'lucide:book-open-text',  1, 0, 0, '', 3),
 -- 平台服务(PaaS 能力入口)
 (110, 0,   '平台服务',   '/paas',           '',                 '', 'lucide:boxes',           1, 0, 0, '', 40),
 (111, 110, '媒资中心',   '/paas/media',     'paas/media/index', '', 'lucide:clapperboard',    1, 0, 0, '', 1),
 (112, 110, '统一存储',   '/paas/storage',   'paas/storage/index','','lucide:hard-drive',      1, 0, 0, '', 2),
 (113, 110, '播放服务',   '/paas/play',      'paas/play/index',  '', 'lucide:play-circle',     1, 0, 0, '', 3),
 (114, 110, '支付中台',   '/paas/pay',       'paas/pay/index',   '', 'lucide:wallet',          1, 0, 0, '', 4),
 (115, 110, '广告平台',   '/paas/ad',        'paas/ad/index',    '', 'lucide:megaphone',       1, 0, 0, '', 5)
ON CONFLICT (id) DO UPDATE SET
    parent_id    = EXCLUDED.parent_id,
    name         = EXCLUDED.name,
    route_url    = EXCLUDED.route_url,
    component    = EXCLUDED.component,
    method       = EXCLUDED.method,
    icon         = EXCLUDED.icon,
    is_menu      = EXCLUDED.is_menu,
    hide_in_menu = EXCLUDED.hide_in_menu,
    affix_tab    = EXCLUDED.affix_tab,
    active_path  = EXCLUDED.active_path,
    sort         = EXCLUDED.sort,
    updated_at   = now();

SELECT setval(pg_get_serial_sequence('manage_permission', 'id'), (SELECT MAX(id) FROM manage_permission));

-- +goose Down
DELETE FROM manage_permission WHERE id BETWEEN 100 AND 115;
SELECT setval(pg_get_serial_sequence('manage_permission', 'id'), (SELECT COALESCE(MAX(id), 1) FROM manage_permission));
