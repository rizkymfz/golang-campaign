package transaction

import "github.com/rizkymfz/golang-campaign/user"

//cocok kalo param lebih dari satu
type GetCampaignTransactionInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
