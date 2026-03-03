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
	DbPort     int
	DbName     string
	DbUser     string
	DbPassword string

	EnableCron bool

	IAmALittleBitch     bool
	IAmALittleBitchCron string
	IAmALittleBitchUrl  string
}

var Env *EnvConfig

func initEnv() error {
	fmt.Printf("config: Initializing env\n")
	err := godotenv.Load()
	if err != nil {
		return err
	}

	fmt.Printf("config: Building singleton\n")
	config := EnvConfig{}

	fmt.Printf("config: Initializing DB values\n")
	config.EnableDB = getEnvBool("ENABLE_DB", false, false)
	if config.EnableDB {
		config.DbType = getEnvString("DB_TYPE", false, "")
		config.DbHost = getEnvString("DB_HOST", false, "")
		config.DbPort = getEnvInt("DB_PORT", false, -1)
		config.DbName = getEnvString("DB_NAME", false, "")
		config.DbUser = getEnvString("DB_USER", false, "")
		config.DbPassword = getEnvString("DB_PASSWORD", false, "")
	} else {
		fmt.Printf("config: DB disabled\n")
	}

	fmt.Printf("config: Initializing Cron values\n")
	config.EnableCron = getEnvBool("ENABLE_CRON", false, false)

	if config.EnableCron {
		fmt.Printf("config: Initializing reset values\n")
		iAmALittleBitch := getEnvBool("I_AM_A_LITTLE_BITCH", false, false)
		iAmALittleBitchCron := getEnvString("I_AM_A_LITTLE_BITCH_CRON", false, "")
		if iAmALittleBitchCron == "" {
			iAmALittleBitchCron = "1 * * * *"
		}
		iAmALittleBitchUrl := getEnvString("I_AM_A_LITTLE_BITCH_WEBHOOK_URL", false, "")
		if iAmALittleBitch && iAmALittleBitchUrl == "" {
			fmt.Println("config: WARNING, reset has been enabled without a webhook for key rotation")
		}
	} else {
		fmt.Printf("config: Cron disabled\n")
	}

	Env = &config

	fmt.Printf("config: Env initialized\n")
	return nil
}

func ReloadConfig() error {
	fmt.Printf("\nReloading config...\n")
	Env = nil
	return initEnv()
}

func InitConfig() error {
	fmt.Printf("\nInitializing config...\n")
	if Env != nil {
		fmt.Printf("\nEnv config already initialized. Skipping.\n")
		return nil
	}
	return initEnv()
}

func getEnvString(key string, required bool, defaultValue string) string {
	strVal := os.Getenv(key)
	if strVal == "" {
		if required {
			panic(fmt.Sprintf("ERROR: Missing required environment variable: %s", key))
		}
		return defaultValue
	}
	return strVal
}

func getEnvInt(key string, required bool, defaultValue int) int {
	strVal := os.Getenv(key)
	if strVal == "" {
		if required {
			panic(fmt.Sprintf("ERROR: Missing required environment variable: %s", key))
		}
		return defaultValue
	}

	i, err := strconv.Atoi(strVal)
	if err != nil {
		fmt.Printf("config: WARNING, env var %s is not an integer. Setting to -1\n", key)
		return -1
	}
	return i
}

func getEnvBool(key string, required, defaultValue bool) bool {
	strVal := os.Getenv(key)
	if strVal == "" {
		if required {
			panic(fmt.Sprintf("ERROR: Missing required environment variable: %s", key))
		}
		return defaultValue
	}

	bVal, err := strconv.ParseBool(strVal)
	if err != nil {
		fmt.Printf("config: WARNING, env var %s is not a boolean. Setting to false\n", key)
		return false
	}
	return bVal
}
