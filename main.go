package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"loringlink/internal/hap"
)

func parseArgs(args []string) (int, string, int, error) {
	if len(args) < 2 {
		return 0, hap.DefaultPLCIP, hap.DefaultPLCPort, errors.New("missing burner percentage argument")
	}

	value, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, hap.DefaultPLCIP, hap.DefaultPLCPort, fmt.Errorf("invalid burner percentage %q: %w", args[1], err)
	}

	targetIP := hap.DefaultPLCIP
	targetPort := hap.DefaultPLCPort
	if len(args) >= 3 {
		targetIP = args[2]
	}
	return value, targetIP, targetPort, nil
}

func main() {
	value, targetIP, targetPort, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintf(os.Stderr, "usage: %s <burner_percentage> [ip_address]\n", os.Args[0])
		os.Exit(1)
	}

	frame, err := hap.SendGas(value, targetIP, targetPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("gas=%d frame=%x\n", value, frame)
}
