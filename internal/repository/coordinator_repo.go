package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type CoordinatorRepo struct{}

func (r *CoordinatorRepo) CreateOrUpdate(data *models.CoordinatorInfo) error {
	_, err := database.FirestoreClient.Collection("coordinator_info").
		Doc(data.CoordinatorEmail).
		Set(database.Ctx, data)
	return err
}

func (r *CoordinatorRepo) GetByEmail(email string) (*models.CoordinatorInfo, error) {
	doc, err := database.FirestoreClient.Collection("coordinator_info").
		Doc(email).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.CoordinatorInfo
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *CoordinatorRepo) Delete(email string) error {
	_, err := database.FirestoreClient.Collection("coordinator_info").
		Doc(email).Delete(database.Ctx)
	return err
}
