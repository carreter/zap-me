# Zap Me! A Pavlok webserver.
This repo provides a small webserver to interface with the [Pavlok](https://pavlok.com/) API.

Built on top of [pavlok-go](https://github.com/carreter/pavlok-go).

## Installing and runnning

Clone the repo and run one of the following in your repo root:

### Run with Docker
```shell
docker build . -t zap-me
docker run -e AUTH_CODE="your passcode here" -e PAVLOK_API_KEY="your API key here" -t zap-me
```

### Run with go
```shell
go mod download
go build -o ./zap-me
AUTH_CODE="your passcode here" PAVLOK_API_KEY="your API key here" ./zap-me
```

## Adding HTTPS to the service
By default, the `zap-me` executable exposes an HTTP (*not HTTPS*) service on port 9000.
This means the passcode is sent over plaintext. To avoid this, you can set up a
[Caddy reverse proxy](https://caddyserver.com/docs/quick-starts/reverse-proxy) with the following:
```shell
sudo caddy reverse-proxy --from localhost --to :9000
```