package handler

import (
	"net/http"

	"github.com/rizkymfz/golang-campaign/helper"
	"github.com/rizkymfz/golang-campaign/transaction"
	"github.com/rizkymfz/golang-campaign/user"

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

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Failed to create campaign", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to create transaction", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "Success create transaction", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Failed to process notification", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	err = h.service.ProcessPayment(input)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Failed to process notification", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	c.JSON(http.StatusOK, input)
}
