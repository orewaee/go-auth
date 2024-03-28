package handlers

import "github.com/gofiber/fiber/v2"

func Ping(ctx *fiber.Ctx) error {
	return ctx.SendString("pong")
}
