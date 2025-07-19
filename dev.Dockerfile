FROM --platform=linux/amd64 golang:1.24-alpine

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV OTEL_GO_AUTO_INCLUDE_DB_STATEMENT=true
ENV OTEL_GO_AUTO_PARSE_DB_STATEMENT=true

WORKDIR /app

RUN go install github.com/air-verse/air@latest

RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
