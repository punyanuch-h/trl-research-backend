package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type AssessmentTrlPart2Repo struct{}

func (r *AssessmentTrlPart2Repo) CreateOrUpdate(data *models.AssessmentTRLPart2) error {
	_, err := database.FirestoreClient.Collection("assessment_trl_part2").
		Doc(data.ID).
		Set(database.Ctx, data)
	return err
}

func (r *AssessmentTrlPart2Repo) GetByID(id string) (*models.AssessmentTRLPart2, error) {
	doc, err := database.FirestoreClient.Collection("assessment_trl_part2").
		Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.AssessmentTRLPart2
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *AssessmentTrlPart2Repo) Delete(id string) error {
	_, err := database.FirestoreClient.Collection("assessment_trl_part2").
		Doc(id).Delete(database.Ctx)
	return err
}
