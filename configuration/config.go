package configuration

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDBHost                     string
	PostgresDBPort                     string
	PostgresDBUser                     string
	PostgresDBPass                     string
	PostgresDBName                     string
	PostgresDatabaseURL                string
	CancelOddRecordsMinutesInterval    int
	NumberOfLatestRecordsForCancelling int
	ServerPort                         string
}

func Load() (config *Config, err error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	return &Config{
		PostgresDBHost:                     viper.GetString("POSTGRES_DB_HOST"),
		PostgresDBPort:                     viper.GetString("POSTGRES_DB_PORT"),
		PostgresDBUser:                     viper.GetString("POSTGRES_DB_USER"),
		PostgresDBPass:                     viper.GetString("POSTGRES_DB_PASS"),
		PostgresDBName:                     viper.GetString("POSTGRES_DB_NAME"),
		PostgresDatabaseURL:                viper.GetString("POSTGRES_DATABASEURL"),
		CancelOddRecordsMinutesInterval:    viper.GetInt("CANCEL_ODD_RECORDS_MINUTES_INTERVAL"),
		NumberOfLatestRecordsForCancelling: viper.GetInt("NUMBER_OF_LATEST_RECORDS_FOR_CANCELLING"),
		ServerPort:                         viper.GetString("SERVER_PORT"),
	}, nil
}
