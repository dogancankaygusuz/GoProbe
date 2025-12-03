# 🚀 GoProbe - Dağıtık Web İzleme Sistemi (Distributed Uptime Monitor)

GoProbe, **Go (Golang)** dili ile geliştirilmiş, yüksek performanslı ve dağıtık mimariye sahip bir sistem izleme aracıdır. Mikroservisler arası iletişimde **gRPC**, veri tutarlılığı için **SQLite/PostgreSQL** ve eşzamanlı işlemler için **Go Concurrency** modellerini kullanır.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![gRPC](https://img.shields.io/badge/gRPC-Protobuf-244c5a?style=flat&logo=google)
![Database](https://img.shields.io/badge/SQLite-PostgreSQL-blue?style=flat&logo=postgresql)

## 🏗 Proje Mimarisi

Sistem iki ana bileşenden oluşur:
1.  **Master Node (Yönetici):** `config.json` dosyasından hedefleri okur, tarama işlemlerini yönetir ve sonuçları veritabanına kaydeder.
2.  **Worker Node (İşçi):** Master'dan gelen gRPC isteklerini karşılar, hedef sitelere HTTP istekleri atar ve analiz sonuçlarını (Süre, Durum Kodu vb.) raporlar.

## ✨ Temel Özellikler
- **Mikroservis Mimarisi:** Servisler arası iletişim hızlı ve güvenli olan gRPC (Protobuf) ile sağlanır.
- **Eşzamanlılık (Concurrency):** Binlerce siteyi aynı anda tarayabilmek için Goroutines ve WaitGroup yapısı kullanılır.
- **Veri Kalıcılığı:** Sonuçlar otomatik olarak SQLite veritabanına kaydedilir.
- **Kolay Konfigürasyon:** İzlenecek siteler JSON dosyası üzerinden yönetilebilir.
- **Docker Desteği:** İstenirse veritabanı Docker üzerinde PostgreSQL olarak çalıştırılabilir.

## 🚀 Kurulum ve Çalıştırma

### 1. Projeyi Klonlayın
git clone https://github.com/dogancankaygusuz/goprobe.git

cd goprobe

### 2. Bağımlılıkları Yükleyin
go mod tidy

### 3. Çalıştırma (Windows)
Projeyi kolayca başlatmak için run.bat dosyasını kullanabilirsiniz:

.\run.bat

Bu komut Worker ve Master servislerini ayrı terminallerde otomatik olarak başlatır.

## Alternatif olarak manuel çalıştırma:

### Terminal 1 (Worker)
go run cmd/worker/main.go

### Terminal 2 (Master)
go run cmd/master/main.go

🛠 Konfigürasyon (config.json)
İzlemek istediğiniz web sitelerini config.json dosyasını düzenleyerek ekleyebilirsiniz:

JSON
{
  "timeout": 10,
  "targets": [
    "https://www.dogancankaygusuz.com",
    "https://github.com/dogancankaygusuz",
    "https://www.linkedin.com/in/dogancan-kaygusuz",
    "https://www.google.com"
  ]
}

## 🗄 Veritabanı
Proje varsayılan olarak kurulum gerektirmeyen SQLite kullanır. Veriler proje dizinindeki goprobe.db dosyasına kaydedilir. Bu dosyayı herhangi bir "SQLite Viewer" ile görüntüleyebilirsiniz.
