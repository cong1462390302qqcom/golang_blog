package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang_blog/config"
	"golang_blog/routes"
	"log"
	"os"
)

func main() {
	err := godotenv.Load() //获取环境变量
	if err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.Infof("Starting blog application...")
	config.InitDatabase() //初始化数据库

	engine := routes.SetupRoutes() //设置路由
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	err = engine.Run(port)
	logrus.WithField("port", port).Info("Server starting")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
