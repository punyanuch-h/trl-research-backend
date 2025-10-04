package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type ResearcherRepo struct{}

func (r *ResearcherRepo) CreateOrUpdate(data *models.ResearcherInfo) error {
	_, err := database.FirestoreClient.Collection("researcher_info").
		Doc(data.ResearcherID).
		Set(database.Ctx, data)
	return err
}

func (r *ResearcherRepo) GetByID(id string) (*models.ResearcherInfo, error) {
	doc, err := database.FirestoreClient.Collection("researcher_info").
		Doc(id).
		Get(database.Ctx)
	if err != nil {
		return nil, err
	}

	var researcher models.ResearcherInfo
	err = doc.DataTo(&researcher)
	return &researcher, err
}

func (r *ResearcherRepo) Delete(id string) error {
	_, err := database.FirestoreClient.Collection("researcher_info").
		Doc(id).
		Delete(database.Ctx)
	return err
}
