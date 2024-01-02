package config

import (
	"errors"
	"os"
	"reflect"
)

type Config struct {
	Provider     string `env:"PROVIDER"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	Scopes       string `env:"SCOPES"`
}

var ErrMissingEnvTag = errors.New("missing 'env' tag in config field")
var ErrUnsupportedFieldType = errors.New("unsupported config field type")

func LoadConfigFromEnv(conf *Config) error {
	t := reflect.TypeOf(conf).Elem()
	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldT := t.Field(i)
		fieldV := v.Field(i)
		key, ok := fieldT.Tag.Lookup("env")
		if !ok {
			return ErrMissingEnvTag
		}

		confV, ok := os.LookupEnv(key)
		kind := fieldV.Kind()
		if ok {
			switch kind {
			case reflect.String:
				fieldV.SetString(confV)
			default:
				return ErrUnsupportedFieldType
			}
		}
	}
	return nil
}
