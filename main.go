package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
)

type apiConfig struct {
	fileserverHits int
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerResets(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerAdmin(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html") // normal header
	w.Write([]byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileserverHits)))
}

func respondWithError(w http.ResponseWriter, code int, msg string) http.ResponseWriter {
	w.WriteHeader(code)
	errResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}
	data, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	return w
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) http.ResponseWriter {
	w.WriteHeader(code)
	data, _ := json.Marshal(payload)
	w.Write(data)

	return w
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type ChirpStorage struct {
	Chirps []Chirp
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w = respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		w = respondWithError(w, 400, "Chirp is too long")
		return
	} else {
		cleanedBody := cleanBody(params.Body)
		w = respondWithJSON(w, 200, Chirp{
			Id:   len(ChirpStorage.Chirps) + 1,
			Body: cleanedBody,
		})
	}
}

var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func cleanBody(s string) string {
	for _, badWord := range badWords {
		s = strings.ReplaceAll(s, badWord, "****")
		s = strings.ReplaceAll(s, capitalizeFirstLetter(badWord), "****")

	}
	return s
}
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func main() {
	serverMux := http.NewServeMux()
	apiCfg := new(apiConfig)

	serverMux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serverMux.Handle("/assets", http.FileServer(http.Dir("assets/")))
	serverMux.HandleFunc("/api/healthz", handler)
	serverMux.HandleFunc("/api/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("/admin/metrics", apiCfg.handlerAdmin)
	serverMux.HandleFunc("/api/reset", apiCfg.handlerResets)
	serverMux.HandleFunc("/api/chirp", apiCfg.handlerChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
