package handlers

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
	"time"
)

func Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	ctx.ClearCookie("refresh_token")

	if !token.VerifyToken(refreshToken, config.RefreshSecret) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	sessions := database.GetCollection("sessions")
	filter := bson.D{{"refresh_token", refreshToken}}

	var session models.Session
	err := sessions.FindOne(context.TODO(), filter).Decode(&session)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return fiber.NewError(fiber.StatusUnauthorized, "Session not found")
	}

	filter = bson.D{{"_id", session.Id}}

	if _, err := sessions.DeleteOne(context.TODO(), filter); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	newAccessToken := token.GenerateToken(jwt.MapClaims{
		"id":  session.User,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}, config.AccessSecret)

	newRefreshToken := token.GenerateToken(jwt.MapClaims{
		"id":  session.User,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	}, config.RefreshSecret)

	newSession := models.Session{
		Id:           primitive.NewObjectID(),
		User:         session.User,
		Ip:           ctx.IP(),
		RefreshToken: newRefreshToken,
		CreatedAt:    time.Now(),
	}

	if _, err := sessions.InsertOne(context.TODO(), newSession); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(time.Minute * 60),
		HTTPOnly: true,
	}

	ctx.Cookie(cookie)
	return ctx.JSON(TokenPair{newAccessToken, newRefreshToken})
}
