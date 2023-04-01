FROM golang:alpine as app-builder
WORKDIR /app
COPY . .
RUN apk add git
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/server

FROM scratch
COPY --from=app-builder /app/server /server
# the tls certificates:
COPY --from=app-builder /app/config/crt/ /config/crt
ENTRYPOINT ["/server"]
EXPOSE 5333