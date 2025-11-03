# SafePaste

🔒 A portable GUI application that masks sensitive information (IP addresses, hostnames) from logs before sharing with AI or colleagues. Features smart unmask to restore original values after AI processing.

![SafePaste](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux-lightgrey)

## ✨ Features

- 🎨 **Modern GUI** - User-friendly 4-panel interface built with Gio UI
- 🔒 **Automatic Masking** - Masks IP addresses and hostnames automatically
- 🔄 **Smart Unmask** - Restore original values after AI processes your code
- 📋 **Easy Sharing** - Copy masked/unmasked text with one click
- 🚀 **Portable** - No installation required, run from USB
- ⚙️ **Customizable** - Define your own hostname patterns and keywords via config.json
- 📜 **Unlimited** - Scroll support for large log files

## 📦 Download

Download the appropriate version for your operating system from the [Releases](https://github.com/00xryu/SafePaste/releases) page:

- **Windows (x64)**: `SafePaste-windows-amd64.zip`
- **Linux (x64)**: `SafePaste-linux-amd64.tar.gz`

## 🚀 Usage

### Windows
1. Download and extract the ZIP file
2. Run `SafePaste.exe`
3. **Top-Left Panel**: Paste your sensitive text
4. Click **"Mask →"** to create masked version (top-right)
5. Copy masked text and share with AI/colleague
6. **Bottom-Left Panel**: Paste AI's response
7. Click **"Unmask →"** to restore original values (bottom-right)

### Linux
```bash
# Extract archive
tar -xzf SafePaste-*.tar.gz
cd SafePaste

# Make executable (first time only)
chmod +x SafePaste-*

# Run
./SafePaste-*
```

## 🔄 Workflow Example

**Step 1 - Mask sensitive data:**
- Original (top-left): `Server 192.168.1.100 connecting to xy123abc456prd`
- Masked (top-right): `Server ip1 connecting to hostname1`

**Step 2 - Get AI help:**
- Share masked version with AI
- AI suggests: `Server ip1 connecting to hostname1 (add connection timeout)`

**Step 3 - Unmask result:**
- AI Response (bottom-left): `Server ip1 connecting to hostname1 (add connection timeout)`
- Unmasked (bottom-right): `Server 192.168.1.100 connecting to xy123abc456prd (add connection timeout)`

## ⚙️ Configuration

Customize masking rules by editing the `config.json` file:

```json
{
  "keywords": [],
  "hostname_pattern": "\\bxy\\d+[a-z]+\\d*prd\\b"
}
```

- **keywords**: Custom words to mask (case-sensitive) - Add your own sensitive keywords if needed
- **hostname_pattern**: Regex pattern to identify hostnames

### Test Cases

**Test 1 - Multiple IPs:**
```
Input:  Connect from 10.0.0.5 to 192.168.1.100
Masked: Connect from ip1 to ip2
```

**Test 2 - Server Logs:**
```
Input:  Server xy123abc456prd at 192.168.1.100 is running
Masked: Server hostname1 at ip1 is running
```

**Test 3 - Network Config:**
```
Input:  Route 10.20.30.40 via xy456def789prd gateway 172.16.0.1
Masked: Route ip1 via hostname1 gateway ip2
```

## 🛠️ Development

### Requirements
- Go 1.24+
- Gio UI dependencies

### Build
```bash
# Download dependencies
go mod download

# Run (development)
go run .

# Build (production)
go build -ldflags="-H windowsgui -s -w" -o SafePaste.exe .
```

### Automated Releases with GitHub Actions

1. Push your code to GitHub
2. Create and push a tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```
3. GitHub Actions automatically builds for all platforms and creates a release!

## 📝 License

GPL-3.0 License - See [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first.

## 💬 Support

Having issues? [Open an issue](https://github.com/00xryu/SafePaste/issues) or submit a PR!

---

⭐ Star the project if you like it!
