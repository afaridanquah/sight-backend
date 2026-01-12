FROM golang:1.25-alpine

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

ENV OTEL_GO_AUTO_INCLUDE_DB_STATEMENT=true
ENV OTEL_GO_AUTO_PARSE_DB_STATEMENT=true

WORKDIR /app

RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/

RUN apt-get install -y -qq \
    tesseract-ocr-eng \
    tesseract-ocr-deu \
    tesseract-ocr-jpn

# RUN export CPATH="/opt/homebrew/include"
# RUN export LIBRARY_PATH="/opt/homebrew/lib"

RUN go install github.com/air-verse/air@latest

RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin


COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
