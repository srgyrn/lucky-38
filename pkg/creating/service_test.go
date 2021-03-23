package creating

import (
	"errors"
	"reflect"
	"testing"
)

func Test_service_CreateDeck(t *testing.T) {
	type fields struct {
		r Repository
	}

	tests := []struct {
		name        string
		fields      fields
		deck        Deck
		want        Deck
		wantErr     bool
		errWantType error
	}{
		{
			name: "wrong value",
			fields: fields{
				r: &mockDB{},
			},
			deck: Deck{
				Shuffled:  false,
				Remaining: 2,
				Cards:     []Card{{Code: "AH"}, {Code: "14H"}},
			},
			want:        Deck{},
			wantErr:     true,
			errWantType: ErrInvalidCard,
		},
		{
			name: "unknown suit",
			fields: fields{
				r: &mockDB{},
			},
			deck: Deck{
				Shuffled:  false,
				Remaining: 2,
				Cards:     []Card{{Code: "AH"}, {Code: "10K"}},
			},
			want:        Deck{},
			wantErr:     true,
			errWantType: ErrInvalidCard,
		},
		{
			name: "partial/cards missing",
			fields: fields{
				r: &mockDB{},
			},
			deck: Deck{
				Shuffled:  false,
				Remaining: 5,
				Cards:     []Card{},
			},
			want:        Deck{},
			wantErr:     true,
			errWantType: ErrInvalidDeck,
		},
		{
			name: "partial/not shuffled",
			fields: fields{
				r: &mockDB{},
			},
			deck: Deck{
				Shuffled:  false,
				Remaining: 4,
				Cards: []Card{
					{Code: "AS"},
					{Code: "KD"},
					{Code: "AC"},
					{Code: "KH"},
				},
			},
			want: Deck{
				Shuffled:  false,
				Remaining: 4,
				Cards: []Card{
					{Code: "AS", Value: "ACE", Suit: "SPADES"},
					{Code: "KD", Value: "KING", Suit: "DIAMONDS"},
					{Code: "AC", Value: "ACE", Suit: "CLUBS"},
					{Code: "KH", Value: "KING", Suit: "HEARTS"},
				},
			},
			wantErr: false,
		},
		{
			name: "handles db fail",
			fields: fields{
				r: &mockDB{
					err: errors.New("create error"),
				},
			},
			deck: Deck{
				Shuffled:  false,
				Remaining: 4,
				Cards: []Card{
					{Code: "AS"},
					{Code: "KD"},
					{Code: "AC"},
					{Code: "KH"},
				},
			},
			want:        Deck{},
			wantErr:     true,
			errWantType: ErrCreate,
		},
		{
			name: "full/not shuffled",
			fields: fields{
				r: &mockDB{},
			},
			deck: Deck{
				Shuffled:  false,
				Remaining: 52,
			},
			want: Deck{
				Shuffled:  false,
				Remaining: 52,
				Cards:     fullDeck,
			},
			wantErr:     false,
			errWantType: ErrCreate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				r: tt.fields.r,
			}
			got, err := s.CreateDeck(tt.deck)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !errors.As(err, &tt.errWantType) {
				t.Errorf("CreateDeck() error want %T, got %T", tt.errWantType, err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDeck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CreateDeck_shuffled(t *testing.T) {
	tests := []struct {
		name      string
		deck      Deck
		initCards []Card
	}{
		{
			name: "partial",
			deck: Deck{
				Shuffled:  true,
				Remaining: 4,
				Cards: []Card{
					{Code: "AS"},
					{Code: "KD"},
					{Code: "AC"},
					{Code: "KH"},
				},
			},
			initCards: []Card{
				{Code: "AS", Value: "ACE", Suit: "SPADES"},
				{Code: "KD", Value: "KINGS", Suit: "DIAMONDS"},
				{Code: "AC", Value: "ACE", Suit: "CLUBS"},
				{Code: "KH", Value: "KINGS", Suit: "HEARTS"},
			},
		},
		{
			name: "full",
			deck: Deck{
				Shuffled:  true,
				Remaining: 52,
			},
			initCards: fullDeck,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{r: &mockDB{}}
			got, err := s.CreateDeck(tt.deck)
			if err != nil {
				t.Errorf("CreateDeck() error = %v", err)
				return
			}

			if reflect.DeepEqual(got.Cards, tt.initCards) {
				t.Errorf("CreateDeck() got = %v, want cards shuffled", got)
			}
		})
	}
}

type mockDB struct {
	err error
}

func (mdb *mockDB) CreateDeck(deck *Deck) error {
	return mdb.err
}

var fullDeck []Card = []Card{
	{Code: "AS", Value: "A", Suit: "SPADES"},
	{Code: "2S", Value: "2", Suit: "SPADES"},
	{Code: "3S", Value: "3", Suit: "SPADES"},
	{Code: "4S", Value: "4", Suit: "SPADES"},
	{Code: "5S", Value: "5", Suit: "SPADES"},
	{Code: "6S", Value: "6", Suit: "SPADES"},
	{Code: "7S", Value: "7", Suit: "SPADES"},
	{Code: "8S", Value: "8", Suit: "SPADES"},
	{Code: "9S", Value: "9", Suit: "SPADES"},
	{Code: "10S", Value: "10", Suit: "SPADES"},
	{Code: "JS", Value: "JACK", Suit: "SPADES"},
	{Code: "QS", Value: "QUEEN", Suit: "SPADES"},
	{Code: "KS", Value: "KING", Suit: "SPADES"},
	{Code: "AD", Value: "A", Suit: "DIAMONDS"},
	{Code: "2D", Value: "2", Suit: "DIAMONDS"},
	{Code: "3D", Value: "3", Suit: "DIAMONDS"},
	{Code: "4D", Value: "4", Suit: "DIAMONDS"},
	{Code: "5D", Value: "5", Suit: "DIAMONDS"},
	{Code: "6D", Value: "6", Suit: "DIAMONDS"},
	{Code: "7D", Value: "7", Suit: "DIAMONDS"},
	{Code: "8D", Value: "8", Suit: "DIAMONDS"},
	{Code: "9D", Value: "9", Suit: "DIAMONDS"},
	{Code: "10D", Value: "10", Suit: "DIAMONDS"},
	{Code: "JD", Value: "JACK", Suit: "DIAMONDS"},
	{Code: "QD", Value: "QUEEN", Suit: "DIAMONDS"},
	{Code: "KD", Value: "KING", Suit: "DIAMONDS"},
	{Code: "AC", Value: "A", Suit: "CLUBS"},
	{Code: "2C", Value: "2", Suit: "CLUBS"},
	{Code: "3C", Value: "3", Suit: "CLUBS"},
	{Code: "4C", Value: "4", Suit: "CLUBS"},
	{Code: "5C", Value: "5", Suit: "CLUBS"},
	{Code: "6C", Value: "6", Suit: "CLUBS"},
	{Code: "7C", Value: "7", Suit: "CLUBS"},
	{Code: "8C", Value: "8", Suit: "CLUBS"},
	{Code: "9C", Value: "9", Suit: "CLUBS"},
	{Code: "10C", Value: "10", Suit: "CLUBS"},
	{Code: "JC", Value: "JACK", Suit: "CLUBS"},
	{Code: "QC", Value: "QUEEN", Suit: "CLUBS"},
	{Code: "KC", Value: "KING", Suit: "CLUBS"},
	{Code: "AH", Value: "A", Suit: "HEARTS"},
	{Code: "2H", Value: "2", Suit: "HEARTS"},
	{Code: "3H", Value: "3", Suit: "HEARTS"},
	{Code: "4H", Value: "4", Suit: "HEARTS"},
	{Code: "5H", Value: "5", Suit: "HEARTS"},
	{Code: "6H", Value: "6", Suit: "HEARTS"},
	{Code: "7H", Value: "7", Suit: "HEARTS"},
	{Code: "8H", Value: "8", Suit: "HEARTS"},
	{Code: "9H", Value: "9", Suit: "HEARTS"},
	{Code: "10H", Value: "10", Suit: "HEARTS"},
	{Code: "JH", Value: "JACK", Suit: "HEARTS"},
	{Code: "QH", Value: "QUEEN", Suit: "HEARTS"},
	{Code: "KH", Value: "KING", Suit: "HEARTS"},
}
