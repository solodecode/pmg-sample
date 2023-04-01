FROM golang:alpine as app-builder
WORKDIR /app
COPY . .
RUN apk add git
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/client

FROM scratch
COPY --from=app-builder /app/client /client
# the tls certificates:
COPY --from=app-builder /app/config/crt/ /config/crt
ENTRYPOINT ["/client"]