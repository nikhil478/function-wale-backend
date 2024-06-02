package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	AWS struct {
		S3 struct {
			Region string `mapstructure:"region"`
			Id     string `mapstructure:"id"`
			Secret string `mapstructure:"secret"`
			Token  string `mapstructure:"token"`
		} `mapstructure:"s3"`
	} `mapstructure:"aws"`
}


func LoadConfig(path string) (*Config, error) {
	var config Config

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")   
	viper.AddConfigPath(".")       
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}