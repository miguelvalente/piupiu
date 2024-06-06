package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	fmt.Println("lol", cfg.fileserverHits)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerResets(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		fmt.Println("midle", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func main() {
	serverMux := http.NewServeMux()
	// apiCfg := apiConfig{}
	apiCfg := new(apiConfig) // automatically initializes fileserverHits to 0

	serverMux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serverMux.Handle("/assets", http.FileServer(http.Dir("assets/")))
	serverMux.HandleFunc("/healthz", handler)
	serverMux.HandleFunc("/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("/reset", apiCfg.handlerResets)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()
}
