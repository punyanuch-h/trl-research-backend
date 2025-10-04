package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type SupporterRepo struct{}

func (r *SupporterRepo) CreateOrUpdate(data *models.Supporter) error {
	_, err := database.FirestoreClient.Collection("supporter").
		Doc(data.CaseID).
		Set(database.Ctx, data)
	return err
}

func (r *SupporterRepo) GetByCaseID(caseID string) (*models.Supporter, error) {
	doc, err := database.FirestoreClient.Collection("supporter").
		Doc(caseID).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.Supporter
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *SupporterRepo) Delete(caseID string) error {
	_, err := database.FirestoreClient.Collection("supporter").
		Doc(caseID).Delete(database.Ctx)
	return err
}
