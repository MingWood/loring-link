package main

import (
	"testing"

	"loringlink/internal/hap"
)

func TestParseArgsGasDefaults(t *testing.T) {
	cmd, value, ip, port, err := parseArgs([]string{"loringlink", "gas", "42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "gas" {
		t.Fatalf("cmd = %q, want gas", cmd)
	}
	if value != 42 {
		t.Fatalf("value = %d, want 42", value)
	}
	if ip != hap.DefaultPLCIP {
		t.Fatalf("ip = %q, want %q", ip, hap.DefaultPLCIP)
	}
	if port != hap.DefaultPLCPort {
		t.Fatalf("port = %d, want %d", port, hap.DefaultPLCPort)
	}
}

func TestParseArgsGasCustomIP(t *testing.T) {
	cmd, value, ip, _, err := parseArgs([]string{"loringlink", "gas", "10", "10.0.0.5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "gas" {
		t.Fatalf("cmd = %q, want gas", cmd)
	}
	if value != 10 {
		t.Fatalf("value = %d, want 10", value)
	}
	if ip != "10.0.0.5" {
		t.Fatalf("ip = %q, want 10.0.0.5", ip)
	}
}

func TestParseArgsGasMissingValue(t *testing.T) {
	if _, _, _, _, err := parseArgs([]string{"loringlink", "gas"}); err == nil {
		t.Fatal("expected error for missing burner percentage argument")
	}
}

func TestParseArgsGasInvalidValue(t *testing.T) {
	if _, _, _, _, err := parseArgs([]string{"loringlink", "gas", "not-a-number"}); err == nil {
		t.Fatal("expected error for invalid burner percentage argument")
	}
}

func TestParseArgsCoolerFanDefaults(t *testing.T) {
	cmd, _, ip, port, err := parseArgs([]string{"loringlink", "coolerfan"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "coolerfan" {
		t.Fatalf("cmd = %q, want coolerfan", cmd)
	}
	if ip != hap.DefaultPLCIP {
		t.Fatalf("ip = %q, want %q", ip, hap.DefaultPLCIP)
	}
	if port != hap.DefaultPLCPort {
		t.Fatalf("port = %d, want %d", port, hap.DefaultPLCPort)
	}
}

func TestParseArgsCoolerFanCustomIP(t *testing.T) {
	cmd, _, ip, _, err := parseArgs([]string{"loringlink", "coolerfan", "10.0.0.5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "coolerfan" {
		t.Fatalf("cmd = %q, want coolerfan", cmd)
	}
	if ip != "10.0.0.5" {
		t.Fatalf("ip = %q, want 10.0.0.5", ip)
	}
}

func TestParseArgsDropDefaults(t *testing.T) {
	cmd, _, ip, port, err := parseArgs([]string{"loringlink", "drop"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "drop" {
		t.Fatalf("cmd = %q, want drop", cmd)
	}
	if ip != hap.DefaultPLCIP {
		t.Fatalf("ip = %q, want %q", ip, hap.DefaultPLCIP)
	}
	if port != hap.DefaultPLCPort {
		t.Fatalf("port = %d, want %d", port, hap.DefaultPLCPort)
	}
}

func TestParseArgsMissingCommand(t *testing.T) {
	if _, _, _, _, err := parseArgs([]string{"loringlink"}); err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestParseArgsUnknownCommand(t *testing.T) {
	if _, _, _, _, err := parseArgs([]string{"loringlink", "bogus"}); err == nil {
		t.Fatal("expected error for unknown command")
	}
}
