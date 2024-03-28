package config

import (
	"errors"
	"github.com/joho/godotenv"
)

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
