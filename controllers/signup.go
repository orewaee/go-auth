package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/orewaee/go-auth/activation"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/email"
	"github.com/orewaee/go-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
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

	go func() {
		time.AfterFunc(5*time.Minute, func() {
			var user models.User
			_ = users.FindOne(context.TODO(), filter).Decode(&user)

			if !user.Activated {
				_, _ = users.DeleteOne(context.TODO(), filter)
			}
		})
	}()

	activations := database.GetCollection("activations")

	newActivation := models.Activation{
		Id:        primitive.NewObjectID(),
		User:      newUser.Id,
		Secret:    activation.GenerateSecret(16),
		CreatedAt: time.Now(),
	}

	if _, err := activations.InsertOne(context.Background(), newActivation); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	go func() {
		err := email.SendMail(
			[]string{newUser.Email},
			"Activate account",
			fmt.Sprintf(
				"<p>Activate your account by clicking on the <a style='color:#16d886' href='http://localhost:8080/activate/%s'>link</a>."+
					"It will be valid for 5 minutes.</p><br><p><b>%s (%s)</b></p>",
				newActivation.Secret, ctx.GetReqHeaders()["User-Agent"][0], ctx.IP(),
			),
		)

		if err != nil {
			log.Println(err)
		}
	}()

	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(newUser)
}
