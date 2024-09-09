-- last_login_at カラムの追加
ALTER TABLE users
ADD COLUMN last_login_at DATETIME NULL AFTER pass;

-- id と mail カラムにインデックスを追加
CREATE INDEX idx_users_id ON users(id);
CREATE INDEX idx_users_mail ON users(mail);
