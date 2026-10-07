package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ngyewch/expframework/codec"
)

func LoadYAMLConfig[T any, PT *T](configFile string) (PT, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	configCodec := codec.YamlCodec[T, PT]{}
	cfg, err := configCodec.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
