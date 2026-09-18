package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	defaultPLCIP   = "192.168.1.69"
	defaultPLCPort = 28784
)

var checksumMasks = [8]uint16{
	0x3133,
	0x6266,
	0xc4cc,
	0xa989,
	0x7303,
	0xe606,
	0xcc0d,
	0x9c1b,
}

func gasChecksum(value int) uint16 {
	c := uint16(0x5d61)
	for i := 0; i < 8; i++ {
		if (value>>i)&1 == 1 {
			c ^= checksumMasks[i]
		}
	}
	return c
}

func buildGasFrame(value int, nonce []byte) []byte {
	if nonce == nil {
		nonce = []byte{0x00, 0x00}
	}
	if len(nonce) != 2 {
		nonce = []byte{0x00, 0x00}
	}
	cksum := gasChecksum(value)
	frame := make([]byte, 19)
	copy(frame[0:3], []byte("HAP"))
	copy(frame[3:5], nonce)
	frame[5] = byte((cksum >> 8) & 0xff)
	frame[6] = byte(cksum & 0xff)
	copy(frame[7:17], []byte{0x0a, 0x00, 0x19, 0x00, 0x01, 0x20, 0x02, 0xa1, 0x00, 0x31})
	binary.LittleEndian.PutUint16(frame[17:19], uint16(value&0xffff))
	return frame
}

func sendGas(value int, targetIP string, targetPort int) error {
	addr := net.UDPAddr{IP: net.ParseIP(targetIP), Port: targetPort}
	if addr.IP == nil {
		return fmt.Errorf("invalid IP address: %q", targetIP)
	}

	conn, err := net.DialUDP("udp4", nil, &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	frame := buildGasFrame(value, nil)
	if _, err := conn.Write(frame); err != nil {
		return err
	}

	fmt.Printf("gas=%d frame=%x\n", value, frame)
	return nil
}

func parseArgs(args []string) (int, string, int, error) {
	if len(args) < 2 {
		return 0, defaultPLCIP, defaultPLCPort, errors.New("missing burner percentage argument")
	}

	value, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, defaultPLCIP, defaultPLCPort, fmt.Errorf("invalid burner percentage %q: %w", args[1], err)
	}

	targetIP := defaultPLCIP
	targetPort := defaultPLCPort
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

	if err := sendGas(value, targetIP, targetPort); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
