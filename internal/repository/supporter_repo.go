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

type SupporterRepo struct {
	Client *firestore.Client
}

func NewSupporterRepo(client *firestore.Client) *SupporterRepo {
	return &SupporterRepo{Client: client}
}

// 🟢 GetSupporterAll
func (r *SupporterRepo) GetSupporterAll() ([]models.Supporter, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("supporters").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var supporters []models.Supporter
	for _, doc := range docs {
		var supporter models.Supporter
		doc.DataTo(&supporter)
		supporters = append(supporters, supporter)
	}
	return supporters, nil
}

// 🟢 GetSupporterByID
func (r *SupporterRepo) GetSupporterByID(supporterID string) (*models.Supporter, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("supporters").Doc(supporterID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var supporter models.Supporter
	doc.DataTo(&supporter)
	return &supporter, nil
}

// 🟢 GetSupporterByCaseID
func (r *SupporterRepo) GetSupporterByCaseID(caseID string) (*models.Supporter, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("supporters").Where("case_id", "==", caseID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var supporter models.Supporter
	doc[0].DataTo(&supporter)
	return &supporter, nil
}

// 🟢 CreateSupporter - auto-generate ID SP-00001
func (r *SupporterRepo) CreateSupporter(supporter *models.Supporter) error {
	ctx := context.Background()

	docs, err := r.Client.Collection("supporters").OrderBy("supporter_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "SP-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["supporter_id"].(string)
		numStr := strings.TrimPrefix(lastID, "SP-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("SP-%05d", n+1)
		}
	}

	supporter.SupporterID = nextID
	now := time.Now()
	supporter.CreatedAt = now
	supporter.UpdatedAt = now

	_, err = r.Client.Collection("supporters").Doc(supporter.SupporterID).Set(ctx, supporter)
	return err
}

// 🟢 UpdateSupporterByID
func (r *SupporterRepo) UpdateSupporterByID(supporterID string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("supporters").Doc(supporterID).Set(ctx, data, firestore.MergeAll)
	return err
}
