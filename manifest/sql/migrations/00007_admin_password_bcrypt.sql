-- 管理员密码由 md5(pwd+salt) 升级为 bcrypt(盐内嵌于哈希)。
-- +goose Up
ALTER TABLE manage_admin
    ALTER COLUMN password TYPE varchar(100);

COMMENT ON COLUMN manage_admin.password IS 'bcrypt hash';
COMMENT ON COLUMN manage_admin.salt IS '已废弃(bcrypt 自带盐), 保留列兼容旧数据';

-- 种子超管 admin / admin123 → bcrypt
UPDATE manage_admin
SET password = '$2a$10$qiHGGPjc74psJ7JZTuyzX.0C6B7LYdD5SFoBDvdHhRXNIeLN/rp16',
    salt = '',
    updated_at = now()
WHERE username = 'admin';

-- +goose Down
-- bcrypt 哈希无法回退为 md5; 仅缩回字段长度(可能导致截断, 仅开发环境使用)。
ALTER TABLE manage_admin
    ALTER COLUMN password TYPE varchar(64);

COMMENT ON COLUMN manage_admin.password IS 'md5(pwd+salt)';
COMMENT ON COLUMN manage_admin.salt IS NULL;
