FROM golang:1.22.4-bookworm AS builder

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 

WORKDIR /build/

COPY . .

RUN go mod download
RUN go build -o verifylab cmd/

FROM scratch
WORKDIR /api/
ENV PATH=/api/bin/:$PATH


COPY --from=builder /build/verifylab ./bin/verifylab
COPY --from=builder /build/env.example .

EXPOSE 9235


CMD [ "/api/verifylab", "-env", "/api/env.example" ]
