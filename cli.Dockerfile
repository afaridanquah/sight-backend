FROM --platform=linux/amd64 golang:1.24-alpine as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o cli ./api/cmd/cli


# FROM scratch
# WORKDIR /api/
# ENV PATH=/api/bin/:$PATH

# COPY --from=builder /build/cli ./bin/cli
# COPY --from=builder /build/env.example .

CMD [ "./cli"]
