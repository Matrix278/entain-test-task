package configuration

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDBHost                  string
	PostgresDBPort                  string
	PostgresDBUser                  string
	PostgresDBPass                  string
	PostgresDBName                  string
	PostgresDatabaseURL             string
	CancelOddRecordsMinutesInterval int
	ServerPort                      string
}

func Load() (config Config, err error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	viper.AutomaticEnv()

	config.PostgresDBHost = viper.GetString("POSTGRES_DB_HOST")
	config.PostgresDBPort = viper.GetString("POSTGRES_DB_PORT")
	config.PostgresDBUser = viper.GetString("POSTGRES_DB_USER")
	config.PostgresDBPass = viper.GetString("POSTGRES_DB_PASS")
	config.PostgresDBName = viper.GetString("POSTGRES_DB_NAME")
	config.PostgresDatabaseURL = viper.GetString("POSTGRES_DATABASEURL")
	config.CancelOddRecordsMinutesInterval = viper.GetInt("CANCEL_ODD_RECORDS_MINUTES_INTERVAL")
	config.ServerPort = viper.GetString("SERVER_PORT")

	return
}
