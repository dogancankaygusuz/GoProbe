package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/dogancankaygusuz/goprobe/internal/grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProbeServiceServer
}

func (s *server) CheckUrl(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	log.Printf("İstek alındı, kontrol ediliyor: %s", req.GetUrl())

	start := time.Now()

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(req.GetUrl())
	elapsed := time.Since(start).Seconds() * 1000 // ms

	// Hata durumu
	if err != nil {
		return &pb.CheckResponse{
			Url:            req.GetUrl(),
			Status:         false,
			ErrorMessage:   err.Error(),
			ResponseTimeMs: elapsed,
		}, nil
	}
	defer resp.Body.Close()

	// Başarılı durum
	return &pb.CheckResponse{
		Url:            req.GetUrl(),
		StatusCode:     int32(resp.StatusCode),
		ResponseTimeMs: elapsed,
		Status:         resp.StatusCode >= 200 && resp.StatusCode < 400,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Port dinlenemiyor: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProbeServiceServer(s, &server{})

	log.Printf("Worker (İşçi) çalışıyor... Port: 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
