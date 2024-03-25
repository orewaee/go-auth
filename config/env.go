package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

var Port string

type env struct {
	key      string
	variable *string
}

var envs = []env{
	{"PORT", &Port},
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
