package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type SignUpBody struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func SignUp(ctx *fiber.Ctx) error {
	var body SignUpBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	users := database.GetCollection("users")
	filter := bson.D{{"email", body.Email}}

	var user models.User
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err == nil {
		return fiber.NewError(fiber.StatusConflict, "User with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 8)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating hash")
	}

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Email:     body.Email,
		Name:      body.Name,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	if _, err := users.InsertOne(context.TODO(), newUser); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(newUser)
}
