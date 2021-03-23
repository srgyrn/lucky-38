package drawing

import (
	"errors"

	"github.com/google/uuid"
)

type (
	Card struct {
		ID    int    `json:"-"`
		Value string `json:"value"`
		Suit  string `json:"suit"`
		Code  string `json:"code"`
	}

	Repository interface {
		FindAvailableCardByDeckID(uuid.UUID) ([]Card, error)
		DrawCards(uuid.UUID, ...Card) error
	}

	Service interface {
		Draw(deckID string, n int) ([]Card, error)
	}

	service struct {
		r Repository
	}
)

var ErrNotFound = errors.New("deck or remaining cards not found")
var ErrInsufficientRemainingCard = errors.New("remaining cards are less than the requested amount to draw")

func NewService(r Repository) Service {
	return &service{r: r}
}

// Draw marks n amount of cards as "drawn" from the deck with given deckID and returns them.
// If n is less than the number of available cards, ErrInsufficientRemainingCard is returned.
func (s *service) Draw(deckID string, n int) ([]Card, error) {
	deckUUID, err := uuid.Parse(deckID)
	if err != nil {
		return []Card{}, err
	}
	cards, err := s.r.FindAvailableCardByDeckID(deckUUID)
	if err != nil {
		return []Card{}, err
	}

	if len(cards) < n {
		return []Card{}, ErrInsufficientRemainingCard
	}

	result := cards[:n]
	err = s.r.DrawCards(deckUUID, result...)
	if err != nil {
		return []Card{}, err
	}

	return result, nil
}
