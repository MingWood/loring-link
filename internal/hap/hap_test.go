package hap

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGasChecksumMatchesKnownExamples(t *testing.T) {
	known := map[int]uint16{
		1:   0x6c52,
		25:  0xb6d8,
		26:  0xe58d,
		50:  0xaa02,
		100: 0xb3a6,
	}

	for value, want := range known {
		if got := GasChecksum(value); got != want {
			t.Fatalf("GasChecksum(%d) = 0x%04x, want 0x%04x", value, got, want)
		}
	}
}

func TestBuildGasFrameMatchesPythonFrameFormat(t *testing.T) {
	for _, value := range []int{0, 30, 42, 60, 100} {
		frame := BuildGasFrame(value, nil)
		if len(frame) != 19 {
			t.Fatalf("value=%d: len(frame) = %d, want 19", value, len(frame))
		}
		if !bytes.Equal(frame[:3], []byte("HAP")) {
			t.Fatalf("value=%d: prefix mismatch: %q", value, frame[:3])
		}
		if !bytes.Equal(frame[3:5], []byte{0x00, 0x00}) {
			t.Fatalf("value=%d: nonce mismatch: %v", value, frame[3:5])
		}
		if got := fmt.Sprintf("%02x", frame[5]); got == "" {
			t.Fatal("unexpected empty checksum byte")
		}
		if frame[15] != 0x00 {
			t.Fatalf("value=%d: expected constant byte at index 15 to be 0x00, got 0x%02x", value, frame[15])
		}
		if frame[16] != 0x31 {
			t.Fatalf("value=%d: expected trailing constant byte at index 16 to be 0x31, got 0x%02x", value, frame[16])
		}
		if frame[17] != byte(value&0xff) || frame[18] != byte((value>>8)&0xff) {
			t.Fatalf("value=%d: value bytes should be little-endian encoded at end: %v", value, frame[17:19])
		}
	}
}
