// Package hap builds and sends Loring/Koyo HAP command UDP frames.
package hap

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	DefaultPLCIP   = "192.168.1.69"
	DefaultPLCPort = 28784
)

// checksumMasks are the per-bit XOR contributions of the low 8 bits of a
// frame's value field. They hold across every known address/opcode; only the
// base constant (the checksum for value=0) changes per target.
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

const gasBaseChecksum = uint16(0x5d61)

var gasAddress = [3]byte{0xa1, 0x00, 0x31}

// buttonAddress is the PLC's shared "momentary button pressed" write
// address. Individual buttons are distinguished only by which value gets
// written while held (0 always means released), so the base checksum below
// is the same for every button using this address.
var buttonAddress = [3]byte{0x87, 0x01, 0x33}

const buttonBaseChecksum = uint16(0xb8bc)

const (
	buttonOffValue   = 0
	coolerFanOnValue = 4
	// coolerFanPressDuration approximates the ~160-200ms click-to-release gap
	// observed between captured button-press frame pairs.
	coolerFanPressDuration = 150 * time.Millisecond

	dropOnValue = 1
	// dropHoldDuration approximates the ~1.2s press-to-release gap observed
	// in the "long press to drop" capture. A direct one-shot write to a
	// separate address (81 00 31 = 5) was also seen on the wire immediately
	// after that long hold, but sending it alone does not reliably trigger a
	// drop on real hardware — the PLC appears to require the sustained
	// button-held state instead.
	dropHoldDuration = 1200 * time.Millisecond
)

func checksumWithBase(base uint16, value int) uint16 {
	c := base
	for i := 0; i < 8; i++ {
		if (value>>i)&1 == 1 {
			c ^= checksumMasks[i]
		}
	}
	return c
}

func GasChecksum(value int) uint16 {
	return checksumWithBase(gasBaseChecksum, value)
}

func buildFrame(baseChecksum uint16, address [3]byte, value int, nonce []byte) []byte {
	if nonce == nil || len(nonce) != 2 {
		nonce = []byte{0x00, 0x00}
	}
	cksum := checksumWithBase(baseChecksum, value)
	frame := make([]byte, 19)
	copy(frame[0:3], []byte("HAP"))
	copy(frame[3:5], nonce)
	frame[5] = byte((cksum >> 8) & 0xff)
	frame[6] = byte(cksum & 0xff)
	copy(frame[7:14], []byte{0x0a, 0x00, 0x19, 0x00, 0x01, 0x20, 0x02})
	copy(frame[14:17], address[:])
	binary.LittleEndian.PutUint16(frame[17:19], uint16(value&0xffff))
	return frame
}

func BuildGasFrame(value int, nonce []byte) []byte {
	return buildFrame(gasBaseChecksum, gasAddress, value, nonce)
}

func BuildCoolerFanFrame(pressed bool, nonce []byte) []byte {
	value := buttonOffValue
	if pressed {
		value = coolerFanOnValue
	}
	return buildFrame(buttonBaseChecksum, buttonAddress, value, nonce)
}

func BuildDropButtonFrame(pressed bool, nonce []byte) []byte {
	value := buttonOffValue
	if pressed {
		value = dropOnValue
	}
	return buildFrame(buttonBaseChecksum, buttonAddress, value, nonce)
}

func dialPLC(targetIP string, targetPort int) (*net.UDPConn, error) {
	addr := net.UDPAddr{IP: net.ParseIP(targetIP), Port: targetPort}
	if addr.IP == nil {
		return nil, fmt.Errorf("invalid IP address: %q", targetIP)
	}
	return net.DialUDP("udp4", nil, &addr)
}

func SendGas(value int, targetIP string, targetPort int) ([]byte, error) {
	conn, err := dialPLC(targetIP, targetPort)
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

// SendCoolerFanButton simulates a physical cooler fan button press: it sends
// the press frame, waits briefly, then sends the release frame, matching the
// two-message click/release pattern.
func SendCoolerFanButton(targetIP string, targetPort int) (pressFrame, releaseFrame []byte, err error) {
	conn, err := dialPLC(targetIP, targetPort)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	pressFrame = BuildCoolerFanFrame(true, nil)
	if _, err = conn.Write(pressFrame); err != nil {
		return nil, nil, err
	}

	time.Sleep(coolerFanPressDuration)

	releaseFrame = BuildCoolerFanFrame(false, nil)
	if _, err = conn.Write(releaseFrame); err != nil {
		return pressFrame, nil, err
	}

	return pressFrame, releaseFrame, nil
}

// SendDrop simulates a long press-and-hold of the drop button: it sends the
// press frame, waits ~1.2s (matching the hold duration observed on the wire),
// then sends the release frame. This mirrors what the physical HMI button
// does and is what actually triggers a drop on real hardware.
func SendDrop(targetIP string, targetPort int) (pressFrame, releaseFrame []byte, err error) {
	conn, err := dialPLC(targetIP, targetPort)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	pressFrame = BuildDropButtonFrame(true, nil)
	if _, err = conn.Write(pressFrame); err != nil {
		return nil, nil, err
	}

	time.Sleep(dropHoldDuration)

	releaseFrame = BuildDropButtonFrame(false, nil)
	if _, err = conn.Write(releaseFrame); err != nil {
		return pressFrame, nil, err
	}

	return pressFrame, releaseFrame, nil
}
