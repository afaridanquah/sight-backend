FROM golang:1.25-alpine

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o guin ./api/services/parser


EXPOSE 9404

CMD [ "/api/guin", "-env", "/api/env.example" ]
