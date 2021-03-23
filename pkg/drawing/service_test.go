package drawing

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_service_Draw(t *testing.T) {
	type fields struct {
		r Repository
	}
	type args struct {
		deckID string
		n      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Card
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				r: &mockRepository{
					cards: []Card{
						{
							ID:    1,
							Value: "ACE",
							Suit:  "SPADES",
							Code:  "AS",
						},
						{
							ID:    2,
							Value: "2",
							Suit:  "SPADES",
							Code:  "2S",
						},
						{
							ID:    3,
							Value: "3",
							Suit:  "SPADES",
							Code:  "3S",
						},
					},
				},
			},
			args:    args{
				deckID: "a251071b-662f-44b6-ba11-e24863039c59",
				n:      2,
			},
			want:    []Card{
				{
					ID:    1,
					Value: "ACE",
					Suit:  "SPADES",
					Code:  "AS",
				},
				{
					ID:    2,
					Value: "2",
					Suit:  "SPADES",
					Code:  "2S",
				},
			},
			wantErr: false,
		},
		{
			name: "requested amount exceeds remaining",
			fields: fields{
				r: &mockRepository{
					cards: []Card{
						{
							ID:    1,
							Value: "ACE",
							Suit:  "SPADES",
							Code:  "AS",
						},
					},
				},
			},
			args: args{
				deckID: "a251071b-662f-44b6-ba11-e24863039c59",
				n:      5,
			},
			want:    []Card{},
			wantErr: true,
		},
		{
			name: "deck not found",
			fields: fields{
				r: &mockRepository{
					err:   ErrNotFound,
					cards: []Card{},
				},
			},
			args: args{
				deckID: "a251071b-662f-44b6-ba11-e24863039c59",
				n:      2,
			},
			want:    []Card{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				r: tt.fields.r,
			}
			got, err := s.Draw(tt.args.deckID, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Draw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Draw() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockRepository struct {
	err   error
	cards []Card
}

func (r *mockRepository) DrawCards(uuid.UUID, ...Card) error {
	return r.err
}

func (r *mockRepository) FindAvailableCardByDeckID(uuid.UUID) ([]Card, error) {
	return r.cards, r.err
}
