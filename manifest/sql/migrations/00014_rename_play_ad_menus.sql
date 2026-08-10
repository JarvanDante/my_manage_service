-- 菜单文案：播放服务 → 统一播放；广告平台 → 广告中台
-- +goose Up

UPDATE manage_permission
SET name = '统一播放', updated_at = now()
WHERE id = 113;

UPDATE manage_permission
SET name = '广告中台', updated_at = now()
WHERE id = 115;

-- +goose Down

UPDATE manage_permission
SET name = '播放服务', updated_at = now()
WHERE id = 113;

UPDATE manage_permission
SET name = '广告平台', updated_at = now()
WHERE id = 115;
