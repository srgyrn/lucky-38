package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/srgyrn/lucky-38/pkg/creating"
	"github.com/srgyrn/lucky-38/pkg/drawing"
	"github.com/srgyrn/lucky-38/pkg/listing"
)

// Handler creates a new router, registers routes and returns the created router.
func Handler(cs creating.Service, ls listing.Service, ds drawing.Service) http.Handler {
	router := httprouter.New()

	router.GET("/health", health())
	router.POST("/decks", createDeck(cs))
	router.GET("/decks/:id", getDeck(ls))
	router.PATCH("/decks/:id/draw/:amount", drawCards(ds))
	return router
}

// health checks if api is responsive
func health() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Why, hello!")
	}
}

// createDeck returns a handler for POST /deck requests
func createDeck(s creating.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var newDeck creating.Deck
		err := decoder.Decode(&newDeck)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newDeck.Remaining = creating.FrenchDeckCardTotal

		if "" != r.URL.Query().Get("cards") {
			cardCodes := strings.Split(r.URL.Query().Get("cards"), ",")
			newDeck.Remaining = len(cardCodes)
			for _, cc := range cardCodes {
				card := creating.Card{
					Code: strings.TrimSpace(cc),
				}

				newDeck.Cards = append(newDeck.Cards, card)
			}
		}

		newDeck, err = s.CreateDeck(newDeck)
		if errors.As(err, &creating.ErrInvalidCard) || errors.Is(err, creating.ErrInvalidDeck) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, creating.ErrCreate) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newDeck)
	}
}

// getDeck returns a handler for GET /deck/<deck_id> requests
func getDeck(s listing.Service) func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		deck, err := s.List(params.ByName("id"))
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, listing.ErrNotFound) {
				status = http.StatusNotFound
			}

			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deck)
	}
}

// drawCards returns a handler for PUT /deck/<deck_id> requests
func drawCards(s drawing.Service) func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		deckID := params.ByName("id")
		amount := params.ByName("amount")

		if "" == amount {
			http.Error(w, "amount cannot be empty", http.StatusBadRequest)
			return
		}

		n, err := strconv.Atoi(amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cards, err := s.Draw(deckID, n)
		if err != nil {
			if errors.Is(err, drawing.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			if errors.Is(err, drawing.ErrInsufficientRemainingCard) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}
}
