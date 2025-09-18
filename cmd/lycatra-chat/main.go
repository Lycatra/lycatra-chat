package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type Health struct {
    Status string `json:"status"`
    Time   string `json:"time"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(Health{Status: "ok", Time: time.Now().Format(time.RFC3339)})
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", healthHandler)

    srv := &http.Server{ 
        Addr: ":8080",
        Handler: mux,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout: 15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout: 60 * time.Second,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    // Graceful shutdown on SIGINT/SIGTERM in future; keep minimal for skeleton
    fmt.Println("lycatra-chat server listening on :8080")

    // Block until stdin closed to keep simple in Windows terminals
    _, _ = os.Stdin.Read(make([]byte, 1))
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
}


