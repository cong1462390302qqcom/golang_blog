package controllers

import (
	"github.com/gin-gonic/gin"
)

type CommentController struct {
}

func (c *CommentController) GetComments(context *gin.Context) {

}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=100"`
}
