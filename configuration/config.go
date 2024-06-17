package configuration

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost                       string
	PostgresPort                       string
	PostgresUser                       string
	PostgresPassword                   string
	PostgresName                       string
	CancelOddRecordsMinutesInterval    int
	NumberOfLatestRecordsForCancelling int
	ServerPort                         string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	return &Config{
		PostgresHost:                       viper.GetString("POSTGRES_HOST"),
		PostgresPort:                       viper.GetString("POSTGRES_PORT"),
		PostgresUser:                       viper.GetString("POSTGRES_USER"),
		PostgresPassword:                   viper.GetString("POSTGRES_PASSWORD"),
		PostgresName:                       viper.GetString("POSTGRES_DB"),
		CancelOddRecordsMinutesInterval:    viper.GetInt("CANCEL_ODD_RECORDS_MINUTES_INTERVAL"),
		NumberOfLatestRecordsForCancelling: viper.GetInt("NUMBER_OF_LATEST_RECORDS_FOR_CANCELLING"),
		ServerPort:                         viper.GetString("SERVER_PORT"),
	}, nil
}
