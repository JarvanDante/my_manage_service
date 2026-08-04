-- site_code 规则改为 2~32 位大写英文字母; 清掉不符合规则的历史站点数据。
-- +goose Up
COMMENT ON COLUMN site.site_code IS '全局唯一; 2~32位大写英文字母; Nacos dataId 前缀';

-- 关联表先清, 再清站点(无 FK, 按依赖顺序)
DELETE FROM site_config_publish
WHERE site_id IN (
    SELECT id FROM site WHERE site_code !~ '^[A-Z]{2,32}$'
);
DELETE FROM site_domain
WHERE site_id IN (
    SELECT id FROM site WHERE site_code !~ '^[A-Z]{2,32}$'
);
DELETE FROM site WHERE site_code !~ '^[A-Z]{2,32}$';

-- +goose Down
COMMENT ON COLUMN site.site_code IS '全局唯一; Nacos dataId 前缀 / 前台 host 识别标识';
