FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o main .
RUN go build -o migrate ./migrate/migrate.go

FROM alpine

# 必要なパッケージをインストール
RUN apk add --no-cache tzdata

# タイムゾーンを日本時間に設定
ENV TZ=Asia/Tokyo

# タイムゾーンデータを適用
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate/migrate .
COPY config.ini .

CMD ["/app/main"]