package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	pb "github.com/dogancankaygusuz/goprobe/internal/grpc/proto"
	"github.com/dogancankaygusuz/goprobe/pkg/config"
	"github.com/dogancankaygusuz/goprobe/pkg/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	latestResults = make(map[string]database.CheckResult)
	mu            sync.Mutex
)

// Web Sunucusu Handler'ı
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var results []database.CheckResult
	for _, result := range latestResults {
		results = append(results, result)
	}
	mu.Unlock()

	// HTML şablonunu yükle ve verileri gönder
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "HTML dosyası yüklenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, results)
}

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Config hatası: %v", err)
	}

	db := database.InitDB()

	// WEB SERVER BAŞLAT
	go func() {
		http.HandleFunc("/", dashboardHandler)
		log.Println("🌍 Dashboard Yayında: http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Web sunucusu hatası: %v", err)
		}
	}()

	// Worker'a Bağlan
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Worker'a bağlanılamadı: %v", err)
	}
	defer conn.Close()
	client := pb.NewProbeServiceClient(conn)

	// Ana Tarama Döngüsü
	for {
		log.Println("----- Tarama Başlıyor -----")
		var wg sync.WaitGroup

		for _, url := range cfg.Targets {
			wg.Add(1)

			go func(targetUrl string) {
				defer wg.Done()

				timeout := time.Duration(cfg.Timeout) * time.Second
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				resp, err := client.CheckUrl(ctx, &pb.CheckRequest{Url: targetUrl})
				result := database.CheckResult{Url: targetUrl}

				if err != nil {
					log.Printf("❌ %s Hata: %v", targetUrl, err)
					result.Status = false
					result.ErrorMessage = err.Error()
				} else {
					result.Url = resp.Url
					result.StatusCode = resp.StatusCode
					result.ResponseTimeMs = resp.ResponseTimeMs
					result.Status = resp.Status
					log.Printf("✅ %s | %d | %.0fms", resp.Url, resp.StatusCode, resp.ResponseTimeMs)
				}
				db.Create(&result)
				mu.Lock()
				latestResults[targetUrl] = result
				mu.Unlock()

			}(url)
		}
		wg.Wait()
		time.Sleep(5 * time.Second)
	}
}
