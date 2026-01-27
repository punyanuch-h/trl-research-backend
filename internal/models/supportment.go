package models

import "time"

type Supportments struct {
	ID                     			string    `gorm:"primaryKey;column:id" json:"id"`
	CaseID                          string    `gorm:"column:case_id;not null" json:"case_id"`
	Case 						    *Cases    `gorm:"foreignKey:CaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"case"`
	SupportResearch                 bool      `gorm:"column:support_research;not null" json:"support_research"`
	SupportVDC                      bool      `gorm:"column:support_vdc;not null" json:"support_vdc"`
	SupportSiEIC                    bool      `gorm:"column:support_sieic;not null" json:"support_sieic"`
	NeedProtectIntellectualProperty bool      `gorm:"column:need_protect_intellectual_property;not null" json:"need_protect_intellectual_property"`
	NeedCoDevelopers                bool      `gorm:"column:need_co_developers;not null" json:"need_co_developers"`
	NeedActivities                  bool      `gorm:"column:need_activities;not null" json:"need_activities"`
	NeedTest                        bool      `gorm:"column:need_test;not null" json:"need_test"`
	NeedCapital                     bool      `gorm:"column:need_capital;not null" json:"need_capital"`
	NeedPartners                    bool      `gorm:"column:need_partners;not null" json:"need_partners"`
	NeedGuidelines                  bool      `gorm:"column:need_guidelines;not null" json:"need_guidelines"`
	NeedCertification               bool      `gorm:"column:need_certification;not null" json:"need_certification"`
	NeedAccount                     bool      `gorm:"column:need_account;not null" json:"need_account"`
	Need                            string    `gorm:"column:need" json:"need"`
	CreatedAt                       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Supportments) TableName() string {
	return "supportments"
}
