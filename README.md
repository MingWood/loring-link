# loring-link

A very small Go command-line utility that sends a burner value to the Loring/Koyo PLC using the same UDP frame layout as the reference Python logic.

## What it does

The program accepts:

- argument 1: burner percentage as an integer
- argument 2: PLC IP address as an optional override (defaults to 192.168.1.69)

It builds the same HAP gas-set frame as the example Python program, computes the reverse-engineered checksum, and sends the payload to the PLC over UDP port 28784.

## Usage

```bash
./loring-link 50
./loring-link 42 192.168.1.69
```

If no IP is supplied, the program uses the default address `192.168.1.69`.

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

## Development

```bash
go test ./...
```
