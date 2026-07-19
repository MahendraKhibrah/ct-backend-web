package main

import (
	"ct-backend/Config"
	"ct-backend/Route"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//if loc, err := time.LoadLocation("Asia/Jakarta"); err != nil {
	//	panic(err)
	//} else {
	//	time.Local = loc
	//}

	db := Config.SetUpDatabaseConnection()

	if db == nil {
		panic("Failed to connect to relational database!")
	}
	defer Config.CloseDatabaseConnection(db)

	println("Connected to relational database")

	minioClient := Config.SetupMinioConnection()
	if minioClient == nil {
		panic("Failed to connect to minio!")
	}

	println("Connected to minio")

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Origin,Content-Type,Accept,User-Agent,Content-Length,Authorization"},
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"},
	}))

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = "0.0.0.0:" + port
	}

	// init route and DI
	Route.Init(server, db)

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
