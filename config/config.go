package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
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

func InitConfig() *Config {
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
