FROM  golang:alpine AS base
WORKDIR /go/quick-function
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /run_api ./cmd/run_api

FROM alpine:latest
COPY --from=base /run_api /app
ENTRYPOINT [ "/app" ]