package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/kamaal111/forex-api/database"
	"github.com/kamaal111/forex-api/utils"
)

type FirestoreRatesRepository struct {
	client *firestore.Client
	ctx    context.Context
}

func NewFirestoreRatesRepository(ctx context.Context, client *firestore.Client) *FirestoreRatesRepository {
	return &FirestoreRatesRepository{client: client, ctx: ctx}
}

func (r *FirestoreRatesRepository) GetLatestRate(base string) (*ExchangeRateRecord, error) {
	documents := r.client.Collection("exchange_rates").
		OrderBy("date", firestore.Desc).
		Where("base", "==", base).
		Limit(1).
		Documents(r.ctx)

	var document *firestore.DocumentSnapshot
	var err error
	for document == nil {
		document, err = documents.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if document == nil {
		return nil, nil
	}

	var record ExchangeRateRecord
	err = document.DataTo(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

var ErrRatesNotFound = errors.New("rates not found")

func GetLatest(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	client, err := database.CreateClient(ctx)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	repo := NewFirestoreRatesRepository(ctx, client)
	service := NewRatesService(repo)

	base := request.URL.Query().Get("base")
	symbols := request.URL.Query().Get("symbols")

	record, err := service.GetLatestRate(base, symbols)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if record == nil {
		utils.ErrorHandler(writer, "Rates not found", http.StatusNotFound)
		return
	}

	output, err := json.Marshal(record)
	if err != nil {
		utils.ErrorHandler(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("content-type", "application/json")
	writer.Write(output)
}
