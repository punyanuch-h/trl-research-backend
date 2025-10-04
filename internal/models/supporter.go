package models

type Supporter struct {
	ID                              string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CaseID                          string `gorm:"type:varchar(10);not null"`
	SupportResearch                 bool   `gorm:"not null"`
	SupportVDC                      bool   `gorm:"not null"`
	SupportSiEIC                    bool   `gorm:"not null"`
	NeedProtectIntellectualProperty bool   `gorm:"not null"`
	NeedCoDevelopers                bool   `gorm:"not null"`
	NeedActivities                  bool   `gorm:"not null"`
	NeedTest                        bool   `gorm:"not null"`
	NeedCapital                     bool   `gorm:"not null"`
	NeedPartners                    bool   `gorm:"not null"`
	NeedGuidelines                  bool   `gorm:"not null"`
	NeedCertification               bool   `gorm:"not null"`
	NeedAccount                     bool   `gorm:"not null"`
	Need                            string `gorm:"type:text"`
	AdditionalDocuments             string `gorm:"type:text"`
}
