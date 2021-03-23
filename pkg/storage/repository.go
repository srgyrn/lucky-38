package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/srgyrn/lucky-38/pkg/creating"
	"github.com/srgyrn/lucky-38/pkg/drawing"
	"github.com/srgyrn/lucky-38/pkg/listing"
)

// Repository holds connection to db and implements creating.Repository
type Repository struct {
	ctx context.Context
	db  *sql.DB
}

func NewRepository(driver, source string) (*Repository, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	return &Repository{db: db}, nil
}

// DrawCards updates drawn status to true of n number of cards from deck with ID deckID
func (r *Repository) DrawCards(deckID uuid.UUID, cards ...drawing.Card) error {
	var whereIn []string
	for _, c := range cards {
		whereIn = append(whereIn, strconv.Itoa(c.ID))
	}

	r.ctx = context.Background()
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return fmt.Errorf("error at creating transaction: %v", err)
	}

	// update cards, set drawn = true
	statement := fmt.Sprintf("UPDATE cards SET drawn = true WHERE card_id IN (%s)", strings.Join(whereIn, ","))
	_, err = tx.ExecContext(r.ctx, statement)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update decks, set remaining = remaining - number_of_cards_drawn
	_, err = tx.ExecContext(r.ctx, fmt.Sprintf("UPDATE decks SET remaining = remaining - %d WHERE deck_id = $1", len(cards)), deckID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//FindAvailableCardByDeckID finds cards that are not drawn from the deck with given ID
func (r *Repository) FindAvailableCardByDeckID(deckID uuid.UUID) ([]drawing.Card, error) {
	query := `SELECT card_id, code, suit, value FROM cards WHERE deck = $1 AND drawn = $2 ORDER BY card_id DESC`
	rows, err := r.db.Query(query, deckID, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []drawing.Card{}, drawing.ErrNotFound
		}

		return []drawing.Card{}, err
	}
	defer rows.Close()

	var cards []drawing.Card
	var cardIDs []string
	for rows.Next() {
		card := drawing.Card{}
		err = rows.Scan(&card.ID, &card.Code, &card.Suit, &card.Value)
		if err != nil {
			return []drawing.Card{}, err
		}
		cards = append(cards, card)
		cardIDs = append(cardIDs, strconv.Itoa(card.ID))
	}

	if 0 == len(cards) {
		return cards, drawing.ErrNotFound
	}

	return cards, nil
}

// Find queries DB for the given deck ID and returns listing.Deck if found.
func (r *Repository) Find(ID uuid.UUID) (listing.Deck, error) {
	var deck listing.Deck
	err := r.db.QueryRow("SELECT deck_id, remaining, shuffled FROM decks WHERE deck_id = $1", ID).Scan(&deck.ID, &deck.Remaining, &deck.Shuffled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return listing.Deck{}, listing.ErrNotFound
		}

		return listing.Deck{}, err
	}

	query := `SELECT card_id, code, suit, value FROM cards WHERE deck = $1 AND drawn = $2`
	rows, err := r.db.Query(query, ID, false)
	if err != nil {
		return listing.Deck{}, err
	}
	defer rows.Close()

	for rows.Next() {
		card := listing.Card{}
		err = rows.Scan(&card.ID, &card.Code, &card.Suit, &card.Value)
		if err != nil {
			return listing.Deck{}, err
		}

		deck.Cards = append(deck.Cards, card)
	}

	return deck, nil
}

// CreateDeck inserts a new deck and cards to DB with given options
func (r *Repository) CreateDeck(deck *creating.Deck) error {
	deck.ID = uuid.New()

	r.ctx = context.Background()
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return fmt.Errorf("error at creating transaction: %v", err)
	}

	err = r.insertDeck(tx, deck)
	if err != nil {
		return err
	}

	deck.Cards, err = r.insertCard(tx, deck.ID, deck.Cards...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) insertDeck(tx *sql.Tx, deck *creating.Deck) error {
	statement := "INSERT INTO decks (deck_id, shuffled, remaining) VALUES ($1, $2, $3)"
	_, err := tx.ExecContext(r.ctx, statement, deck.ID, deck.Shuffled, deck.Remaining)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *Repository) insertCard(tx *sql.Tx, deckID uuid.UUID, cards ...creating.Card) ([]creating.Card, error) {
	var id int
	err := r.db.QueryRow("SELECT nextval(pg_get_serial_sequence('cards', 'card_id')) AS id").Scan(&id)
	if err != nil {
		tx.Rollback()
		return []creating.Card{}, errors.New("error retrieving next id")
	}

	var result []creating.Card
	statement := "INSERT INTO cards (code, value, suit, drawn, deck) VALUES ($1, $2, $3, $4, $5)"
	for _, c := range cards {
		if _, err := tx.ExecContext(r.ctx, statement, c.Code, c.Value, c.Suit, false, deckID); err != nil {
			tx.Rollback()
			return []creating.Card{}, fmt.Errorf("error at inserting card %v, err: %v", c, err)
		}

		c.ID = id
		result = append(result, c)
		id++
	}

	return result, nil
}
