package model

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Dsn      string `yaml:"dsn"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

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

func initConfig() *Config {
	// 读取YAML文件
	wd, _ := os.Getwd()
	configPath := filepath.Join(wd, "config", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	// 解析YAML到结构体
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("解析YAML失败: %v", err)
	}

	// 使用配置
	fmt.Printf("服务器地址: %s:%d\n",
		config.Server.Host, config.Server.Port)
	fmt.Printf("数据库用户: %s\n", config.Database.Username)
	fmt.Println("dsn: ", config.Database.Dsn)
	return &config

}

func InitDB() (*gorm.DB, *Config) {
	config := initConfig()
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
