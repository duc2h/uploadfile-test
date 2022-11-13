package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GoogleCredentialPath string `mapstructure:"GOOGLE_CREDENTIAL_PATH"`
	GoogleBucket         string `mapstructure:"GOOGLE_BUCKET"`
	GoogleProjectID      string `mapstructure:"GOOGLE_PROJECT_ID"`
	ObjectPath           string `mapstructure:"OBJECT_PATH"`

	AwsAccessKey string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey string `mapstructure:"AWS_SECRET_KEY"`
	AwsToken     string `mapstructure:"AWS_TOKEN"`
	AwsBucket    string `mapstructure:"AWS_BUCKET"`
	AwsRegion    string `mapstructure:"AWS_REGION"`
	AwsKey       string `mapstructure:"AWS_KEY"`
	AwsACL       string `mapstructure:"AWS_ACL"`

	Timeout time.Duration `mapstructure:"TIME_OUT"`
}

func LoadConfig(path string, env string, ref interface{}) error {
	viper.AddConfigPath(path)
	name := "app." + env
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&ref)
	if err != nil {
		return err
	}
	return nil
}
