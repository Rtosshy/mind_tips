# 使用するGoのバージョンを指定
FROM golang:1.23

# ワーキングディレクトリを作成
WORKDIR /app

# Goモジュールファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o main ./cmd/main.go

# コンテナがリッスンするポートを指定
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]
