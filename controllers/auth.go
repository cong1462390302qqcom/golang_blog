package controllers

import (
	"github.com/gin-gonic/gin"
	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"
)

type AuthController struct {
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (ac *AuthController) Register(c *gin.Context) {
	var registerRequest = RegisterRequest{}
	err := c.ShouldBindJSON(&registerRequest)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var existingUser models.User
	// 检查用户名是否已存在
	err = config.DB.Where("username=?", registerRequest.Username).First(&existingUser).Error
	if err == nil {
		utils.BadRequest(c, "username already exists")
		return
	}
	// 检查用邮箱是否已存在
	if err = config.DB.Where("email=?", registerRequest.Email).First(existingUser).Error; err == nil {
		utils.BadRequest(c, "email already exists")
		return
	}

	// 创建新用户
	user := models.User{
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
		Password: registerRequest.Password, //在数据落库之前会调用钩子函数，将密码加密
	}
	//落库
	if err = config.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}
	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}
	utils.Success(c, AuthResponse{Token: token, User: user})
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginRequset LoginRequest
	if err := c.ShouldBindJSON(&loginRequset); err != nil {
		utils.BadRequest(c, err.Error())
	}
	//用户名密码验证
	var user models.User
	if err := config.DB.Where("username=?", loginRequset.Username).First(&user).Error; err != nil {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}
	if res := user.CheckPassword(loginRequset.Password); !res {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}
	//生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}
	utils.Success(c, AuthResponse{Token: token, User: user})
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	//获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found ")
		return
	}
	utils.Success(c, user)
}
