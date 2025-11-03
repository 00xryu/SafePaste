package safe_paste

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Keywords        []string `json:"keywords"`
	HostnamePattern string   `json:"hostname_pattern"`
}

// MaskResult holds both the masked text and the mapping for unmasking
type MaskResult struct {
	MaskedText string
	Mapping    map[string]string // masked -> original (e.g., "ip1" -> "192.168.1.100")
}

var (
	ipRegex      = regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	localhostIPs = []string{"127.0.0.1", "localhost", "0.0.0.0", "::1"}
)

// getConfigPath exe ile aynı dizinde config.json bulur (portable)
func getConfigPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "config.json" // fallback: mevcut dizin
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "config.json")
}

func loadConfig() Config {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", configPath)
		return Config{
			Keywords:        []string{},
			HostnamePattern: "\\bxy\\d+[a-z]+\\d*prd\\b",
		}
	}
	var cfg Config

	json.Unmarshal(data, &cfg)
	return cfg
}

func MaskText(input string) string {
	result := MaskTextWithMapping(input)
	return result.MaskedText
}

// MaskTextWithMapping returns both masked text and the mapping
func MaskTextWithMapping(input string) MaskResult {
	cfg := loadConfig()
	hostnameRegex := regexp.MustCompile(cfg.HostnamePattern)
	ipMap := make(map[string]string)
	hostnameMap := make(map[string]string)
	keywordMap := make(map[string]string)
	reverseMapping := make(map[string]string) // masked -> original

	ipCounter, hostnameCounter, keywordCounter := 1, 1, 1

	// IPs replace (do this first to avoid conflicts)
	ips := ipRegex.FindAllString(input, -1)
	for _, ip := range ips {
		skip := false
		for _, local := range localhostIPs {
			if ip == local {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if ipMap[ip] == "" {
			masked := fmt.Sprintf("ip%d", ipCounter)
			ipMap[ip] = masked
			reverseMapping[masked] = ip
			ipCounter++
		}
		input = strings.ReplaceAll(input, ip, ipMap[ip])
	}

	// Hostnames replace
	hostnames := hostnameRegex.FindAllString(input, -1)
	for _, hn := range hostnames {
		if hostnameMap[hn] == "" {
			masked := fmt.Sprintf("hostname%d", hostnameCounter)
			hostnameMap[hn] = masked
			reverseMapping[masked] = hn
			hostnameCounter++
		}
		input = strings.ReplaceAll(input, hn, hostnameMap[hn])
	}

	// Keywords replace (do this last to catch remaining sensitive words)
	for _, kw := range cfg.Keywords {
		if keywordMap[kw] == "" {
			masked := fmt.Sprintf("kw%d", keywordCounter)
			keywordMap[kw] = masked
			reverseMapping[masked] = kw
			keywordCounter++
		}
		input = strings.ReplaceAll(input, kw, keywordMap[kw])
	}

	return MaskResult{
		MaskedText: input,
		Mapping:    reverseMapping,
	}
}

// UnmaskText replaces masked values (ip1, hostname1, kw1) with original values
func UnmaskText(maskedText string, mapping map[string]string) string {
	result := maskedText
	// Replace in reverse order: keywords, hostnames, then IPs
	// This ensures longer patterns are replaced first
	for masked, original := range mapping {
		result = strings.ReplaceAll(result, masked, original)
	}
	return result
}
