package safe_paste

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Keywords        []string `json:"keywords"`
	HostnamePattern string   `json:"hostname_pattern"`
	Theme           string   `json:"theme"` // "light" or "dark"
}

// MaskResult holds both the masked text and the mapping for unmasking
type MaskResult struct {
	MaskedText string
	Mapping    map[string]string // masked -> original (e.g., "ip1" -> "192.168.1.100")
}

var (
	// IPv4: matches xxx.xxx.xxx.xxx
	ipv4Regex = regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	// IPv6: matches full and compressed forms
	ipv6Regex    = regexp.MustCompile(`\b(?:[0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}\b`)
	localhostIPs = []string{"127.0.0.1", "localhost", "0.0.0.0", "::1", "::"}
)

// isValidIPv4 checks if an IPv4 address has valid octets (0-255)
func isValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// getConfigPath exe ile aynı dizinde config.json bulur (portable)
func getConfigPath() string {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		configInExe := filepath.Join(exeDir, "config.json")
		if _, err := os.Stat(configInExe); err == nil {
			return configInExe
		}
	}
	return "config.json" // fallback: mevcut dizin
}

func LoadConfig() Config {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", configPath)
		return Config{
			Keywords:        []string{},
			HostnamePattern: "\\bxy-[a-z0-9.-]+\\b",
			Theme:           "light",
		}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	if cfg.Theme == "" {
		cfg.Theme = "light"
	}
	// fmt.Printf("Loaded config from: %s\nHostname pattern: %s\nKeywords: %v\n", configPath, cfg.HostnamePattern, cfg.Keywords)
	return cfg
}

func SaveConfig(cfg Config) error {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func MaskText(input string) string {
	result := MaskTextWithMapping(input)
	return result.MaskedText
}

// MaskTextWithMapping returns both masked text and the mapping
func MaskTextWithMapping(input string) MaskResult {
	cfg := LoadConfig()
	hostnameRegex := regexp.MustCompile(cfg.HostnamePattern)
	ipMap := make(map[string]string)
	hostnameMap := make(map[string]string)
	keywordMap := make(map[string]string)
	reverseMapping := make(map[string]string) // masked -> original

	ipCounter, hostnameCounter, keywordCounter := 1, 1, 1

	// IPv4 replace (do this first to avoid conflicts)
	ipv4s := ipv4Regex.FindAllString(input, -1)
	for _, ip := range ipv4s {
		// Skip localhost IPs
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
		// Skip invalid IPs (e.g., 256.256.256.256)
		if !isValidIPv4(ip) {
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

	// IPv6 replace
	ipv6s := ipv6Regex.FindAllString(input, -1)
	for _, ip := range ipv6s {
		// Skip localhost IPv6
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
