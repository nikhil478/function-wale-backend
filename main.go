package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nikhil478/function-wale-server/config"
	_ "github.com/nikhil478/function-wale-server/docs"
	"github.com/nikhil478/function-wale-server/models"
	router "github.com/nikhil478/function-wale-server/routers"
	"github.com/nikhil478/function-wale-server/storage"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal("could not load config: ", err)
	}

	dbConfig := &storage.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	awsConfig := &storage.AwsS3Config{
		Region: cfg.AWS.S3.Region,
		Id:     cfg.AWS.S3.Id,
		Secret: cfg.AWS.S3.Secret,
		Token:  cfg.AWS.S3.Token,
	}

	awsS3, err := storage.NewAwsConnection(awsConfig)
	if err != nil {
		log.Fatal("could not load AWS S3")
	}

	db, err := storage.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("could not load the db")
	}

	err = models.MigrateOrganization(db)
	if err != nil {
		log.Fatal("could not migrate organization table")
	}

	err = models.MigrateFiles(db)
	if err != nil {
		log.Fatal("could not migrate files table")
	}

	err = models.MigrateUser(db)
	if err != nil {
		log.Fatal("could not migrate files table")
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	router.SetupRoutes(app, db, awsS3)

	app.Listen(":8000")
}
