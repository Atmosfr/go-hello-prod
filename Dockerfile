FROM golang:1.25.5-alpine3.23@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w -extldflags=-static" -o /hello-service . && chmod +x /hello-service

# Final stage
FROM cgr.dev/chainguard/wolfi-base@sha256:ca099f560465d417619fc46f4f14719e791858c36a4a76559a6cd6d5546fec1f

COPY --from=builder /hello-service /hello-service

USER 65532

EXPOSE 8080

ENTRYPOINT [ "/hello-service" ]

