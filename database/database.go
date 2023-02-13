package database

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/kamaal111/forex-api/utils"
)

func CreateClient(ctx context.Context) (*firestore.Client, error) {
	projectID := utils.UnwrapEnvironment("GCP_PROJECT_ID")
	return firestore.NewClient(ctx, projectID)
}
