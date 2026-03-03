package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	EnableDB   bool
	DbType     string
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string

	EnableCron bool

	IAmALittleBitch     bool
	IAmALittleBitchCron string
	IAmALittleBitchUrl  string
}

var Env *EnvConfig

func InitConfig() {
	fmt.Printf("\nInitializing config...\n")
	godotenv.Load()

	fmt.Printf("config: Initializing DB values\n")
	enableDbStr := os.Getenv("ENABLE_DB")
	enableDb, err := strconv.ParseBool(enableDbStr)
	if err != nil {
		enableDb = false
	}
	dbType := os.Getenv("DB_TYPE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	fmt.Printf("config: Initializing Cron values\n")
	enableCronStr := os.Getenv("ENABLE_CRON")
	enableCron, err := strconv.ParseBool(enableCronStr)
	if err != nil {
		enableCron = false
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
		EnableDB:   enableDb,
		DbType:     dbType,
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
		DbUser:     dbUser,
		DbPassword: dbPassword,

		EnableCron: enableCron,

		IAmALittleBitch:     iAmALittleBitch,
		IAmALittleBitchCron: iAmALittleBitchCron,
		IAmALittleBitchUrl:  iAmALittleBitchUrl,
	}

	Env = &config

	fmt.Printf("config: Env initialized\n")
}
