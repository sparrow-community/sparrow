package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparrow-community/sparrow/gateway"
	"github.com/sparrow-community/sparrow/processing"
)

func main() {
	dataDir := flag.String("data-dir", "data", "engine data directory (events.log + deployments)")
	listen := flag.String("listen", ":50051", "gRPC listen address")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		log.Fatalf("mkdir data-dir: %v", err)
	}
	eng, err := processing.Open(ctx, *dataDir)
	if err != nil {
		log.Fatalf("open engine: %v", err)
	}
	defer func() {
		if err := eng.Close(); err != nil {
			log.Printf("close engine: %v", err)
		}
	}()

	go fireDueLoop(ctx, eng)

	lis, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("sparrow listening on %s (data-dir=%s)", lis.Addr(), *dataDir)

	srv := gateway.NewServer(eng)
	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func fireDueLoop(ctx context.Context, eng *processing.Engine) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := eng.FireDue(ctx); err != nil && ctx.Err() == nil {
				log.Printf("fire due: %v", err)
			}
		}
	}
}
