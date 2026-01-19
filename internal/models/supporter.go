package models

import "time"

type Supporter struct {
	SupporterID                     string    `gorm:"primaryKey;column:supporter_id" json:"supporter_id"`
	CaseID                          string    `gorm:"column:case_id" json:"case_id"`
	SupportResearch                 bool      `gorm:"column:support_research" json:"support_research"`
	SupportVDC                      bool      `gorm:"column:support_vdc" json:"support_vdc"`
	SupportSiEIC                    bool      `gorm:"column:support_sieic" json:"support_sieic"`
	NeedProtectIntellectualProperty bool      `gorm:"column:need_protect_intellectual_property" json:"need_protect_intellectual_property"`
	NeedCoDevelopers                bool      `gorm:"column:need_co_developers" json:"need_co_developers"`
	NeedActivities                  bool      `gorm:"column:need_activities" json:"need_activities"`
	NeedTest                        bool      `gorm:"column:need_test" json:"need_test"`
	NeedCapital                     bool      `gorm:"column:need_capital" json:"need_capital"`
	NeedPartners                    bool      `gorm:"column:need_partners" json:"need_partners"`
	NeedGuidelines                  bool      `gorm:"column:need_guidelines" json:"need_guidelines"`
	NeedCertification               bool      `gorm:"column:need_certification" json:"need_certification"`
	NeedAccount                     bool      `gorm:"column:need_account" json:"need_account"`
	Need                            string    `gorm:"column:need" json:"need"`
	CreatedAt                       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Supporter) TableName() string {
	return "supporters"
}
