FROM golang:1.22-alpine

RUN apk add --no-cache git curl

RUN git clone https://github.com/cosmtrek/air /tmp/air \
    && cd /tmp/air \
    && go build -o /go/bin/air

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]

EXPOSE 8080
