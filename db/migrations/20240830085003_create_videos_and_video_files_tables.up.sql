
-- テーブル: videos
CREATE TABLE videos (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,                       -- ユーザーID（外部キー）
    title VARCHAR(255) NOT NULL,                         -- 動画のタイトル
    description TEXT NULL,                               -- 動画の説明
    created DATETIME DEFAULT CURRENT_TIMESTAMP,          -- 作成日時
    modified DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新日時
    deleted DATETIME NULL,                               -- 削除日時（ソフトデリート用）
    FOREIGN KEY (user_id) REFERENCES users(id)           -- 外部キー制約（usersテーブル）
);


-- テーブル: video_files
CREATE TABLE video_files (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    video_id BIGINT UNSIGNED NOT NULL,                   -- videosテーブルとのリレーション用外部キー
    file_path VARCHAR(255) NOT NULL,                     -- 動画ファイルの保存パス
    thumbnail_path VARCHAR(255) NULL,                    -- サムネイル画像の保存パス
    duration INT UNSIGNED NULL,                          -- 動画の再生時間（秒単位）
    file_size BIGINT UNSIGNED NULL,                      -- 動画ファイルのサイズ（バイト単位）
    format VARCHAR(50) NOT NULL,                         -- 動画ファイルのフォーマット
    status ENUM('pending', 'processing', 'completed', 'failed') NOT NULL DEFAULT 'pending', -- 動画の処理ステータス
    created DATETIME DEFAULT CURRENT_TIMESTAMP,          -- 作成日時
    modified DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新日時
    deleted DATETIME NULL,                               -- 削除日時（ソフトデリート用）
    FOREIGN KEY (video_id) REFERENCES videos(id)         -- 外部キー制約（videosテーブル）
);
