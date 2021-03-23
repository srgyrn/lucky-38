package storage_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/srgyrn/lucky-38/pkg/config"
	"github.com/srgyrn/lucky-38/pkg/creating"
	"github.com/srgyrn/lucky-38/pkg/drawing"
	"github.com/srgyrn/lucky-38/pkg/listing"
	"github.com/srgyrn/lucky-38/pkg/storage"
)

func TestRepository_CreateDeck(t *testing.T) {
	r := getRepository(t)
	defer r.TestTeardown(t)

	deck := creating.Deck{
		Shuffled:  false,
		Remaining: 3,
		Cards: []creating.Card{
			{
				Code:  "AS",
				Value: "ACE",
				Suit:  "SPADES",
			},
			{
				Code:  "2S",
				Value: "2",
				Suit:  "SPADES",
			},
			{
				Code:  "3S",
				Value: "3",
				Suit:  "SPADES",
			},
		},
	}
	want := creating.Deck{
		Shuffled:  false,
		Remaining: 3,
		Cards: []creating.Card{
			{
				ID:    1,
				Code:  "AS",
				Value: "ACE",
				Suit:  "SPADES",
			},
			{
				ID:    2,
				Code:  "2S",
				Value: "2",
				Suit:  "SPADES",
			},
			{
				ID:    3,
				Code:  "3S",
				Value: "3",
				Suit:  "SPADES",
			},
		},
	}

	err := r.CreateDeck(&deck)
	if err != nil {
		t.Errorf("CreateDeck() error = %v", err)
		return
	}
	want.ID = deck.ID
	if !reflect.DeepEqual(deck, want) {
		t.Errorf("CreateDeck() got = %v, want %v", deck, want)
	}
}

func TestRepository_Find(t *testing.T) {
	r := getRepository(t)

	t.Run("missing ID", func(t *testing.T) {
		defer r.TestTeardown(t)

		deckID, _ := uuid.Parse("69077400-88cd-11eb-8dcd-0242ac130003")
		_, err := r.Find(deckID)
		if !errors.Is(err, listing.ErrNotFound) {
			t.Errorf("Find() want %T, got = %v", listing.ErrNotFound, err)
		}
	})

	t.Run("valid find", func(t *testing.T) {
		defer r.TestTeardown(t)

		migration := getMigrationSQL(t, filepath.Join("testdata", "migrations", "insert_deck.sql"))
		r.TestInitData(t, migration)

		deckID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")
		got, err := r.Find(deckID)
		if err != nil {
			t.Errorf("Find() error = %v", err)
			return
		}

		want := listing.Deck{
			ID:        deckID,
			Shuffled:  false,
			Remaining: 4,
			Cards: []listing.Card{
				{
					ID:    1,
					Code:  "AS",
					Value: "ACE",
					Suit:  "SPADES",
				},
				{
					ID:    2,
					Code:  "2S",
					Value: "2",
					Suit:  "SPADES",
				},
				{
					ID:    3,
					Code:  "3S",
					Value: "3",
					Suit:  "SPADES",
				},
				{
					ID:    4,
					Code:  "4S",
					Value: "4",
					Suit:  "SPADES",
				},
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Find() got = %v, want %v", got, want)
		}
	})
}

func getRepository(t *testing.T) *storage.Repository {
	t.Helper()
	conf, err := config.Load("../../")
	if err != nil {
		t.Fatalf("config.Load() err: %v", err)
	}

	if conf.Driver == "" || conf.Source == "" {
		t.Fatalf("unexpected value for driver(%s) or source(%s)", conf.Driver, conf.Source)
	}

	r, err := storage.NewRepository(conf.Driver, conf.Source)
	if err != nil {
		t.Fatalf("storage.NewRepository() error = %v", err)
	}

	return r
}

func TestRepository_DrawCards(t *testing.T) {
	r := getRepository(t)
	defer r.TestTeardown(t)

	migration := getMigrationSQL(t, filepath.Join("testdata", "migrations", "draw_card_insert.sql"))
	r.TestInitData(t, migration)

	deckID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")
	drawnCards := []drawing.Card{{ID: 4}, {ID: 5}}
	n := 2
	wantRemaining := r.TestDeckRemaining(t, deckID) - n

	err := r.DrawCards(deckID, drawnCards...)
	if err != nil {
		t.Errorf("DrawCards() error = %v", err)
		return
	}

	gotRemaining := r.TestDeckRemaining(t, deckID)
	if gotRemaining != wantRemaining {
		t.Errorf("deck remaining: %d, want: %d", gotRemaining, wantRemaining)
		return
	}

	drawnCardCount := r.TestCountDrawnCards(t, deckID)
	if drawnCardCount != n {
		t.Errorf("drawn card count %d, want %d", drawnCardCount, n)
	}
}

func TestRepository_FindAvailableCardByDeckID(t *testing.T) {
	r := getRepository(t)
	defer r.TestTeardown(t)

	migration := getMigrationSQL(t, filepath.Join("testdata", "migrations", "insert_deck.sql"))
	r.TestInitData(t, migration)

	deckID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")
	want := []drawing.Card{
		{
			ID:    4,
			Value: "4",
			Suit:  "SPADES",
			Code:  "4S",
		},
		{
			ID:    3,
			Value: "3",
			Suit:  "SPADES",
			Code:  "3S",
		},
		{
			ID:    2,
			Value: "2",
			Suit:  "SPADES",
			Code:  "2S",
		},
		{
			ID:    1,
			Value: "ACE",
			Suit:  "SPADES",
			Code:  "AS",
		},
	}

	missingDeckID, _ := uuid.Parse("69077400-88cd-11eb-8dcd-0242ac130003")
	_, err := r.FindAvailableCardByDeckID(missingDeckID)
	if err == nil {
		t.Errorf("DrawCards() want error = %v got %v", drawing.ErrNotFound, err)
		return
	}

	got, err := r.FindAvailableCardByDeckID(deckID)
	if err != nil {
		t.Errorf("DrawCards() error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAvailableCardByDeckID() = %v, want %v", got, want)
	}
}

func getMigrationSQL(t *testing.T, filepath string) string {
	t.Helper()
	f, err := os.Open(filepath)
	if err != nil {
		t.Fatalf("could get migration file: %v", err)
	}
	defer f.Close()

	query, _ := ioutil.ReadAll(f)
	return string(query)
}