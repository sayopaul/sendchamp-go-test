package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Debug       string `mapstructure:"DEBUG"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	DBUsername  string `mapstructure:"DB_USERNAME"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBSslMode   string `mapstructure:"DB_SSL_MODE"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBName      string `mapstructure:"DB_NAME"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	AMQPUrl     string `mapstructure:"AMQP_URL"`
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported type %s", v.Type())
	}
	return nil
}

func (c *Config) PopulateFieldsForProduction() {
	v := reflect.ValueOf(c).Elem()

	// Build map of fields keyed by the tagName
	fields := make(map[string]reflect.Value)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tagName := fieldInfo.Tag.Get("mapstructure")
		fields[tagName] = v.Field(i)
	}

	// Update config field for each parameter
	for k, v := range fields {
		if !v.IsValid() {
			continue // ignore unrecognized config values
		}
		//populate config struct
		err := populate(v, os.Getenv(k))
		if err != nil {
			log.Fatal("error populating config struct field", err)
		}
	}

}

// returns an instance of Config, for access to env variables app-wide
func LoadConfig() Config {
	config := Config{}

	environment := os.Getenv("ENVIRONMENT")
	//if local environment, use local.env, if production, use the OS' environment variables (usually injected in by the cloud provider e.g using Docker)
	if environment == "LOCAL" || environment == "" {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err, viper.ConfigFileUsed())
		}
		//unmarshal env to cinfig struct
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatal("environment cant be loaded: ", err)
		}
	} else {
		//we are in production, do not use .env, use the OS' environment variables
		config.PopulateFieldsForProduction()
	}
	return config
}
