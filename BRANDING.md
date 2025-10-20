# SafePaste - Logo & Branding İçin Öneriler

## 🎨 Logo Fikirleri

### Konsept
- **Ana Tema**: Güvenlik + Clipboard/Yapıştır
- **Renkler**: 
  - Mavi (#00ADD8 - Go rengi)
  - Yeşil (#4CAF50 - güvenlik)
  - Turuncu (#FF9800 - uyarı/maskeleme)

### Logo Öğeleri
1. 🔒 **Kilit simgesi** + 📋 **Pano**
2. 👁️‍🗨️ **Göz (kapalı)** + 📝 **Metin**
3. 🛡️ **Kalkan** + 📄 **Belge**

## 📸 İhtiyaç Listesi

### Repository için
- [ ] `logo.png` (512x512) - Ana logo
- [ ] `logo-small.png` (128x128) - Küçük icon
- [ ] `icon.ico` (Windows için)
- [ ] `icon.icns` (macOS için)
- [ ] `screenshot.png` - Uygulama ekran görüntüsü
- [ ] `banner.png` (1280x640) - GitHub social preview

### Release için
- [ ] Release notları için markdown template

## 🖼️ Logo Oluşturma Araçları

### Online (Ücretsiz)
1. **Canva** - canva.com
   - Çok kolay, drag & drop
   - Hazır şablonlar var

2. **Figma** - figma.com
   - Profesyonel
   - Ücretsiz plan yeterli

3. **LogoMakr** - logomakr.com
   - Hızlı logo oluşturma

### AI Destekli
1. **DALL-E** / **Midjourney** - AI ile logo oluştur
   Prompt: "minimalist lock and clipboard icon, security app logo, flat design, blue and green colors"

2. **Looka** - looka.com
   - AI powered logo maker

## 📝 README için Screenshot Alma

1. Uygulamayı çalıştır
2. Örnek metin ekle
3. Windows: `Win + Shift + S` (Snipping Tool)
4. macOS: `Cmd + Shift + 4`
5. `screenshot.png` olarak kaydet

## 🏷️ Icon Oluşturma (Windows)

Bir .png'den .ico oluşturmak için:

### Online
- https://convertio.co/png-ico/
- https://icoconvert.com/

### Komut satırı (ImageMagick)
```bash
# Yükle: https://imagemagick.org/
magick convert logo.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico
```

## 📦 Asset'leri Projeye Ekleme

```
SafePaste/
├── assets/
│   ├── logo.png
│   ├── logo-small.png
│   ├── icon.ico
│   ├── icon.icns
│   └── screenshot.png
└── README.md (screenshot link)
```

## 🎯 Sonraki Adımlar

1. Logo tasarla/oluştur
2. Screenshot al
3. Icon'ları oluştur
4. `assets/` klasörüne ekle
5. README'ye screenshot ekle
6. Git commit & push
7. Tag oluştur: `git tag v1.0.0-pre.1`
8. Push: `git push origin v1.0.0-pre.1`

## 💡 Pre-release Versiyon İsimlendirme

- `v1.0.0-alpha.1` - İlk alpha
- `v1.0.0-beta.1` - Beta test
- `v1.0.0-rc.1` - Release candidate
- `v1.0.0-pre.1` - Genel pre-release
- `v1.0.0` - Resmi release

GitHub Actions otomatik olarak bunları pre-release olarak işaretleyecek!
