# SafePaste

🔒 A portable GUI application that safely masks sensitive information (IP addresses, hostnames, passwords) from logs before sharing.

![SafePaste](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux-lightgrey)

## ✨ Features

- 🎨 **Modern GUI** - User-friendly interface built with Gio UI
- 🔒 **Automatic Masking** - Masks IP addresses, hostnames, and custom keywords
- 📋 **Easy Sharing** - Copy masked text with one click
- 🚀 **Portable** - No installation required, run from USB
- ⚙️ **Customizable** - Define your own rules via config.json
- 📜 **Unlimited** - Scroll support for large log files

## 📦 Download

Download the appropriate version for your operating system from the [Releases](https://github.com/00xryu/SafePaste/releases) page:

- **Windows (x64)**: `SafePaste-windows-amd64.zip`
- **Linux (x64)**: `SafePaste-linux-amd64.tar.gz`

## 🚀 Usage

### Windows
1. Download and extract the ZIP file
2. Run `SafePaste.exe`
3. Paste your text in the left panel
4. Click the "Mask" button
5. Copy the masked text from the right panel

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

## ⚙️ Configuration

Customize masking rules by editing the `config.json` file:

```json
{
  "keywords": ["password", "secret", "token", "api_key"],
  "hostname_pattern": "\\bxy\\d+[a-z]+\\d*prd\\b"
}
```

- **keywords**: Custom words to mask (case-sensitive)
- **hostname_pattern**: Regex pattern to identify hostnames

### Example Usage

**Input:**
```
This log file contains password: Abc123!
Server xy123abc456prd is running and connecting from 192.168.1.100.
```

**Output:**
```
This log file contains kw1: Abc123!
Server hostname1 is running and connecting from ip1.
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
