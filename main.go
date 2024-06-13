package main

import (
	"crypto/rand"
	"encoding/hex"
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

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		subject, w := authenticateUser(w, req)

		if subject == "" {
			return
		}
		userId, err := strconv.Atoi(subject)
		if err != nil {
			w = respondWithError(w, 500, err.Error())
		}

		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			w = respondWithError(w, 500, "Something went wrong")
			return
		}

		if len(params.Body) > 140 {
			w = respondWithError(w, 400, "Chirp is too long")
			return
		} else {
			cleanedBody := cleanBody(params.Body)
			chirp, err := cfg.DB.CreateChirp(cleanedBody, userId)
			if err != nil {
				w = respondWithError(w, 500, "Something went wrong making chirps")
				return
			}
			w = respondWithJSON(w, 201, chirp)
		}
	} else if req.Method == http.MethodGet {
		s := req.URL.Query().Get("author_id")
		if s == "" {
			chirpId, err := strconv.Atoi(req.PathValue("chirpId"))
			if err != nil {
				chirps, _ := cfg.DB.GetChirps()
				w = respondWithJSON(w, 200, chirps)
				return
			}

			chirp, err := cfg.DB.GetChirp(chirpId)
			if err != nil {
				w = respondWithError(w, 404, "Chirp Id does not exist")
			}
			w = respondWithJSON(w, 200, chirp)

		} else {
			authorId, err := strconv.Atoi(s)
			if err != nil {
				w.WriteHeader(500)
			}
			sortMethod := req.URL.Query().Get("sort")
			if sortMethod == "" {
				sortMethod == "asc"
			}

			chirps, _ := cfg.DB.GetUserChirps(authorId, sortMethod)
			w = respondWithJSON(w, 200, chirps)
			return

		}

	} else if req.Method == http.MethodDelete {
		subject, w := authenticateUser(w, req)

		if subject == "" {
			w.WriteHeader(204)
			return
		}

		userId, err := strconv.Atoi(subject)
		if err != nil {
			w = respondWithError(w, 500, err.Error())
		}

		chirpId, err := strconv.Atoi(req.PathValue("chirpId"))
		if err != nil {
			chirps, _ := cfg.DB.GetChirps()
			w = respondWithJSON(w, 200, chirps)
			return
		}

		err = cfg.DB.DeleteChirp(chirpId, userId)
		if err != nil {
			w = respondWithError(w, 403, err.Error())
		}

		w.WriteHeader(204)

	}
}

func authenticateUser(w http.ResponseWriter, req *http.Request) (string, http.ResponseWriter) {

	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		w = respondWithError(w, 500, "No Authorization Token")
		return "", w
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return "", w
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return "", w
	}

	return claims.Subject, w
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, req *http.Request) {
	// helperPrintHeaders(req)
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
		subject, w := authenticateUser(w, req)
		if subject == "" {
			return
		}

		userId, err := strconv.Atoi(subject)
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

func makeRefreshToken() (string, error) {
	c := 128
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}

	refreshToken := hex.EncodeToString(b)
	return refreshToken, nil
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	// helperPrintHeaders(req)
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

		token := createJWT(user)

		refreshToken, err := makeRefreshToken()
		if err != nil {
			w = respondWithJSON(w, 401, err.Error())
		}
		cfg.DB.AddRefreshTokenToUser(user, refreshToken)

		userOut := UserOutLogin{
			Email:        user.Email,
			UserId:       user.Id,
			IsChirpyRed:  user.IsChirpyRed,
			Token:        token,
			RefreshToken: refreshToken,
		}
		w = respondWithJSON(w, 200, userOut)

	} else {
		w = respondWithJSON(w, 405, "Method not allowed")
	}

}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	// helperPrintHeadersreq)

	if req.Method == http.MethodPost {
		tokenString := req.Header.Get("Authorization")
		if tokenString == "" {
			w = respondWithError(w, 500, "No Authorization Token")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		newToken, valid, err := cfg.DB.RefreshTokenValid(tokenString)
		if err != nil {
			w = respondWithError(w, 401, "Not  Authorized")
			return
		}

		payload := map[string]interface{}{
			"token": newToken,
		}

		if valid {
			w = respondWithJSON(w, 200, payload)
		} else {
			w = respondWithError(w, 401, "Fuck you")

		}

	} else {
		w = respondWithJSON(w, 405, "Method not allowed")
	}

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		tokenString := req.Header.Get("Authorization")
		if tokenString == "" {
			w = respondWithError(w, 500, "No Authorization Token")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		_, valid, err := cfg.DB.RefreshTokenValid(tokenString)
		if err != nil {
			w = respondWithError(w, 401, "Not  Authorized")
			return
		}

		if valid {
			cfg.DB.RevokeToken(tokenString)
			w.WriteHeader(204)
			return
		}

	} else {
		w = respondWithJSON(w, 405, "Method not allowed")
	}

}

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, req *http.Request) {
	apiKeyString := req.Header.Get("Authorization")
	if apiKeyString == "" {
		w.WriteHeader(401)
		return
	}

	apiKeyString = strings.TrimPrefix(apiKeyString, "ApiKey ")
	fmt.Println(apiKeyString)
	key := string(os.Getenv("POLKA_KEY"))
	if apiKeyString != key {
		w.WriteHeader(401)
		return
	}

	if req.Method == http.MethodPost {
		type parameters struct {
			Event string `json:"event"`
			Data  struct {
				UserId int `json:"user_id"`
			} `json:"data"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			w = respondWithError(w, 500, "Something went wrong")
			return
		}

		if params.Event != "user.upgraded" {
			w.WriteHeader(204)
			return
		}

		_, err = cfg.DB.GetUserById(params.Data.UserId)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		cfg.DB.UpgradeUser(params.Data.UserId)
		w.WriteHeader(204)

	}
}

func helperPrintHeaders(req *http.Request) {
	// Print the method and endpoint
	fmt.Printf("===========")
	fmt.Printf("Method: %s\n", req.Method)
	fmt.Printf("Endpoint: %s\n", req.URL.Path)

	// Print all headers
	for header, values := range req.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", header, value)
		}
	}
	fmt.Printf("\n\n")
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
	serverMux.HandleFunc("/api/refresh", apiCfg.handlerRefresh)
	serverMux.HandleFunc("/api/revoke", apiCfg.handlerRevoke)
	serverMux.HandleFunc("/api/polka/webhooks", apiCfg.handlerWebhook)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
