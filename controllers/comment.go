package controllers

import (
	"github.com/gin-gonic/gin"
	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"
	"strconv"
)

type CommentController struct {
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=100"`
}

func (cc *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var comments []models.Comment
	if err = config.DB.Preload("User").Preload("Post").Find(&comments, "post_id=?", postID).Error; err != nil {
		utils.NotFound(c, "comment not found")
		return
	}
	utils.Success(c, comments)
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}
	var commentRequest CreateCommentRequest
	if err = c.ShouldBindJSON(&commentRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	//评论的文章是否存在
	var post models.Post
	if err = config.DB.Where("id=?", postId).First(&post).Error; err != nil {
		utils.NotFound(c, "post not found")
		return
	}
	comment := models.Comment{
		Content: commentRequest.Content,
		UserID:  userID.(uint),
		PostID:  uint(postId),
	}
	if err = config.DB.Save(&comment).Error; err != nil {
		utils.InternalServerError(c, "Failed to create comment")
	}
	config.DB.Preload("User").First(&comment, comment.ID)
	utils.Success(c, comment)
}
