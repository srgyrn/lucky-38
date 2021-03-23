package storage

import (
	"testing"

	"github.com/google/uuid"
)

func (r *Repository) TestTeardown(t *testing.T) {
	if _, err := r.db.Exec("TRUNCATE cards"); err != nil {
		t.Fatalf("TRUNCATE cards err: %v", err)
	}

	if _, err := r.db.Exec("ALTER SEQUENCE cards_card_id_seq RESTART WITH 1"); err != nil {
		t.Fatalf("resetting card_id failed, err: %v", err)
	}

	if _, err := r.db.Exec("DELETE FROM decks"); err != nil {
		t.Fatalf("DELETE FROM decks err: %v", err)
	}
}

func (r *Repository) TestCountCards(t *testing.T) int {
	var c int
	err := r.db.QueryRow("SELECT COUNT(card_id) FROM cards").Scan(&c)
	if err != nil {
		t.Fatalf("Scan() err: %v", err)
	}

	return c
}

func (r *Repository) TestCountDrawnCards(t *testing.T, deckID uuid.UUID) int {
	var c int
	err := r.db.QueryRow("SELECT COUNT(card_id) FROM cards where deck = $1 AND drawn = true", deckID).Scan(&c)
	if err != nil {
		t.Fatalf("Scan() err: %v", err)
	}

	return c
}

func (r *Repository) TestCountDecks(t *testing.T) int {
	var c int
	err := r.db.QueryRow("SELECT COUNT(deck_id) FROM decks").Scan(&c)
	if err != nil {
		t.Fatalf("Scan() err: %v", err)
	}

	return c
}

func (r *Repository) TestDeckRemaining(t *testing.T, ID uuid.UUID) int {
	var c int
	err := r.db.QueryRow("SELECT remaining FROM decks WHERE deck_id =  $1", ID).Scan(&c)
	if err != nil {
		t.Fatalf("Scan() err: %v", err)
	}

	return c
}

func (r *Repository) TestInitData(t *testing.T, statement string)  {
	t.Helper()
	_, err := r.db.Exec(statement)
	if err != nil {
		t.Fatalf("TestInitData() err: %v\nstatement: %s", err, statement)
	}
}