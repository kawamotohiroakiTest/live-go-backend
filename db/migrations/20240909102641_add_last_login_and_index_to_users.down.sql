-- last_login_at カラムの削除
ALTER TABLE users
DROP COLUMN last_login_at;

-- インデックスの削除
DROP INDEX idx_users_id ON users;
DROP INDEX idx_users_mail ON users;
