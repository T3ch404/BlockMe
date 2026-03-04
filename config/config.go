package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DbType string
	DbHost string
	DbPort int
	DbName string
	DbUser string
	DbPass string

	IAmALittleBitch     bool
	IAmALittleBitchCron string
	IAmALittleBitchUrl  string
}

var Env *EnvConfig
var ErrMissingRequiredEnvVar = errors.New("missing required config")

func InitConfig() error {
	fmt.Printf("\nInitializing config...\n")
	godotenv.Load()

	fmt.Printf("config: Initializing DB values\n")
	dbType := strings.ToLower(os.Getenv("DB_TYPE"))
	if dbType != "postgres" && dbType != "sqlite" {
		fmt.Println("DB_TYPE invalid or not set - Defaulting to sqlite")
		dbType = "sqlite"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" && dbType == "postgres" {
		return fmt.Errorf("%w: %s", ErrMissingRequiredEnvVar, "DB_HOST")
	}
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" && dbType == "postgres" {
		dbPortStr = "0"
	}
	dbPort, _ := strconv.Atoi(dbPortStr)
	if dbPort <= 0 || dbPort > 65535 {
		switch dbType {
		case "postgres":
			fmt.Println("Invalid DB_PORT - Continuing with Postgres default 5432")
			dbPort = 5432
			break
		default:
			dbPort = 0
		}
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" && dbType == "postgres" {
		fmt.Println("WARN: Invalid or missing DB_NAME - Continuing with default 'blockme'")
		dbName = "blockme"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" && dbType == "postgres" {
		return fmt.Errorf("%w: %s", ErrMissingRequiredEnvVar, "DB_USER")
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" && dbType == "postgres" {
		return fmt.Errorf("%w: %s", ErrMissingRequiredEnvVar, "DB_PASS")
	}

	fmt.Printf("config: Initializing reset values\n")
	iAmALittleBitchStr := os.Getenv("I_AM_A_LITTLE_BITCH")
	iAmALittleBitch, err := strconv.ParseBool(iAmALittleBitchStr)
	if err != nil {
		iAmALittleBitch = false
	}
	iAmALittleBitchCron := os.Getenv("I_AM_A_LITTLE_BITCH_CRON")
	if iAmALittleBitchCron == "" {
		iAmALittleBitchCron = "1 * * * *"
	}
	iAmALittleBitchUrl := os.Getenv("I_AM_A_LITTLE_BITCH_WEBHOOK_URL")
	if iAmALittleBitch && iAmALittleBitchUrl == "" {
		fmt.Println("config: WARNING, reset has been enabled without a webhook for key rotation")
	}

	config := EnvConfig{
		DbType: dbType,
		DbHost: dbHost,
		DbPort: dbPort,
		DbName: dbName,
		DbUser: dbUser,
		DbPass: dbPass,

		IAmALittleBitch:     iAmALittleBitch,
		IAmALittleBitchCron: iAmALittleBitchCron,
		IAmALittleBitchUrl:  iAmALittleBitchUrl,
	}

	Env = &config

	fmt.Printf("config: Env initialized\n")
	return nil
}
