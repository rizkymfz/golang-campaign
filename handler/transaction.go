package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("failed to get transaction campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	data, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		msg := err.Error()
		if msg == "" {
			msg = "Failed get detail campaign"
		}
		response := helper.ApiResponse("error", http.StatusBadRequest, msg, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "transaction detail", transaction.FormatCampaignTransactions(data))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		msg := err.Error()
		if msg == "" {
			msg = "Failed get user transactions"
		}
		response := helper.ApiResponse("error", http.StatusBadRequest, msg, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "user transactions", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}
