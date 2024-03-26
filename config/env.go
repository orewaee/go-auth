package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

var Port string
var MongoUri string
var AccessSecret string
var RefreshSecret string

type env struct {
	key      string
	variable *string
}

var envs = []env{
	{"PORT", &Port},
	{"MONGO_URI", &MongoUri},
	{"ACCESS_SECRET", &AccessSecret},
	{"REFRESH_SECRET", &RefreshSecret},
}

func loadEnv(env env) bool {
	value, ok := os.LookupEnv(env.key)
	if ok {
		*env.variable = value
	}
	return ok
}

func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	for _, env := range envs {
		if !loadEnv(env) {
			return errors.New(env.key + " does not exist")
		}
	}

	return nil
}
