# SafePaste

🔒 Hassas bilgileri (IP adresleri, hostname'ler, şifreler) loglardan güvenli bir şekilde maskeleyip paylaşmanızı sağlayan portable GUI uygulaması.

![SafePaste](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux-lightgrey)

## ✨ Özellikler

- 🎨 **Modern GUI** - Gio UI ile yapılmış kullanıcı dostu arayüz
- 🔒 **Otomatik Maskeleme** - IP adresleri, hostname'ler ve özel keyword'leri maskeler
- 📋 **Kolay Paylaşım** - Maskelenmiş metni tek tıkla kopyala
- 🚀 **Portable** - Kurulum gerektirmez, USB'den bile çalıştırabilirsiniz
- ⚙️ **Özelleştirilebilir** - config.json ile kuralları kendiniz belirleyin
- 📜 **Sınırsız** - Büyük log dosyaları için scroll desteği

## 📦 İndirme

[Releases](https://github.com/00xryu/SafePaste/releases) sayfasından işletim sisteminize uygun versiyonu indirin:

- **Windows (x64)**: `SafePaste-windows-amd64.zip`
- **Linux (x64)**: `SafePaste-linux-amd64.tar.gz`

## 🚀 Kullanım

### Windows
1. ZIP dosyasını indirin ve çıkarın
2. `SafePaste.exe`'yi çalıştırın
3. Sol panele metninizi yapıştırın
4. "Maskele" butonuna basın
5. Sağ panelden maskelenmiş metni kopyalayın

### Linux/macOS
```bash
# Arşivi çıkar
tar -xzf SafePaste-*.tar.gz
cd SafePaste

# Çalıştırılabilir yap (sadece ilk seferde)
chmod +x SafePaste-*

# Çalıştır
./SafePaste-*
```

## ⚙️ Yapılandırma

`config.json` dosyasını düzenleyerek maskeleme kurallarını özelleştirebilirsiniz:

```json
{
  "keywords": ["password", "secret", "token", "api_key"],
  "hostname_pattern": "\\bxy\\d+[a-z]+\\d*prd\\b"
}
```

- **keywords**: Maskelenecek özel kelimeler (büyük/küçük harf duyarlı)
- **hostname_pattern**: Hostname'leri tanımak için regex deseni

### Örnek Kullanım

**Giriş:**
```
Bu log dosyasında password: Abc123! var. 
Sunucu xy123abc456prd adresinde çalışıyor ve 192.168.1.100 IP'sinden bağlanıyor.
```

**Çıkış:**
```
Bu log dosyasında kw1: Abc123! var. 
Sunucu hostname1 adresinde çalışıyor ve ip1 IP'sinden bağlanıyor.
```

## 🛠️ Geliştirme

### Gereksinimler
- Go 1.24+
- Gio UI bağımlılıkları

### Build
```bash
# Bağımlılıkları indir
go mod download

# Çalıştır (development)
go run .

# Build (production)
go build -ldflags="-H windowsgui -s -w" -o SafePaste.exe .
```

### GitHub Actions ile Otomatik Release

1. Kodu GitHub'a push'la
2. Tag oluştur ve push'la:
```bash
git tag v1.0.0
git push origin v1.0.0
```
3. GitHub Actions otomatik olarak tüm platformlar için build yapıp release oluşturur!

## 📝 Lisans

MIT License - Detaylar için [LICENSE](LICENSE) dosyasına bakın.

## 🤝 Katkıda Bulunma

Pull request'ler hoş karşılanır! Büyük değişiklikler için lütfen önce bir issue açın.

## 💬 Destek

Sorun mu yaşıyorsunuz? [Issue açın](https://github.com/00xryu/SafePaste/issues) veya PR gönderin!

---

⭐ Projeyi beğendiyseniz yıldız vermeyi unutmayın!
