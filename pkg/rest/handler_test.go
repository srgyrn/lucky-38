package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/srgyrn/lucky-38/pkg/creating"
	"github.com/srgyrn/lucky-38/pkg/drawing"
	"github.com/srgyrn/lucky-38/pkg/listing"
)

func Test_createDeck(t *testing.T) {
	deckID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")
	type requestParams struct {
		body  string
		query map[string]string
	}
	tests := []struct {
		name       string
		reqParams  requestParams
		service    creating.Service
		want       creating.Deck
		wantStatus int
	}{
		{
			name: "valid",
			reqParams: requestParams{
				body:  `{"shuffled": false, "partial": true}`,
				query: map[string]string{"cards": "AS,KD,AC,KH"},
			},
			service: &mockCreateService{
				out: creating.Deck{
					ID:        deckID,
					Shuffled:  false,
					Remaining: 4,
					Cards: []creating.Card{
						{Code: "AS", Value: "ACE", Suit: "SPADES"},
						{Code: "KD", Value: "KING", Suit: "DIAMONDS"},
						{Code: "AC", Value: "ACE", Suit: "CLUBS"},
						{Code: "KH", Value: "KING", Suit: "HEARTS"},
					},
				},
			},
			want: creating.Deck{
				ID:        deckID,
				Shuffled:  false,
				Remaining: 4,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles error from service",
			reqParams: requestParams{
				body: `{"shuffled": false, "partial": true}`,
			},
			service:    &mockCreateService{err: &creating.InvalidCardErr{Card: creating.Card{}}},
			want:       creating.Deck{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "handles error from DB",
			reqParams: requestParams{
				body: `{"shuffled": false, "partial": true}`,
			},
			service:    &mockCreateService{err: creating.ErrCreate},
			want:       creating.Deck{},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := httprouter.New()
			router.POST("/deck", createDeck(tt.service))

			queryStr := url.Values{}
			if tt.reqParams.query != nil {
				for key, val := range tt.reqParams.query {
					queryStr.Add(key, val)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/deck?"+queryStr.Encode(), bytes.NewBufferString(tt.reqParams.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tt.wantStatus != rr.Code {
				t.Errorf("createDeck() status %d, want %d", rr.Code, tt.wantStatus)
				t.Logf(rr.Body.String())
				return
			}

			var got creating.Deck
			json.Unmarshal(rr.Body.Bytes(), &got)

			if http.StatusOK == tt.wantStatus && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDeck() got %v, want %v", got, tt.want)
			}
		})
	}
}

type mockCreateService struct {
	out creating.Deck
	err error
}

func (ms *mockCreateService) CreateDeck(deck creating.Deck) (creating.Deck, error) {
	return ms.out, ms.err
}

type mockListService struct {
	out listing.Deck
	err error
}

func (mls *mockListService) List(ID string) (listing.Deck, error) {
	return mls.out, mls.err
}

func Test_getDeck(t *testing.T) {
	deckID, _ := uuid.Parse("a251071b-662f-44b6-ba11-e24863039c59")
	type args struct {
		service listing.Service
		deckID  string
	}
	tests := []struct {
		name       string
		args       args
		want       listing.Deck
		wantStatus int
	}{
		{
			name: "handles not found",
			args: args{
				service: &mockListService{err: listing.ErrNotFound},
				deckID:  "asdf",
			},
			want:       listing.Deck{},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "handles db error",
			args: args{
				service: &mockListService{err: errors.New("test error")},
				deckID:  "a251071b-662f-44b6-ba11-e24863039c59",
			},
			want:       listing.Deck{},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "valid",
			args: args{
				service: &mockListService{
					out: listing.Deck{
						ID:        deckID,
						Shuffled:  false,
						Remaining: 2,
						Cards: []listing.Card{
							{
								ID:    1,
								Code:  "AS",
								Value: "ACE",
								Suit:  "SPADES",
							},
							{
								ID:    2,
								Code:  "AS",
								Value: "2",
								Suit:  "SPADES",
							},
						},
					},
				},
				deckID: "a251071b-662f-44b6-ba11-e24863039c59",
			},
			want: listing.Deck{
				ID:        deckID,
				Shuffled:  false,
				Remaining: 2,
				Cards: []listing.Card{
					{
						Code:  "AS",
						Value: "ACE",
						Suit:  "SPADES",
					},
					{
						Code:  "AS",
						Value: "2",
						Suit:  "SPADES",
					},
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := httprouter.New()
			router.GET("/deck/:id", getDeck(tt.args.service))

			req := httptest.NewRequest(http.MethodGet, "/deck/"+tt.args.deckID, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tt.wantStatus != rr.Code {
				t.Errorf("getDeck() status code %d, want %d", rr.Code, tt.wantStatus)
				return
			}

			var got listing.Deck
			json.Unmarshal(rr.Body.Bytes(), &got)
			if http.StatusOK == tt.wantStatus && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDeck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_drawCards(t *testing.T) {
	type args struct {
		amount int
		s      drawing.Service
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantResponse []drawing.Card
	}{
		{
			name: "handles not found",
			args: args{
				amount: 2,
				s: &mockDrawingService{
					out: []drawing.Card{},
					err: drawing.ErrNotFound,
				},
			},
			wantStatus:   http.StatusNotFound,
			wantResponse: []drawing.Card{},
		},
		{
			name: "handles invalid card amount",
			args: args{
				amount: 10,
				s: &mockDrawingService{
					out: []drawing.Card{},
					err: drawing.ErrInsufficientRemainingCard,
				},
			},
			wantStatus:   http.StatusBadRequest,
			wantResponse: []drawing.Card{},
		},
		{
			name: "valid",
			args: args{
				amount: 2,
				s: &mockDrawingService{
					out: []drawing.Card{
						{
							ID:    1,
							Value: "ACE",
							Suit:  "SPADES",
							Code:  "AS",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantResponse: []drawing.Card{
				{
					Value: "ACE",
					Suit:  "SPADES",
					Code:  "AS",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := httprouter.New()
			router.PATCH("/deck/:id/draw/:amount", drawCards(tt.args.s))

			uri := fmt.Sprintf("/deck/test-test-test/draw/%d", tt.args.amount)
			req := httptest.NewRequest(http.MethodPatch, uri, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tt.wantStatus != rr.Code {
				t.Errorf("getDeck() status code %d, want %d", rr.Code, tt.wantStatus)
				return
			}

			var got []drawing.Card
			json.Unmarshal(rr.Body.Bytes(), &got)
			if http.StatusOK == tt.wantStatus && !reflect.DeepEqual(got, tt.wantResponse) {
				t.Errorf("getDeck() = %v, want %v", got, tt.wantResponse)
			}
		})
	}
}

type mockDrawingService struct {
	out []drawing.Card
	err error
}

func (ms *mockDrawingService) Draw(deckID string, n int) ([]drawing.Card, error) {
	return ms.out, ms.err
}
