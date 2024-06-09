package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type apiConfig struct {
	fileserverHits int
	DB             DB
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

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
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
			chirp, err := cfg.DB.CreateChirp(cleanedBody)
			if err != nil {
				w = respondWithError(w, 500, "Something went wrong making chirps")
				return
			}
			w = respondWithJSON(w, 201, chirp)
		}
	} else if req.Method == http.MethodGet {

		chirpId, err := strconv.Atoi(req.PathValue("chirpId"))
		if err != nil {
			chirps, _ := cfg.DB.GetChirps()
			w = respondWithJSON(w, 200, chirps)
		} else {
			chirp, err := cfg.DB.GetChirp(chirpId)
			if err != nil {
				w = respondWithError(w, 404, "Chirp Id does not exist")
			}
			w = respondWithJSON(w, 200, chirp)
		}
	}
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			w = respondWithError(w, 500, "Something went wrong")
			return
		}

		if cfg.DB.UserExists(params.Email) {
			w = respondWithError(w, 403, "User already Exists")
			return

		}

		hashed_password, err := hash(params.Password)
		if err != nil {
			w = respondWithError(w, 500, "Error hashing password")
		}

		user, err := cfg.DB.CreateUser(params.Email, hashed_password)
		if err != nil {
			w = respondWithError(w, 500, "Something went wrong making chirps")
			return
		}

		w = respondWithJSON(w, 201, user)
	} else if req.Method == http.MethodPut {
		tokenString := req.Header.Get("Authorization")
		if tokenString == "" {
			w = respondWithError(w, 500, "No Authorization Token")
			return
		}

		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			w = respondWithError(w, 500, err.Error())
		}

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			w = respondWithError(w, 500, "Something went wrong")
			return
		}

		hashedPassword, err := hash(params.Password)
		if err != nil {
			w = respondWithError(w, 500, err.Error())
		}

		userOut, err := cfg.DB.UpdateUser(userId, params.Email, hashedPassword)
		if err != nil {
			w = respondWithError(w, 500, err.Error())
		}

		respondWithJSON(w, 200, userOut)

	}

}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Expires  int    `json:"expires_in_seconds,omitempty"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			w = respondWithError(w, 500, "Something went wrong")
			return
		}

		user, err := cfg.DB.GetUserByEmail(params.Email)

		if err != nil {
			w = respondWithError(w, 500, "User not found")
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
		if err != nil {
			w = respondWithJSON(w, 401, "Wrong credentials")
		}

		token := createJWT(user, params.Expires)

		userOut := UserOutLogin{
			Email: user.Email,
			Token: token,
		}
		w = respondWithJSON(w, 200, userOut)

	} else {
		w = respondWithJSON(w, 405, "Method not allowed")
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
	godotenv.Load()
	serverMux := http.NewServeMux()

	err := os.Remove("database.json")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to delete existing database file: %v", err)
	} else {
		log.Println("Existing database file deleted or not found.")
	}
	db_, _ := NewDB("database.json")
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             *db_,
	}

	serverMux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serverMux.Handle("/assets", http.FileServer(http.Dir("assets/")))
	serverMux.HandleFunc("/api/healthz", handler)
	serverMux.HandleFunc("/api/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("/admin/metrics", apiCfg.handlerAdmin)
	serverMux.HandleFunc("/api/reset", apiCfg.handlerResets)
	serverMux.HandleFunc("/api/chirps", apiCfg.handlerChirp)
	serverMux.HandleFunc("/api/chirps/{chirpId}", apiCfg.handlerChirp)
	serverMux.HandleFunc("/api/users", apiCfg.handlerUser)
	serverMux.HandleFunc("/api/login", apiCfg.handlerLogin)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
