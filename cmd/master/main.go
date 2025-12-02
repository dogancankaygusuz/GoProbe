package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/dogancankaygusuz/goprobe/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Worker'a bağlanılamadı: %v", err)
	}
	defer conn.Close()

	client := pb.NewProbeServiceClient(conn)

	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://go.dev",
		"https://api.boredapi.com/api/activity", // Yavaş/Kapalı site
	}

	for {
		log.Println("----- Taramayı Başlat (Concurrent) -----")
		startTotal := time.Now()

		// WaitGroup: Tüm goroutine'lerin bitmesini beklemek için sayaç
		var wg sync.WaitGroup

		for _, url := range urls {
			wg.Add(1) // Sayacı 1 artır

			// Her URL için ayrı bir Goroutine (iş parçacığı) başlatıyoruz
			go func(targetUrl string) {
				defer wg.Done() // İş bitince sayacı 1 azalt

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				resp, err := client.CheckUrl(ctx, &pb.CheckRequest{Url: targetUrl})

				if err != nil {
					// Hata durumunda sadece log basıyoruz (ileride DB'ye yazacağız)
					log.Printf("❌ HATA [%s]: %v", targetUrl, err)
				} else {
					statusIcon := "✅"
					if !resp.Status {
						statusIcon = "🔻"
					}
					log.Printf("%s Site: %s | Kod: %d | Süre: %.0fms",
						statusIcon, resp.Url, resp.StatusCode, resp.ResponseTimeMs)
				}
			}(url)
		}

		// Tüm goroutine'ler bitene kadar burada bekle
		wg.Wait()

		totalDuration := time.Since(startTotal)
		log.Printf("----- Tarama Bitti (Toplam Süre: %v) -----\n", totalDuration)

		time.Sleep(5 * time.Second)
	}
}
