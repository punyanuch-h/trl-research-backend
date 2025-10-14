package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"trl-research-backend/internal/models"
)

type AssessmentTrlRepo struct {
	Client *firestore.Client
}

func NewAssessmentTrlRepo(client *firestore.Client) *AssessmentTrlRepo {
	return &AssessmentTrlRepo{Client: client}
}

// 🟢 GetAssessmentTrlAll
func (r *AssessmentTrlRepo) GetAssessmentTrlAll() ([]models.AssessmentTrl, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("assessment_trl").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var assessments []models.AssessmentTrl
	for _, doc := range docs {
		var a models.AssessmentTrl
		doc.DataTo(&a)
		assessments = append(assessments, a)
	}
	return assessments, nil
}

// 🟢 GetAssessmentTrlByID
func (r *AssessmentTrlRepo) GetAssessmentTrlByID(id string) (*models.AssessmentTrl, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("assessment_trl").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var a models.AssessmentTrl
	doc.DataTo(&a)
	return &a, nil
}

// 🟢 GetAssessmentTrlByCaseID
func (r *AssessmentTrlRepo) GetAssessmentTrlByCaseID(caseID string) (*models.AssessmentTrl, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("assessment_trl").Where("case_id", "==", caseID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var a models.AssessmentTrl
	doc[0].DataTo(&a)
	return &a, nil
}

// 🟢 CreateAssessmentTrl - auto generate ID AS-00001
func (r *AssessmentTrlRepo) CreateAssessmentTrl(a *models.AssessmentTrl) error {
	ctx := context.Background()

	docs, err := r.Client.Collection("assessment_trl").OrderBy("id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "AS-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["id"].(string)
		numStr := strings.TrimPrefix(lastID, "AS-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("AS-%05d", n+1)
		}
	}

	a.ID = nextID
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err = r.Client.Collection("assessment_trl").Doc(a.ID).Set(ctx, a)
	return err
}

// 🟢 UpdateAssessmentTrlByID
func (r *AssessmentTrlRepo) UpdateAssessmentTrlByID(id string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("assessment_trl").Doc(id).Set(ctx, data, firestore.MergeAll)
	return err
}
