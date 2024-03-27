package controllers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/orewaee/go-auth/config"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/models"
	"github.com/orewaee/go-auth/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type SignInBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func SignIn(ctx *fiber.Ctx) error {
	var body SignInBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	users := database.GetCollection("users")
	filter := bson.D{{"email", body.Email}}

	var user models.User
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if !user.Activated {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "Account not activated")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Wrong password")
	}

	accessToken := token.GenerateToken(jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}, config.AccessSecret)

	refreshToken := token.GenerateToken(jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	}, config.RefreshSecret)

	newSession := models.Session{
		Id:           primitive.NewObjectID(),
		User:         user.Id,
		Ip:           ctx.IP(),
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
	}

	sessions := database.GetCollection("sessions")

	_, err = sessions.InsertOne(context.TODO(), newSession)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Minute * 60),
		HTTPOnly: true,
	}

	ctx.Cookie(cookie)
	return ctx.JSON(TokenPair{accessToken, refreshToken})
}
