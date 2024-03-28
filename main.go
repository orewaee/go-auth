package main

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/orewaee/go-auth/config"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/handlers"
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

	app.Get("/ping", handlers.Ping)

	app.Post("/signup", handlers.SignUp)
	app.Post("/signin", handlers.SignIn)
	app.Post("/refresh", handlers.Refresh)
	app.Get("/activate/:secret", handlers.Activate)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AccessSecret)},
	}))

	app.Get("/user", handlers.User)

	log.Fatalln(app.Listen(":" + config.Port))
}
