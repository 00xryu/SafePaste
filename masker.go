package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func loadConfig() Config {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Config load error:", err)
		return Config{
			Keywords:        []string{},
			HostnamePattern: "^xy\\d+[a-z]+\\d*prd$",
		}
	}
	var cfg Config

	json.Unmarshal(data, &cfg)
	return cfg
}

func maskText(input string) string {
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

func main() {

	input := "Bağlan 10.20.20.20'ye, xy01dbp339prd'de sorun var. Şirket_sifresi: abc123"
	fmt.Println("Orijinal:", input)
	masked := maskText(input)
	fmt.Println("Maskelenmiş:", masked)
}
