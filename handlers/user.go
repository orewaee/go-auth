package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/dto"
	"github.com/orewaee/go-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func User(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(dto.Error{Message: err.Error()})
	}

	users := database.GetCollection("users")
	filter := bson.D{{"_id", objectId}}

	var user models.User
	err = users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(dto.Error{Message: "user not found"})
	}

	return ctx.JSON(user)
}
