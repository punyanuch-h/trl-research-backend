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

type CaseRepo struct {
	Client *firestore.Client
}

func NewCaseRepo(client *firestore.Client) *CaseRepo {
	return &CaseRepo{Client: client}
}

// 🟢 GetCaseAll - fetch all cases
func (r *CaseRepo) GetCaseAll() ([]models.CaseInfo, error) {
	fmt.Println("GetCaseAll from repo")
	fmt.Println("r", r)
	ctx := context.Background()
	fmt.Println("ctx", ctx)
	docs, err := r.Client.Collection("cases").Documents(ctx).GetAll()
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	fmt.Println(docs)

	var cases []models.CaseInfo
	for _, doc := range docs {
		var cs models.CaseInfo
		doc.DataTo(&cs)
		cases = append(cases, cs)
	}
	fmt.Println(cases)
	return cases, nil
}

// 🟢 GetCaseByID
func (r *CaseRepo) GetCaseByID(caseID string) (*models.CaseInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("cases").Doc(caseID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var cs models.CaseInfo
	doc.DataTo(&cs)
	return &cs, nil
}

// 🟢 CreateCase - auto generate CaseID (CS-00001)
func (r *CaseRepo) CreateCase(cs *models.CaseInfo) error {
	ctx := context.Background()

	docs, err := r.Client.Collection("cases").OrderBy("case_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "CS-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["case_id"].(string)
		numStr := strings.TrimPrefix(lastID, "CS-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("CS-%05d", n+1)
		}
	}

	cs.CaseID = nextID
	now := time.Now()
	cs.CreatedAt = now
	cs.UpdatedAt = now

	_, err = r.Client.Collection("cases").Doc(cs.CaseID).Set(ctx, cs)
	return err
}

// 🟢 UpdateCaseByID
func (r *CaseRepo) UpdateCaseByID(caseID string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("cases").Doc(caseID).Set(ctx, data, firestore.MergeAll)
	return err
}
