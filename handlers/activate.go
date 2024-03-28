package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Activate(ctx *fiber.Ctx) error {
	secret := ctx.Params("secret")

	activations := database.GetCollection("activations")

	var activation models.Activation
	err := activations.FindOne(context.TODO(), bson.D{{"secret", secret}}).Decode(&activation)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return fiber.NewError(fiber.StatusNotFound, "Activation error")
	}

	users := database.GetCollection("users")
	filter := bson.D{{"_id", activation.User}}
	update := bson.D{{"$set", bson.D{{"activated", true}}}}

	var user models.User
	err = users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	_, err = users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Activation error")
	}

	_, err = activations.DeleteOne(context.TODO(), bson.D{{"_id", activation.Id}})
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Activation error")
	}

	return ctx.SendString("Account activated")
}
