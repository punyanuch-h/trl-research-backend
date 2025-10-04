package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type AssessmentTrlPart1Repo struct{}

func (r *AssessmentTrlPart1Repo) CreateOrUpdate(data *models.AssessmentTRLPart1) error {
	_, err := database.FirestoreClient.Collection("assessment_trl_part1").
		Doc(data.ID).
		Set(database.Ctx, data)
	return err
}

func (r *AssessmentTrlPart1Repo) GetByID(id string) (*models.AssessmentTRLPart1, error) {
	doc, err := database.FirestoreClient.Collection("assessment_trl_part1").
		Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.AssessmentTRLPart1
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *AssessmentTrlPart1Repo) Delete(id string) error {
	_, err := database.FirestoreClient.Collection("assessment_trl_part1").
		Doc(id).Delete(database.Ctx)
	return err
}
