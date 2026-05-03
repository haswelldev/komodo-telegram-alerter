package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sshaddicts/komodo-telegram-alerter/internal/render"
	"github.com/sshaddicts/komodo-telegram-alerter/internal/telegram"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /alert", handleAlert)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("[INFO] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[WARN] shutdown: %v", err)
	}
	log.Println("[INFO] stopped")
}

func handleAlert(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	chatID := r.URL.Query().Get("chat_id")
	if token == "" || chatID == "" {
		http.Error(w, `{"error":"missing token or chat_id"}`, http.StatusBadRequest)
		return
	}

	var alert render.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		log.Printf("[ERROR] decode body: %v", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	msg := render.Render(alert)

	if err := telegram.Send(token, chatID, msg); err != nil {
		log.Printf("[ERROR] send telegram: %v", err)
		http.Error(w, `{"error":"failed to send message"}`, http.StatusBadGateway)
		return
	}

	log.Printf("[INFO] sent alert type=%s level=%s resolved=%v", alert.Data.Type, alert.Level, alert.Resolved)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"success":true}`))
}
