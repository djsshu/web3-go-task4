package handler

import (
	"github.com/gin-gonic/gin"
	. "go_task4/model"
	"gorm.io/gorm"
	"net/http"
)

// 核心业务逻辑实现

func CreateComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)
	postID := c.Param("id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	comment := Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  post.ID,
	}

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func ListComments(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	postID := c.Param("id")
	var comments []Comment
	if err := db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
