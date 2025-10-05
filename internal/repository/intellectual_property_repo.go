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

type IntellectualPropertyRepo struct {
	Client *firestore.Client
}

func NewIntellectualPropertyRepo(client *firestore.Client) *IntellectualPropertyRepo {
	return &IntellectualPropertyRepo{Client: client}
}

// 🟢 GetIPAll - fetch all intellectual property records
func (r *IntellectualPropertyRepo) GetIPAll() ([]models.IntellectualProperty, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("intellectual_properties").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var ips []models.IntellectualProperty
	for _, doc := range docs {
		var ip models.IntellectualProperty
		doc.DataTo(&ip)
		ips = append(ips, ip)
	}
	return ips, nil
}

// 🟢 GetIPByID
func (r *IntellectualPropertyRepo) GetIPByID(ipID string) (*models.IntellectualProperty, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("intellectual_properties").Doc(ipID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var ip models.IntellectualProperty
	doc.DataTo(&ip)
	return &ip, nil
}

// 🟢 CreateIP - auto generate ID IP-00001
func (r *IntellectualPropertyRepo) CreateIP(ip *models.IntellectualProperty) error {
	ctx := context.Background()

	docs, err := r.Client.Collection("intellectual_properties").OrderBy("id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "IP-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["id"].(string)
		numStr := strings.TrimPrefix(lastID, "IP-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("IP-%05d", n+1)
		}
	}

	ip.ID = nextID
	now := time.Now()
	ip.CreatedAt = now
	ip.UpdatedAt = now

	_, err = r.Client.Collection("intellectual_properties").Doc(ip.ID).Set(ctx, ip)
	return err
}

// 🟢 UpdateIPByID
func (r *IntellectualPropertyRepo) UpdateIPByID(ipID string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("intellectual_properties").Doc(ipID).Set(ctx, data, firestore.MergeAll)
	return err
}
