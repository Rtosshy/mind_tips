# 使用するGoのバージョンを指定
FROM golang:1.23

# ワーキングディレクトリを作成
WORKDIR /app

RUN go install github.com/air-verse/air@latest

# Goモジュールファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# コンテナがリッスンするポートを指定
EXPOSE 8080

# アプリケーションを実行
CMD ["air"]
