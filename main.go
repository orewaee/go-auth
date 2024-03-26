package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/orewaee/go-auth/config"
	"github.com/orewaee/go-auth/controllers"
	"github.com/orewaee/go-auth/database"
	"log"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	if err := database.Load(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := database.Unload(); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST",
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	app.Get("/ping", controllers.Ping)

	app.Post("/signup", controllers.SignUp)
	app.Post("/signin", controllers.SignIn)
	app.Post("/refresh", controllers.Refresh)

	log.Fatalln(app.Listen(":" + config.Port))
}
