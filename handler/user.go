package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	jwtService  auth.Service
}

func NewUserHandler(userService user.Service, jwtService auth.Service) *userHandler {
	return &userHandler{userService, jwtService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Register failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Register failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.jwtService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Register failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)
	response := helper.ApiResponse("success", http.StatusOK, "Account has been registered", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Login failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ApiResponse("error", http.StatusBadRequest, "Register failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.jwtService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.ApiResponse("error", http.StatusBadRequest, "Login failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)
	response := helper.ApiResponse("success", http.StatusOK, "Login success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Email checking failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Something wrong"}

		response := helper.ApiResponse("error", http.StatusUnprocessableEntity, "Email checking failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	msg := "Email has been registered"

	if isEmailAvailable {
		msg = "Email is available"
	}

	response := helper.ApiResponse("success", http.StatusOK, msg, data)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to upload avatar image", data)

		c.JSON(http.StatusBadRequest, response)
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to upload avatar image", data)

		c.JSON(http.StatusBadRequest, response)
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("error", http.StatusBadRequest, "Failed to upload avatar image", data)

		c.JSON(http.StatusBadRequest, response)
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("success", http.StatusOK, "Success upload avatar", data)

	c.JSON(http.StatusOK, response)
}
