# docker build -t orchard/orchard-gateway-go:latest .
FROM golang:1.11.4-alpine as builder
RUN apk update --no-cache && apk add git
RUN mkdir /staging
WORKDIR /staging
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags orchard-gateway-go -o build/gateway-msvc -ldflags '-w' ./cmd/gateway-msvc/

FROM scratch
COPY --from=builder staging/build/gateway-msvc app
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/app"]