package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

type AdminRepo struct {
	Client *firestore.Client
}

func NewAdminRepo(client *firestore.Client) *AdminRepo {
	return &AdminRepo{Client: client}
}

// 🟢 Get all admins
func (r *AdminRepo) GetAdminAll() ([]models.AdminInfo, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("admins").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var admins []models.AdminInfo
	for _, doc := range docs {
		var admin models.AdminInfo
		doc.DataTo(&admin)
		admins = append(admins, admin)
	}
	return admins, nil
}

// 🟢 Get admin by ID
func (r *AdminRepo) GetAdminByID(adminID string) (*models.AdminInfo, error) {
	ctx := context.Background()

	// Query by admin_id field instead of document ID
	docs, err := r.Client.Collection("admins").Where("admin_id", "==", adminID).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("admin not found")
	}

	var admin models.AdminInfo
	docs[0].DataTo(&admin)
	return &admin, nil
}

// 🟢 Get admin by email
func (r *AdminRepo) GetAdminByEmail(email string) (*models.AdminInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("admins").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}

	var admin models.AdminInfo
	doc.DataTo(&admin)
	return &admin, nil
}

// 🟢 Create admin (auto-generate AdminID)
func (r *AdminRepo) CreateAdmin(admin *models.AdminInfo) error {
	ctx := context.Background()

	// find last ID to generate next
	docs, err := r.Client.Collection("admins").OrderBy("admin_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "SI-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["admin_id"].(string)
		numStr := strings.TrimPrefix(lastID, "SI-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("SI-%05d", n+1)
		}
	}

	// assign values
	admin.AdminID = nextID
	now := time.Now()
	admin.CreatedAt = now

	// save to Firestore using email as document ID
	_, err = r.Client.Collection("admins").Doc(admin.AdminEmail).Set(ctx, admin)
	return err
}

// 🟢 Login with password verification
func (r *AdminRepo) Login(email string, password string) (*models.AdminInfo, error) {
	ctx := context.Background()

	// Query by email field instead of using email as document ID
	docs, err := r.Client.Collection("admin_info").Where("admin_email", "==", email).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("admin not found")
	}

	var admin models.AdminInfo
	docs[0].DataTo(&admin)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &admin, nil
}

// 🟢 Update password
func (r *AdminRepo) UpdatePasswordByEmail(email string, password string) error {
	ctx := context.Background()
	_, err := r.Client.Collection("admins").Doc(email).Set(ctx, map[string]interface{}{"admin_password": password}, firestore.MergeAll)
	return err
}

// 🟢 Delete admin
func (r *AdminRepo) DeleteAdmin(email string) error {
	ctx := context.Background()
	_, err := r.Client.Collection("admins").Doc(email).Delete(ctx)
	return err
}
