FROM --platform=linux/amd64 golang:1.24-alpine as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /build/
COPY . .

RUN go mod download
RUN go build -o sight ./api/cmd/cli


FROM scratch
WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/sight ./bin/sight
COPY --from=builder /build/env.example .

CMD [ "./api/sight"]
