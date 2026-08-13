package controllers

import (
	"github.com/gin-gonic/gin"
	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"
	"strconv"
)

type PostController struct {
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var createPost CreatePostRequest
	err := c.ShouldBindJSON(&createPost)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	post := models.Post{
		Title:   createPost.Title,
		Content: createPost.Content,
		UserID:  userID.(uint),
	}
	if err = config.DB.Create(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to create post")
		return
	}
	config.DB.Preload("User").First(&post, post.ID)
	utils.Success(c, post)
}
func (pc *PostController) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}
	var updateRequest UpdatePostRequest
	if err = c.ShouldBindJSON(&updateRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	var post models.Post
	err = config.DB.First(&post, postID).Error
	if err != nil {
		utils.NotFound(c, "Post not found")
		return
	}
	if userID != post.UserID {
		utils.Forbidden(c, "You can only update your own posts")
		return
	}
	//更新文章内容
	post.Title = updateRequest.Title
	post.Content = updateRequest.Content
	if err = config.DB.UpdateColumns(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to update post")
		return
	}
	config.DB.Preload("User").First(&post, postID)
	utils.Success(c, post)
}
func (pc *PostController) DeletePost(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	var post models.Post
	if err = config.DB.First(&post, postId).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}
	if userID != post.UserID {
		utils.Forbidden(c, "You can only delete your own posts")
		return
	}
	//软删除
	if err = config.DB.Delete(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete post "+err.Error())
		return
	}
	if err = config.DB.Delete("Comment", "post_id=?", postId).Error; err != nil {
		utils.InternalServerError(c, "Filed to delete comment "+err.Error())
		return
	}
}
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post

	//分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	config.DB.Preload("User").Order("created_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&posts)

	//获取数据总量
	var total int64
	config.DB.Model(&models.Post{}).Count(&total)
	utils.Success(c, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"posts":     posts,
	})
}

func (pc *PostController) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}
	var post models.Post
	if err = config.DB.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}
	utils.Success(c, post)
}
