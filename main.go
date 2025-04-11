package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sabrek15/chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db 	*database.Queries
	platform string
}


type errorResponse struct {
	Error string `json:"error"`
}


type createUser struct {
	Email	string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

const metric string = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf(metric, cfg.fileserverHits.Load())
	w.Write([]byte(msg))

}

func (cfg *apiConfig) userResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err := cfg.db.DeteleUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete users")
		return 
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func cleanChirp(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")

	for i, word := range words {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	responseBody, _ := json.Marshal(errorResponse{Error: msg})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(responseBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	responseBody, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(responseBody)
}


func (cfg *apiConfig) PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	defer r.Body.Close()
	var req createUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went Wrong")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) postChirpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	defer r.Body.Close()
	type parameters struct {
		Body	string `json:"body"`
		UserID 	uuid.UUID `json:"user_id"`
	}
	var req parameters
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went Wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanChirp(req.Body)

	chirp, err := cfg.db.CreateChrips(r.Context(), database.CreateChripsParams{Body: cleanedBody, UserID: req.UserID})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't create chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	defer r.Body.Close()
	
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Fatal(err)
	}
	respondWithJSON(w, http.StatusCreated, chirps)
}

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chirpIDstr := r.PathValue("chirpid")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't parse chirp id")
		return
	}
	
	defer r.Body.Close()
	
	chirp, err := cfg.db.GetChirpsByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if dbURL == "" {
		log.Fatal("DB_URL not found in env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open the db: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("couldn't ping the db: %s", err)
	}

	dbQueries := database.New(db)


	cfg := apiConfig{db: dbQueries, platform: platform}

	serverHandler := http.NewServeMux()

	serverHandler.Handle("/app/", http.StripPrefix("/app/", middlewareLog(cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))))
	// serverHandler.Handle("/assets", http.FileServer(http.Dir(".")))

	serverHandler.HandleFunc("GET /api/healthz", readinessHandler)
	serverHandler.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serverHandler.HandleFunc("POST /admin/reset", cfg.userResetHandler)
	serverHandler.HandleFunc("POST /api/users", cfg.PostUsersHandler)
	serverHandler.HandleFunc("POST /api/chirps", cfg.postChirpsHandler)
	serverHandler.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	serverHandler.HandleFunc("GET /api/chirps/{chirpid}", cfg.getChirpByID)

	server := &http.Server{
		Addr:    ":8080",
		Handler: serverHandler,
	}

	server.ListenAndServe()
}
