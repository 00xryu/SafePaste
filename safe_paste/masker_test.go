package safe_paste

import (
	"testing"
)

func TestMaskTextWithMapping(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedMasked string
		checkMapping   map[string]string // masked -> original
	}{
		{
			name:           "Basic IPv4 Masking",
			input:          "Connection to 192.168.1.100 failed",
			expectedMasked: "Connection to ip1 failed",
			checkMapping: map[string]string{
				"ip1": "192.168.1.100",
			},
		},
		{
			name:           "Multiple IPv4 Masking",
			input:          "From 10.0.0.1 to 10.0.0.2",
			expectedMasked: "From ip1 to ip2",
			checkMapping: map[string]string{
				"ip1": "10.0.0.1",
				"ip2": "10.0.0.2",
			},
		},
		{
			name:           "Invalid IPv4 (Out of range)",
			input:          "IP is 256.256.256.256",
			expectedMasked: "IP is 256.256.256.256",
			checkMapping:   map[string]string{},
		},
		{
			name:           "Localhost IPv4 Ignored",
			input:          "Ping 127.0.0.1",
			expectedMasked: "Ping 127.0.0.1",
			checkMapping:   map[string]string{},
		},
		{
			name:           "Basic IPv6 Masking",
			input:          "IPv6 address 2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expectedMasked: "IPv6 address ip1",
			checkMapping: map[string]string{
				"ip1": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
		},
		{
			name:           "Hostname Masking (Default Pattern xy-*)",
			input:          "Server xy-database-01 is down",
			expectedMasked: "Server hostname1 is down",
			checkMapping: map[string]string{
				"hostname1": "xy-database-01",
			},
		},
		{
			name:           "Mixed Content",
			input:          "Check xy-backend at 192.168.1.50",
			expectedMasked: "Check hostname1 at ip1",
			checkMapping: map[string]string{
				"hostname1": "xy-backend",
				"ip1":       "192.168.1.50",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskTextWithMapping(tt.input)

			if result.MaskedText != tt.expectedMasked {
				t.Errorf("MaskedText = %v, want %v", result.MaskedText, tt.expectedMasked)
			}

			for masked, original := range tt.checkMapping {
				if gotOriginal, ok := result.Mapping[masked]; !ok || gotOriginal != original {
					t.Errorf("Mapping[%v] = %v, want %v", masked, gotOriginal, original)
				}
			}
		})
	}
}

func TestUnmaskText(t *testing.T) {
	mapping := map[string]string{
		"ip1":       "192.168.1.100",
		"hostname1": "xy-server",
	}
	input := "Error on hostname1 at ip1"
	expected := "Error on xy-server at 192.168.1.100"

	result := UnmaskText(input, mapping)
	if result != expected {
		t.Errorf("UnmaskText = %v, want %v", result, expected)
	}
}

func TestIsValidIPv4(t *testing.T) {
	tests := []struct {
		ip    string
		valid bool
	}{
		{"192.168.1.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"256.0.0.0", false},
		{"1.2.3", false},
		{"1.2.3.4.5", false},
		{"a.b.c.d", false},
		{"-1.0.0.0", false},
	}

	for _, tt := range tests {
		if got := isValidIPv4(tt.ip); got != tt.valid {
			t.Errorf("isValidIPv4(%q) = %v, want %v", tt.ip, got, tt.valid)
		}
	}
}

func TestAIWorkflow(t *testing.T) {
	input := "Fix the bug in 192.168.1.100 server"

	// 1. Mask
	result := MaskTextWithMapping(input)
	// Expected: "Fix the bug in ip1 server"

	// 2. Simulate AI response (AI keeps the tokens but changes text)
	aiResponse := "I have fixed the issue on ip1 server by restarting the service."

	// 3. Unmask
	unmaskedResponse := UnmaskText(aiResponse, result.Mapping)

	expected := "I have fixed the issue on 192.168.1.100 server by restarting the service."

	if unmaskedResponse != expected {
		t.Errorf("AI Workflow failed.\nGot: %s\nWant: %s", unmaskedResponse, expected)
	}
}
