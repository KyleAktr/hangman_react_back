package main

import (
	"encoding/json"
	"net/http"
)

type GameState struct {
	Word        string   `json:"word"`        // Le mot à deviner
	Guesses     []string `json:"guesses"`     // Lettres devinées
	Attempts    int      `json:"attempts"`    // Essais restants
	IsCompleted bool     `json:"isCompleted"` // Jeu terminé ?
	Message     string   `json:"message"`     // Message à afficher
}

var game GameState

// Middleware pour activer CORS
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/start", startGame)
	mux.HandleFunc("/guess", guessLetter)

	http.ListenAndServe(":8080", enableCors(mux))
}

func startGame(w http.ResponseWriter, r *http.Request) {
	game = GameState{
		Word:        "example", // Ton mot ici
		Guesses:     []string{},
		Attempts:    6,
		IsCompleted: false,
		Message:     "Game started!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func guessLetter(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Letter string `json:"letter"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	// Exemple simple de logique
	if input.Letter != "" && len(input.Letter) == 1 {
		game.Guesses = append(game.Guesses, input.Letter)
		game.Attempts--
	}
	if game.Attempts == 0 {
		game.IsCompleted = true
		game.Message = "Game over!"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}
