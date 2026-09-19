// Command server exposes a REST endpoint for sending Loring/Koyo HAP gas
// commands over UDP. It is a separate binary from the CLI in the module
// root so that the plain CLI stays free of any HTTP dependency.
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
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

func gasHandler(plcIP string, plcPort int) http.HandlerFunc {
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

		var req gasRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if req.Value < 0 || req.Value > 100 {
			http.Error(w, "value must be between 0 and 100", http.StatusBadRequest)
			return
		}

		targetIP := plcIP
		if req.IP != "" {
			if net.ParseIP(req.IP) == nil {
				http.Error(w, fmt.Sprintf("invalid ip address: %q", req.IP), http.StatusBadRequest)
				return
			}
			targetIP = req.IP
		}

		frame, err := hap.SendGas(req.Value, targetIP, plcPort)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to send gas command: %v", err), http.StatusBadGateway)
			return
		}
		log.Printf("gas=%d frame=%x", req.Value, frame)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gasResponse{Value: req.Value, Frame: fmt.Sprintf("%x", frame)})
	}
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

	log.Printf("listening on %s (PLC target %s:%d)", *listenAddr, *plcIP, *plcPort)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
