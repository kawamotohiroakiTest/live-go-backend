FROM golang:1.23-alpine

# 必要なパッケージをインストール (ffmpegを含む)
RUN apk add --no-cache git curl tzdata ffmpeg

# タイムゾーンをAsia/Tokyoに設定
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" > /etc/timezone

# airをインストール
RUN git clone https://github.com/cosmtrek/air /tmp/air \
    && cd /tmp/air \
    && go build -o /go/bin/air

WORKDIR /app

# Goモジュールをダウンロード
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# エントリーポイントスクリプトをコピーして実行可能にする
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# ポート8080を公開
EXPOSE 8080

# コンテナ起動時にエントリーポイントスクリプトを実行
ENTRYPOINT ["/app/entrypoint.sh"]
