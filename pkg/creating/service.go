package creating

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// FrenchDeckCardTotal holds the total number of cards that a French playing card deck has
const FrenchDeckCardTotal = 52

var suits = map[byte]string{
	'S': "SPADES",
	'D': "DIAMONDS",
	'H': "HEARTS",
	'C': "CLUBS",
}

var alphaValues = map[string]string{
	"A": "ACE",
	"K": "KING",
	"Q": "QUEEN",
	"J": "JACK",
}

type (
	Service interface {
		CreateDeck(Deck) (Deck, error)
	}
	Repository interface {
		CreateDeck(*Deck) error
	}

	service struct {
		r Repository
	}

	checkFn        func(Deck) error
	InvalidCardErr struct {
		Card Card
	}
)

var ErrInvalidDeck = errors.New("could not create deck")
var ErrCreate = errors.New("insert failed")
var ErrInvalidCard *InvalidCardErr

func (err *InvalidCardErr) Error() string {
	return fmt.Sprintf("invalid card: %v", err.Card)
}

func NewService(r Repository) Service {
	return &service{r: r}
}

// CreateDeck prepares cards in a deck, then communicates with repository to insert the deck and cards to DB.
//
// If any of the Deck.Cards have invalid value and/or suit (i.e. 50K or 10T), ErrInvalidCard is returned.
// If Deck.Cards length and Deck.Remaining are not equal, ErrInvalidDeck is returned.
// In case Repository returns an error, ErrCreate is returned.
func (s *service) CreateDeck(d Deck) (Deck, error) {
	checkCardSuit := func(deck Deck) error {
		for _, c := range deck.Cards {
			code := strings.TrimSpace(c.Code)
			n := len(code)
			if n < 2 || n > 3 {
				return &InvalidCardErr{c}
			}

			if _, ok := suits[code[n-1]]; !ok {
				return &InvalidCardErr{c}
			}
		}
		return nil
	}
	checkCardVal := func(deck Deck) error {
		for _, c := range deck.Cards {
			code := strings.TrimSpace(c.Code)
			n := len(code)
			if n < 2 || n > 3 {
				return &InvalidCardErr{c}
			}

			val := code[:n-1]
			if _, ok := alphaValues[val]; ok {
				continue
			}

			if valInt, _ := strconv.Atoi(val); 2 > valInt || 11 < valInt {
				return &InvalidCardErr{c}
			}
		}

		return nil
	}
	checkCardAmount := func(deck Deck) error {
		if deck.Remaining != FrenchDeckCardTotal && deck.Remaining != len(deck.Cards) {
			return ErrInvalidDeck
		}

		return nil
	}

	// Fill missing values of Card from code if deck is partial
	for i := 0; i < len(d.Cards); i++ {
		card := &d.Cards[i]
		code := strings.TrimSpace(card.Code)
		val := code[:len(code)-1]
		suit := code[len(code)-1]

		card.Value = val
		if v, ok := alphaValues[val]; ok {
			card.Value = v
		}

		card.Suit = suits[suit]
	}

	// validate cards
	for _, fn := range []checkFn{checkCardAmount, checkCardVal, checkCardSuit} {
		if err := fn(d); err != nil {
			return Deck{}, err
		}
	}

	// Generate cards in order if not partial
	if FrenchDeckCardTotal == d.Remaining {
		d.Cards = fullDeckGen(FrenchDeckCardTotal)
	}

	if d.Shuffled {
		d.shuffleCards()
	}

	err := s.r.CreateDeck(&d)
	if err != nil {
		return Deck{}, ErrCreate
	}

	return d, nil
}

// fullDeckGen generates all cards of a deck in order both value and suit
func fullDeckGen(cardAmount int) []Card {
	cards := make([]Card, cardAmount, cardAmount)
	var index int

	for _, s := range []byte{'S', 'D', 'C', 'H'} {
		suitCode := string(s)
		suit := suits[s]

		for j := 1; j <= 10; j++ {
			val := strconv.Itoa(j)
			if 1 == j {
				val = "A"
			}
			cards[index] = Card{
				Code:  val + suitCode,
				Value: val,
				Suit:  suit,
			}
			index++
		}

		for _, v := range []string{"JACK", "QUEEN", "KING"} {
			cards[index] = Card{
				Code:  string(v[0]) + suitCode,
				Value: v,
				Suit:  suit,
			}
			index++
		}
	}

	return cards
}

func (d *Deck) shuffleCards() {
	for i := 0; i < len(d.Cards); i++ {
		j := i + (rand.Intn(len(d.Cards)) % (len(d.Cards) - i))
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}
