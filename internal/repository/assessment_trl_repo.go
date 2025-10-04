package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type AssessmentTrlRepo struct{}

// CreateOrUpdate → ใช้ CaseID เป็น document key
func (r *AssessmentTrlRepo) CreateOrUpdate(data *models.AssessmentTRL) error {
	_, err := database.FirestoreClient.Collection("assessment_trl").
		Doc(data.CaseID). // ใช้ CaseID เป็น unique key
		Set(database.Ctx, data)
	return err
}

func (r *AssessmentTrlRepo) GetByCaseID(caseID string) (*models.AssessmentTRL, error) {
	doc, err := database.FirestoreClient.Collection("assessment_trl").
		Doc(caseID).Get(database.Ctx)
	if err != nil {
		return nil, err
	}

	var obj models.AssessmentTRL
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *AssessmentTrlRepo) Delete(caseID string) error {
	_, err := database.FirestoreClient.Collection("assessment_trl").
		Doc(caseID).Delete(database.Ctx)
	return err
}
