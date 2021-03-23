package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/srgyrn/lucky-38/pkg/config"
	"github.com/srgyrn/lucky-38/pkg/creating"
)
func TestRepository_insertDeck(t *testing.T) {
	r := getRepository(t)
	defer r.TestTeardown(t)
	tx, _ := r.db.BeginTx(r.ctx, nil)
	deckID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")
	deck := creating.Deck{ID: deckID, Shuffled: false, Remaining: 3}
	err := r.insertDeck(tx, &deck)
	if err != nil {
		t.Errorf("insertDeck() err: %v", err)
		return
	}

	tx.Commit()

	if 0 == r.TestCountDecks(t) {
		t.Errorf("deck count does not match")
	}
}


func TestRepository_insertCard(t *testing.T) {
	r := getRepository(t)

	deckInsertID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")

	type tests struct {
		name    string
		deckID  string
		cards   []creating.Card
		want    []creating.Card
		wantErr bool
	}

	for _, tt := range []tests{
		{
			name:   "valid insert",
			deckID: "a251071b-662f-44b6-ba11-e24863039c59",
			cards: []creating.Card{
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
			want: []creating.Card{
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
			wantErr: false,
		},
		{
			name: "missing deck",
			deckID: "69077400-88cd-11eb-8dcd-0242ac130003",
			cards: []creating.Card{
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
			},
			want:    []creating.Card{
				{
					ID: 1,
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
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			defer r.TestTeardown(t)

			deckID, _ := uuid.Parse(tt.deckID)
			tx, _ := r.db.BeginTx(r.ctx, nil)
			err := r.insertDeck(tx, &creating.Deck{ID: deckInsertID, Shuffled: false, Remaining: 3})
			if err != nil {
				t.Fatalf("insertDeck() err: %v", err)
			}

			got, err := r.insertCard(tx, deckID, tt.cards...)
			if err != nil {
				t.Errorf("insertCard() err: %v", err)
				return
			}

			if tt.wantErr {
				t.Logf("successfully received error: %v", err)
			}

			if err = tx.Commit(); (err != nil) != tt.wantErr {
				t.Fatalf("tx.Commit() err %v", err)
			}

			var wantCount int
			if !tt.wantErr {
				wantCount = len(tt.want)
			}

			if count := r.TestCountCards(t); !tt.wantErr && count != wantCount {
				t.Errorf("Card count %d, want %d", count, wantCount)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insertCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getRepository(t *testing.T) *Repository {
	t.Helper()
	conf, err := config.Load("../../")
	if err != nil {
		t.Fatalf("config.Load() err: %v", err)
	}

	if conf.Driver == "" || conf.Source == "" {
		t.Fatalf("unexpected value for driver(%s) or source(%s)", conf.Driver, conf.Source)
	}

	r, err := NewRepository(conf.Driver, conf.Source)
	if err != nil {
		t.Fatalf("storage.NewRepository() error = %v", err)
	}

	r.ctx = context.TODO()

	return r
}
