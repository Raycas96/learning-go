package main

import (
	"context"
	"fmt"
	api "micro-vuln-scanner/internal/httpapi"
	generator "micro-vuln-scanner/internal/simulator"
	store "micro-vuln-scanner/internal/store"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer cancel()
	vulnStore := store.NewStore(100)
gen := generator.NewGenerator(vulnStore, 5*time.Second)
done := gen.Start(ctx)

handler := api.NewHandler(vulnStore)
mux := http.NewServeMux()
mux.HandleFunc("/api/vulnerabilities", handler.GetVulnerabilities)

server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}

go func() {
    fmt.Println("Server is running on port 8080...")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        fmt.Println("server error:", err)
        cancel()
    }
}()

<-ctx.Done()

shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer shutdownCancel()
_ = server.Shutdown(shutdownCtx)

<-done
}