package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type TrlQuestionRepo struct{}

// ใช้ composite key (assessment_id + question_code) เป็น DocID
func docID(assessmentID, questionCode string) string {
	return assessmentID + "_" + questionCode
}

func (r *TrlQuestionRepo) CreateOrUpdate(data *models.TrlQuestion) error {
	_, err := database.FirestoreClient.Collection("trl_questions").
		Doc(docID(data.AssessmentID, data.QuestionCode)).
		Set(database.Ctx, data)
	return err
}

func (r *TrlQuestionRepo) Get(assessmentID, questionCode string) (*models.TrlQuestion, error) {
	doc, err := database.FirestoreClient.Collection("trl_questions").
		Doc(docID(assessmentID, questionCode)).
		Get(database.Ctx)
	if err != nil {
		return nil, err
	}

	var obj models.TrlQuestion
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *TrlQuestionRepo) GetByAssessment(assessmentID string) ([]models.TrlQuestion, error) {
	docs, err := database.FirestoreClient.Collection("trl_questions").
		Where("assessment_id", "==", assessmentID).
		Documents(database.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var results []models.TrlQuestion
	for _, d := range docs {
		var q models.TrlQuestion
		_ = d.DataTo(&q)
		results = append(results, q)
	}
	return results, nil
}

func (r *TrlQuestionRepo) Delete(assessmentID, questionCode string) error {
	_, err := database.FirestoreClient.Collection("trl_questions").
		Doc(docID(assessmentID, questionCode)).
		Delete(database.Ctx)
	return err
}
