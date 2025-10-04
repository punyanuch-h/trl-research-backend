package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type AdminRepo struct{}

func (r *AdminRepo) CreateOrUpdate(data *models.AdminInfo) error {
	_, err := database.FirestoreClient.Collection("admin_info").
		Doc(data.AdminID).
		Set(database.Ctx, data)
	return err
}

func (r *AdminRepo) GetByID(id string) (*models.AdminInfo, error) {
	doc, err := database.FirestoreClient.Collection("admin_info").
		Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.AdminInfo
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *AdminRepo) DeleteAdmin(id string) error {
	_, err := database.FirestoreClient.Collection("admin_info").
		Doc(id).Delete(database.Ctx)
	return err
}
