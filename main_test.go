package main

import (
	"testing"

	"loringlink/internal/hap"
)

func TestParseArgsDefaults(t *testing.T) {
	value, ip, port, err := parseArgs([]string{"loringlink", "42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

func TestParseArgsCustomIP(t *testing.T) {
	value, ip, _, err := parseArgs([]string{"loringlink", "10", "10.0.0.5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 10 {
		t.Fatalf("value = %d, want 10", value)
	}
	if ip != "10.0.0.5" {
		t.Fatalf("ip = %q, want 10.0.0.5", ip)
	}
}

func TestParseArgsMissingValue(t *testing.T) {
	if _, _, _, err := parseArgs([]string{"loringlink"}); err == nil {
		t.Fatal("expected error for missing burner percentage argument")
	}
}

func TestParseArgsInvalidValue(t *testing.T) {
	if _, _, _, err := parseArgs([]string{"loringlink", "not-a-number"}); err == nil {
		t.Fatal("expected error for invalid burner percentage argument")
	}
}
