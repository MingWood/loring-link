# loring-link

A very small Go command-line utility and server example that sends burner, cooler fan, and drop commands to a Loring/Koyo PLC. The built files under the dist folder can be run as executables depending on which hardware you have:
darwin-arm64 => macOS on Apple silicon
darwin-amd64 => macOS on intel/AMD silicon
linux-amd64 => Linux distribution on intel/AMD silicon
windows-amd64 => Windows OS on intel/AMD silicon

## Prerequisites

Using this program assumes you are already connected to a Loring in a way that logging already works through something like Cropster or Artisan.

This is NOT an officially published access protocol by Loring for public use. Please use at your own risk.

## Quick start example

Download or clone this repo.

Find loring-link-server file in the dist folder for your computers OS and double click to open or call the file in terminal `./loring-link-server-darwin-arm64`. Then navigate to http://localhost:8080/ on your browser. 

Depending on your computers security you may need to change some permissions to allow this file to run including running `sudo chmod 755 <filname_here>` if on a Mac and you didn't use git to clone the repo.

## What it does

The CLI and REST server both talk to the Loring/Koyo PLC over UDP port 28784 using the reverse-engineered HAP frame format. Three commands are supported:

- `gas <percentage>` — sets the burner/gas value (0-100)
- `coolerfan` — simulates a cooler fan button press: sends a press frame, waits briefly, then sends a release frame, matching the physical button's click/release behavior
- `drop` — simulates long press on drop button. On the real HMI this only fires after the button has been held down for a sustained period of 1.2s; the CLI and server send it immediately when called, so it's up to the caller (or the web UI) to gate on a press-and-hold gesture if you want to avoid accidental drops

All three optionally accept a PLC IP address override. If none is supplied, `192.168.1.69` (the Loring S7 PLC IP) is used. Loring S15 and S35 use `192.168.1.199`, but confirm this using the HMI.

## CLI Usage

For example for a M1 Macbook from the dist/ folder in terminal:

```bash
./loring-link-darwin-arm64 gas 50
./loring-link-darwin-arm64 gas 42 192.168.1.199
./loring-link-darwin-arm64 coolerfan
./loring-link-darwin-arm64 coolerfan 192.168.1.199
./loring-link-darwin-arm64 drop
./loring-link-darwin-arm64 drop 192.168.1.199
```

The CLI can be easily called by Artisan if you'd like to control the burner through the Artisan UI. Simply set up a slider with "Action" of "Call Program" and the "Command" to call the CLI program with the `gas` command and the value of the slider. Something like:

```bash
/Users/mingwood/Downloads/loring-link/dist/loring-link-darwin-arm64 gas {}
```

## REST server + web UI

`cmd/server` builds a separate, standalone binary that exposes `POST /gas`, `POST /coolerfan`, and `POST /drop` HTTP endpoints, plus two embedded web UIs:

- `/` — a button grid for common burner values (20-100), an "I'm feeling lucky" button for a random value, Cooler Fan and Drop (hold 1.5s) buttons, and a Loring selector (S7/S15/S35) that automatically applies the `192.168.1.199` IP override for S15/S35
- `/entry` — a free-entry numeric input (clamped to 20-100, press Enter to send) with the same roaster selector and Cooler Fan/Drop controls, plus a green/red border flash showing whether the last gas request succeeded

The server shares the frame/checksum/UDP logic with the CLI via `internal/hap` but has no effect on the CLI binary above, which stays dependency-free.

For example from an M1 Macbook in the dist/ folder:

```bash
./loring-link-server-darwin-arm64
```

By default it listens on `:8080` and targets `192.168.1.69:28784`.

Then open `http://localhost:8080/` or `http://localhost:8080/entry` in a browser, or call the endpoints directly. Each accepts an optional `"ip"` field to override the default PLC address for that request:

```bash
curl -X POST http://localhost:8080/gas -d '{"value": 50}'
curl -X POST http://localhost:8080/gas -d '{"value": 50, "ip": "192.168.1.199"}'

curl -X POST http://localhost:8080/coolerfan -d '{}'
curl -X POST http://localhost:8080/coolerfan -d '{"ip": "192.168.1.199"}'

curl -X POST http://localhost:8080/drop -d '{}'
curl -X POST http://localhost:8080/drop -d '{"ip": "192.168.1.199"}'
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
