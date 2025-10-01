package model

import (
	"github.com/golang-jwt/jwt/v5"
	. "go_task4/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

var (
	db     *gorm.DB
	JwtKey = []byte("your_secret_key")
)

func InitDB() (*gorm.DB, *Config) {
	config := InitConfig()
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dsn := config.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("数据库连接失败")
	}

	// 自动迁移
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	return db, config
}
