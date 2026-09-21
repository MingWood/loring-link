// Command server exposes a REST endpoint for sending Loring/Koyo HAP gas
// commands over UDP. It is a separate binary from the CLI in the module
// root so that the plain CLI stays free of any HTTP dependency.
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"loringlink/internal/hap"
)

//go:embed static/index.html
var indexHTML string

//go:embed static/entry.html
var entryHTML string

type gasRequest struct {
	Value int    `json:"value"`
	IP    string `json:"ip,omitempty"`
}

type gasResponse struct {
	Value int    `json:"value"`
	Frame string `json:"frame"`
}

type ipOverrideRequest struct {
	IP string `json:"ip,omitempty"`
}

type coolerFanResponse struct {
	Press   string `json:"press"`
	Release string `json:"release"`
}

type dropResponse struct {
	Frame string `json:"frame"`
}

// withCORS wraps a POST-only handler with the headers needed for cross-origin
// callers (e.g. the web UI opened as a local file) and answers preflight
// OPTIONS requests directly.
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

// resolveTargetIP returns override if set (validating it), otherwise
// defaultIP.
func resolveTargetIP(defaultIP, override string) (string, error) {
	if override == "" {
		return defaultIP, nil
	}
	if net.ParseIP(override) == nil {
		return "", fmt.Errorf("invalid ip address: %q", override)
	}
	return override, nil
}

// decodeIPOverride reads an optional JSON body of the form {"ip": "..."}. An
// empty body is treated as no override.
func decodeIPOverride(r *http.Request) (string, error) {
	var req ipOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		return "", err
	}
	return req.IP, nil
}

func gasHandler(plcIP string, plcPort int) http.HandlerFunc {
	return withCORS(func(w http.ResponseWriter, r *http.Request) {
		var req gasRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if req.Value < 0 || req.Value > 100 {
			http.Error(w, "value must be between 0 and 100", http.StatusBadRequest)
			return
		}

		targetIP, err := resolveTargetIP(plcIP, req.IP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		frame, err := hap.SendGas(req.Value, targetIP, plcPort)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to send gas command: %v", err), http.StatusBadGateway)
			return
		}
		log.Printf("gas=%d frame=%x", req.Value, frame)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gasResponse{Value: req.Value, Frame: fmt.Sprintf("%x", frame)})
	})
}

func coolerFanHandler(plcIP string, plcPort int) http.HandlerFunc {
	return withCORS(func(w http.ResponseWriter, r *http.Request) {
		override, err := decodeIPOverride(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		targetIP, err := resolveTargetIP(plcIP, override)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		press, release, err := hap.SendCoolerFanButton(targetIP, plcPort)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to send cooler fan command: %v", err), http.StatusBadGateway)
			return
		}
		log.Printf("coolerfan press=%x release=%x", press, release)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(coolerFanResponse{
			Press:   fmt.Sprintf("%x", press),
			Release: fmt.Sprintf("%x", release),
		})
	})
}

func dropHandler(plcIP string, plcPort int) http.HandlerFunc {
	return withCORS(func(w http.ResponseWriter, r *http.Request) {
		override, err := decodeIPOverride(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		targetIP, err := resolveTargetIP(plcIP, override)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		frame, err := hap.SendDrop(targetIP, plcPort)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to send drop command: %v", err), http.StatusBadGateway)
			return
		}
		log.Printf("drop frame=%x", frame)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dropResponse{Frame: fmt.Sprintf("%x", frame)})
	})
}

func main() {
	listenAddr := flag.String("listen", ":8080", "address to listen on")
	plcIP := flag.String("plc-ip", hap.DefaultPLCIP, "PLC IP address")
	plcPort := flag.Int("plc-port", hap.DefaultPLCPort, "PLC UDP port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})
	http.HandleFunc("/entry", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/entry" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, entryHTML)
	})
	http.HandleFunc("/gas", gasHandler(*plcIP, *plcPort))
	http.HandleFunc("/coolerfan", coolerFanHandler(*plcIP, *plcPort))
	http.HandleFunc("/drop", dropHandler(*plcIP, *plcPort))

	log.Printf("listening on %s (PLC target %s:%d)", *listenAddr, *plcIP, *plcPort)
	fmt.Println("open browser to " + hyperlink(browserURL(*listenAddr)))
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

// browserURL turns a listen address like ":8080" or "127.0.0.1:8080" into a
// browser-friendly http://localhost:<port> URL.
func browserURL(listenAddr string) string {
	_, port, err := net.SplitHostPort(listenAddr)
	if err != nil {
		return "http://localhost" + listenAddr
	}
	return "http://localhost:" + port
}

// hyperlink wraps url in the OSC 8 terminal escape sequence so terminals that
// support clickable links (most modern ones) render it as one.
func hyperlink(url string) string {
	return "\x1b]8;;" + url + "\x1b\\" + url + "\x1b]8;;\x1b\\"
}
