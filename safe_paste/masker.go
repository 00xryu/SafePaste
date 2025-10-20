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
		fmt.Println("Config yüklenemedi:", configPath)
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
	cfg := loadConfig()
	hostnameRegex := regexp.MustCompile(cfg.HostnamePattern)
	ipMap := make(map[string]string)
	hostnameMap := make(map[string]string)
	keywordMap := make(map[string]string)

	ipCounter, hostnameCounter, keywordCounter := 1, 1, 1

	// Keywords replace
	for _, kw := range cfg.Keywords {
		if keywordMap[kw] == "" {
			keywordMap[kw] = fmt.Sprintf("kw%d", keywordCounter)
			keywordCounter++
		}
		input = strings.ReplaceAll(input, kw, keywordMap[kw])
	}

	// IPs replace
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
			ipMap[ip] = fmt.Sprintf("ip%d", ipCounter)
			ipCounter++
		}
		input = strings.ReplaceAll(input, ip, ipMap[ip])
	}

	// Hostnames replace
	hostnames := hostnameRegex.FindAllString(input, -1)
	for _, hn := range hostnames {
		if hostnameMap[hn] == "" {
			hostnameMap[hn] = fmt.Sprintf("hostname%d", hostnameCounter)
			hostnameCounter++
		}
		input = strings.ReplaceAll(input, hn, hostnameMap[hn])
	}

	return input
}
