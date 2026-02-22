# http-debug

A simple Go HTTP server that returns the status code specified in the URL path and dumps received headers to stdout and the response body.

## Build & Run

```bash
go build -o http-debug main.go
./http-debug
```

## Usage

Request a specific status code (e.g., 200, 404, 500):

```bash
curl -v http://localhost:8080/200
```

With custom headers:

```bash
curl -v -H "X-Debug: true" http://localhost:8080/418
```

## Docker

```bash
docker build -t http-debug .
docker run -p 8080:8080 http-debug
```
