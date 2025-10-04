package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type CaseRepo struct{}

func (r *CaseRepo) CreateOrUpdate(data *models.CaseInfo) error {
	_, err := database.FirestoreClient.Collection("case_info").
		Doc(data.CaseID).
		Set(database.Ctx, data)
	return err
}

func (r *CaseRepo) GetByID(id string) (*models.CaseInfo, error) {
	doc, err := database.FirestoreClient.Collection("case_info").
		Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.CaseInfo
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *CaseRepo) Delete(id string) error {
	_, err := database.FirestoreClient.Collection("case_info").
		Doc(id).Delete(database.Ctx)
	return err
}
