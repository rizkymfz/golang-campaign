package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserId(userID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.
		Where("campaign_id = ?", campaignID).
		Preload("User").
		Find(&transactions).
		Order("id desc").Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByUserId(userID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.
		Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").
		Preload("User").
		Where("user_id = ?", userID).
		Find(&transactions).
		Order("id desc").Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
