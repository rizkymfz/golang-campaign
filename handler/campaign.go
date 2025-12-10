package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {

		response := helper.ApiResponse("error", http.StatusBadRequest, "error get campaigns", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "list of campaigns", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaignByID(c *gin.Context) {
	idStr := c.Param("id")
	campaignID, err := strconv.Atoi(idStr)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Invalid campaign ID", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data, err := h.service.GetCampaignByID(campaignID)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed get detail campaign", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "list of campaigns", campaign.FormatCampaignDetail(data))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

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

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to create campaign", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "Success create campaign", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	campaignID, err := strconv.Atoi(idStr)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Invalid campaign ID", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input campaign.UpdateCampaignInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Failed to update campaign", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(campaignID, input)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("success", http.StatusOK, "Success create campaign", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	input.User = currentUser

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Failed upload image", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to upload image", data)

		c.JSON(http.StatusBadRequest, response)
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to upload image", data)

		c.JSON(http.StatusBadRequest, response)
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to upload image", data)

		c.JSON(http.StatusBadRequest, response)
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("success", http.StatusOK, "Success upload image", data)

	c.JSON(http.StatusOK, response)
}
