package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/orewaee/go-auth/config"
	"github.com/orewaee/go-auth/database"
	"github.com/orewaee/go-auth/dto"
	"github.com/orewaee/go-auth/models"
	"github.com/orewaee/go-auth/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func SignIn(ctx *fiber.Ctx) error {
	var body dto.SignInBody
	if err := ctx.BodyParser(&body); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(dto.Error{Message: err.Error()})
	}

	if err := body.Validate(); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(dto.Error{Message: err.Error()})
	}

	users := database.GetCollection("users")
	filter := bson.D{{"email", body.Email}}

	var user models.User
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(dto.Error{Message: "user not found"})
	}

	if !user.Activated {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return ctx.JSON(dto.Error{Message: "account not activated"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.JSON(dto.Error{Message: "wrong password"})
	}

	accessToken := token.GenerateToken(jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}, config.AccessSecret)

	refreshToken := token.GenerateToken(jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	}, config.RefreshSecret)

	sessions := database.GetCollection("sessions")

	newSession := models.Session{
		Id:           primitive.NewObjectID(),
		User:         user.Id,
		Ip:           ctx.IP(),
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
	}

	_, err = sessions.InsertOne(context.TODO(), newSession)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(dto.Error{Message: err.Error()})
	}

	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Minute * 60),
		HTTPOnly: true,
	}

	ctx.Cookie(cookie)
	return ctx.JSON(dto.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken})
}
