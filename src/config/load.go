package config

import (
	"errors"
	"flag"
	"io/fs"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

const defaultFilePath = "./config/config.yml"

func LoadConfig(cfg interface{}) error {
	filePath := flag.String("config-file", defaultFilePath, "The config file to be used")
	flag.Parse()

	if err := cleanenv.ReadConfig(*filePath, cfg); err != nil {
		if *filePath == defaultFilePath && errors.Is(err, fs.ErrNotExist) {
			slog.Debug("Skipping missing config file", "error", err)

			if err := cleanenv.ReadEnv(cfg); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(cfg)
}
