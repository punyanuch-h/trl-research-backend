package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"trl-research-backend/internal/models"
)

type CoordinatorRepo struct {
	Client *firestore.Client
}

func NewCoordinatorRepo(client *firestore.Client) *CoordinatorRepo {
	return &CoordinatorRepo{Client: client}
}

// 🟢 GetCoordinatorAll - fetch all coordinators
func (r *CoordinatorRepo) GetCoordinatorAll() ([]models.CoordinatorInfo, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("coordinators").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var coordinators []models.CoordinatorInfo
	for _, doc := range docs {
		var coordinator models.CoordinatorInfo
		doc.DataTo(&coordinator)
		coordinators = append(coordinators, coordinator)
	}
	return coordinators, nil
}

// 🟢 GetCoordinatorByEmail - use email as document ID
func (r *CoordinatorRepo) GetCoordinatorByEmail(email string) (*models.CoordinatorInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("coordinators").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}

	var coordinator models.CoordinatorInfo
	doc.DataTo(&coordinator)
	return &coordinator, nil
}

// 🟢 CreateCoordinator
func (r *CoordinatorRepo) CreateCoordinator(coordinator *models.CoordinatorInfo) error {
	ctx := context.Background()
	now := time.Now()
	coordinator.CreatedAt = now
	coordinator.UpdatedAt = now

	_, err := r.Client.Collection("coordinators").Doc(coordinator.CoordinatorEmail).Set(ctx, coordinator)
	return err
}

// 🟢 UpdateCoordinatorByEmail
func (r *CoordinatorRepo) UpdateCoordinatorByEmail(email string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("coordinators").Doc(email).Set(ctx, data, firestore.MergeAll)
	return err
}
