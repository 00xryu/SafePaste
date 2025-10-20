## 🎉 SafePaste v1.0.0-pre.1

İlk pre-release versiyonu! 🚀

### ✨ Özellikler

- 🎨 Modern GUI arayüz
- 🔒 Otomatik IP adresi maskeleme
- 🏷️ Hostname maskeleme (regex destekli)
- 🔑 Özel keyword maskeleme
- 📋 Tek tıkla kopyalama
- ⚙️ Yapılandırılabilir config.json
- 📜 Büyük dosyalar için scroll desteği
- 🚀 Portable (kurulum gerektirmez)

### 📦 İndirme

Platformunuza uygun dosyayı seçin:

| Platform | Dosya | Boyut |
|----------|-------|-------|
| Windows (x64) | `SafePaste-windows-amd64.zip` | ~XX MB |
| Windows (ARM64) | `SafePaste-windows-arm64.zip` | ~XX MB |
| Linux (x64) | `SafePaste-linux-amd64.tar.gz` | ~XX MB |
| macOS (Intel) | `SafePaste-macos-amd64.tar.gz` | ~XX MB |
| macOS (Apple Silicon) | `SafePaste-macos-arm64.tar.gz` | ~XX MB |

### 🚀 Hızlı Başlangıç

1. Platformunuza uygun dosyayı indirin
2. Arşivi çıkarın
3. `SafePaste` uygulamasını çalıştırın
4. `config.json` dosyasını düzenleyerek kuralları özelleştirin

### 📝 Kullanım

```bash
# Linux/macOS
./SafePaste-*

# Windows
SafePaste.exe
```

### ⚙️ Yapılandırma

`config.json` dosyası exe ile aynı dizinde olmalıdır:

```json
{
  "keywords": ["password", "secret", "token"],
  "hostname_pattern": "\\bxy\\d+[a-z]+\\d*prd\\b"
}
```

### 🐛 Bilinen Sorunlar

Bu bir pre-release versiyonudur. Karşılaştığınız sorunları lütfen [Issues](https://github.com/00xryu/SafePaste/issues) sayfasında bildirin.

### 🔄 Değişiklikler

- İlk release

### 🙏 Teşekkürler

SafePaste'i denediğiniz için teşekkürler! Geri bildirimlerinizi bekliyoruz.

---

**Not**: Bu bir pre-release versiyonudur. Üretim ortamında kullanmadan önce kapsamlı test yapmanız önerilir.
