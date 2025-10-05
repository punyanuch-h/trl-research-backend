package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"cloud.google.com/go/firestore"
)

type ResearcherRepo struct {
	Client *firestore.Client
}

func NewResearcherRepo(client *firestore.Client) *ResearcherRepo {
	return &ResearcherRepo{Client: client}
}

// 🟢 GetResearcherAll - fetch all researchers
func (r *ResearcherRepo) GetResearcherAll() ([]models.ResearcherInfo, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("researchers").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var researchers []models.ResearcherInfo
	for _, doc := range docs {
		var researcher models.ResearcherInfo
		doc.DataTo(&researcher)
		researchers = append(researchers, researcher)
	}
	return researchers, nil
}

// 🟢 GetResearcherByID - fetch one researcher by ID
func (r *ResearcherRepo) GetResearcherByID(researcherID string) (*models.ResearcherInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("researchers").Doc(researcherID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var researcher models.ResearcherInfo
	doc.DataTo(&researcher)
	return &researcher, nil
}

// 🟢 CreateResearcher - auto-generate ResearcherID and create new record
func (r *ResearcherRepo) CreateResearcher(researcher *models.ResearcherInfo) error {
	ctx := context.Background()

	// find last ID to generate next
	docs, err := r.Client.Collection("researchers").OrderBy("researcher_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "RS-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["researcher_id"].(string)
		numStr := strings.TrimPrefix(lastID, "RS-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("RS-%05d", n+1)
		}
	}

	researcher.ResearcherID = nextID
	now := time.Now()
	researcher.CreatedAt = now
	researcher.UpdatedAt = now

	_, err = r.Client.Collection("researchers").Doc(researcher.ResearcherID).Set(ctx, researcher)
	return err
}

// 🟢 UpdateResearcherByID - partial update with UpdatedAt
func (r *ResearcherRepo) UpdateResearcherByID(researcherID string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("researchers").Doc(researcherID).Set(ctx, data, firestore.MergeAll)
	return err
}
