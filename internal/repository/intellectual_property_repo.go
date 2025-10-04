package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type IntellectualPropertyRepo struct{}

func (r *IntellectualPropertyRepo) CreateOrUpdate(data *models.IntellectualProperty) error {
	_, err := database.FirestoreClient.Collection("intellectual_property").
		Doc(data.IPRequestNumber).
		Set(database.Ctx, data)
	return err
}

func (r *IntellectualPropertyRepo) GetByID(id string) (*models.IntellectualProperty, error) {
	doc, err := database.FirestoreClient.Collection("intellectual_property").
		Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.IntellectualProperty
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *IntellectualPropertyRepo) Delete(id string) error {
	_, err := database.FirestoreClient.Collection("intellectual_property").
		Doc(id).Delete(database.Ctx)
	return err
}
