package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/orewaee/go-auth/config"
	"github.com/orewaee/go-auth/controllers"
	"log"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	log.Println(config.Port)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST",
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	app.Get("/ping", controllers.Ping)

	log.Fatalln(app.Listen(":" + config.Port))
}
