-- テーブルの作成
CREATE TABLE user_video_interactions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id INT UNSIGNED NOT NULL,  -- users テーブルの外部キー
    video_id BIGINT UNSIGNED NOT NULL,  -- videos テーブルの外部キー
    event_type ENUM('play', 'pause', 'complete', 'like', 'dislike') NOT NULL,  -- インタラクションの種類
    event_value INT UNSIGNED NULL,  -- オプションで、視聴時間や評価スコアを保存
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
);

-- インデックスの追加
CREATE INDEX idx_user_video_interactions_id ON user_video_interactions(id);
CREATE INDEX idx_user_video_interactions_user_id ON user_video_interactions(user_id);
CREATE INDEX idx_user_video_interactions_video_id ON user_video_interactions(video_id);
CREATE INDEX idx_user_video_interactions_event_type ON user_video_interactions(event_type);
CREATE INDEX idx_user_video_interactions_event_value ON user_video_interactions(event_value);
