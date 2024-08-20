FROM golang:1.22-alpine

RUN apk add --no-cache git curl

WORKDIR /app

RUN git clone https://github.com/air-verse/air.git /tmp/air \
    && cd /tmp/air \
    && go build -o /go/bin/air

COPY go.mod ./
RUN go mod download

COPY . .

CMD ["air"]

EXPOSE 8080
