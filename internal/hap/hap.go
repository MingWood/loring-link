// Package hap builds and sends Loring/Koyo HAP gas-command UDP frames.
package hap

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	DefaultPLCIP   = "192.168.1.69"
	DefaultPLCPort = 28784
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

func GasChecksum(value int) uint16 {
	c := uint16(0x5d61)
	for i := 0; i < 8; i++ {
		if (value>>i)&1 == 1 {
			c ^= checksumMasks[i]
		}
	}
	return c
}

func BuildGasFrame(value int, nonce []byte) []byte {
	if nonce == nil {
		nonce = []byte{0x00, 0x00}
	}
	if len(nonce) != 2 {
		nonce = []byte{0x00, 0x00}
	}
	cksum := GasChecksum(value)
	frame := make([]byte, 19)
	copy(frame[0:3], []byte("HAP"))
	copy(frame[3:5], nonce)
	frame[5] = byte((cksum >> 8) & 0xff)
	frame[6] = byte(cksum & 0xff)
	copy(frame[7:17], []byte{0x0a, 0x00, 0x19, 0x00, 0x01, 0x20, 0x02, 0xa1, 0x00, 0x31})
	binary.LittleEndian.PutUint16(frame[17:19], uint16(value&0xffff))
	return frame
}

// SendGas builds the frame for value and sends it to targetIP:targetPort over
// UDP, returning the frame that was sent.
func SendGas(value int, targetIP string, targetPort int) ([]byte, error) {
	addr := net.UDPAddr{IP: net.ParseIP(targetIP), Port: targetPort}
	if addr.IP == nil {
		return nil, fmt.Errorf("invalid IP address: %q", targetIP)
	}

	conn, err := net.DialUDP("udp4", nil, &addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	frame := BuildGasFrame(value, nil)
	if _, err := conn.Write(frame); err != nil {
		return nil, err
	}

	return frame, nil
}
