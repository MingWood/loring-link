# loring-link

A very small Go command-line utility and server example that sends a burner value to the Loring/Koyo PLC. The built files under the dist folder can be run as executables depending on which hardware you have:
- darwin-arm64 => macOS on Apple silicon
- darwin-amd64 => macOS on intel/AMD silicon
- linux-amd64 => Linux distribution on intel/AMD silicon
- windows-amd64 => Windows OS on intel/AMD silicon

## Prerequisites

Using this program assumes you are already connected to a Loring in a way that logging already works through something like Cropster or Artisan.

This is not an officially published access protocol by Loring for public use. Please use at your own risk.

## Quick start example

Find loring-link-server file in the dist folder for your computers OS and double click to open. Then navigate to http://localhost:8080/ on your browser.

## What it does

The cli or server program accepts:

- argument 1: burner percentage as an integer
- argument 2: PLC IP address as an optional override (defaults to 192.168.1.69 which is the Loring S7 PLC IP)

It builds a HAP gas-set frame, computes the reverse-engineered checksum, and sends the payload to the PLC over UDP port 28784.

## CLI Usage
For example for a M1 Macbook from the dist/ folder in terminal to change the burner setting live on a Loring:
```bash
./loring-link-darwin-arm64 50
./loring-link-darwin-arm64 42 192.168.1.199
```
If no IP is supplied, the program uses the default address `192.168.1.69`. Loring S15, S35 use `192.168.1.199` but should be confirmed using the HMI.

The CLI can be easily called by Artisan if you'd like to control burner through the Artisan UI. Simply set up a slider with "Action" of "Call Program" and the "Command" to call the CLI program with the value of the slider. Something like:
```bash
/Users/mingwood/Downloads/loring-link/dist/loring-link-darwin-arm64 {}
```

## REST server + web UI

`cmd/server` builds a separate, standalone binary that exposes a `POST /gas` HTTP endpoint plus two embedded web UIs (`/` for the button grid, `/entry` for free-entry with a roaster selector). It shares the frame/checksum/UDP logic with the CLI via `internal/hap` but has no effect on the CLI binary above, which stays dependency-free.

For example from an M1 Macbook in the dist/ folder:

```bash
./loring-link-server-darwin-arm64
```

By default it listens on `:8080` and targets `192.168.1.69:28784`.

Then open `http://localhost:8080/` or `http://localhost:8080/entry` in a browser, or call the endpoint directly:

```bash
curl -X POST http://localhost:8080/gas -d '{"value": 50}'
curl -X POST http://localhost:8080/gas -d '{"value": 50, "ip": "192.168.1.199"}'
```

## Build standalone executables

From the repo root:

### macOS (Intel)

```bash
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-darwin-amd64 .
```

### macOS (Apple Silicon)

```bash
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/loring-link-darwin-arm64 .
```

### Linux

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-linux-amd64 .
```

### Windows

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-windows-amd64.exe .
```

### Build all targets at once

```bash
mkdir -p dist
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/loring-link-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-windows-amd64.exe .
```

### Build standalone server executables

```bash
mkdir -p dist
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-server-darwin-amd64 ./cmd/server
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/loring-link-server-darwin-arm64 ./cmd/server
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-server-linux-amd64 ./cmd/server
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/loring-link-server-windows-amd64.exe ./cmd/server
```

## Development

```bash
go test ./...
```
