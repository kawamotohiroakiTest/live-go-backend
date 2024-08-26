FROM golang:1.22-alpine

RUN apk add --no-cache git curl

RUN git clone https://github.com/cosmtrek/air /tmp/air \
    && cd /tmp/air \
    && go build -o /go/bin/air

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Entrypoint scriptを追加して実行権限を付与
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Airを使ってアプリケーションを起動
CMD ["/app/entrypoint.sh"]

EXPOSE 8080
