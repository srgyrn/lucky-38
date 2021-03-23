package creating

import "github.com/google/uuid"

type Deck struct {
	ID        uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"-"`
}
