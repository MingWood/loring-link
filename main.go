package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"loringlink/internal/hap"
)

func parseArgs(args []string) (cmd string, value int, targetIP string, targetPort int, err error) {
	targetIP = hap.DefaultPLCIP
	targetPort = hap.DefaultPLCPort

	if len(args) < 2 {
		return "", 0, targetIP, targetPort, errors.New("missing command")
	}
	cmd = args[1]

	switch cmd {
	case "gas":
		if len(args) < 3 {
			return cmd, 0, targetIP, targetPort, errors.New("missing burner percentage argument")
		}
		value, err = strconv.Atoi(args[2])
		if err != nil {
			return cmd, 0, targetIP, targetPort, fmt.Errorf("invalid burner percentage %q: %w", args[2], err)
		}
		if len(args) >= 4 {
			targetIP = args[3]
		}
		return cmd, value, targetIP, targetPort, nil

	case "coolerfan", "drop":
		if len(args) >= 3 {
			targetIP = args[2]
		}
		return cmd, 0, targetIP, targetPort, nil

	default:
		return cmd, 0, targetIP, targetPort, fmt.Errorf("unknown command %q", cmd)
	}
}

func usage(progName string) {
	fmt.Fprintf(os.Stderr, "usage: %s <command> <value> [ip_address]\n", progName)
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  gas <burner_percentage> [ip_address]")
	fmt.Fprintln(os.Stderr, "  coolerfan [ip_address]")
	fmt.Fprintln(os.Stderr, "  drop [ip_address]")
}

func main() {
	cmd, value, targetIP, targetPort, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		usage(os.Args[0])
		os.Exit(1)
	}

	switch cmd {
	case "gas":
		frame, err := hap.SendGas(value, targetIP, targetPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("gas=%d frame=%x\n", value, frame)

	case "coolerfan":
		press, release, err := hap.SendCoolerFanButton(targetIP, targetPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("coolerfan press=%x release=%x\n", press, release)

	case "drop":
		press, release, err := hap.SendDrop(targetIP, targetPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("drop press=%x release=%x\n", press, release)
	}
}
