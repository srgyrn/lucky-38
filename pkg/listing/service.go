package listing

import (
	"errors"

	"github.com/google/uuid"
)

type (
	Deck struct {
		ID        uuid.UUID `json:"deck_id"`
		Shuffled  bool      `json:"shuffled"`
		Remaining int       `json:"remaining"`
		Cards     []Card    `json:"cards"`
	}

	Card struct {
		ID    int    `json:"-"`
		Code  string `json:"code"`
		Value string `json:"value"`
		Suit  string `json:"suit"`
	}

	Service interface {
		List(ID string) (Deck, error)
	}

	Repository interface {
		Find(ID uuid.UUID) (Deck, error)
	}

	service struct {
		r Repository
	}
)

var ErrNotFound = errors.New("deck not found")

func NewService(r Repository) Service {
	return &service{r: r}
}

// List uses Repository to retrieve Deck by given deck id from DB.
// If deck is not found, ErrNotFound is returned.
func (s *service) List(ID string) (Deck, error) {
	deckID, err := uuid.Parse(ID)
	if err != nil {
		return Deck{}, err
	}
	deck, err := s.r.Find(deckID)
	if err != nil {
		return Deck{}, ErrNotFound
	}
	return deck, nil
}
