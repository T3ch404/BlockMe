package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	LogFormat   string
	LogLevel    string
	LogLocation string

	DbType string
	DbHost string
	DbPort int
	DbName string
	DbUser string
	DbPass string

	ResetMe        bool
	ResetMeCron    string
	ResetMeWebhook string

	IgnoreList []string
}

var Env *EnvConfig
var ErrMissingRequiredEnvVar = errors.New("missing required config")

func InitConfig() error {
	fmt.Printf("\nInitializing config...\n")
	_ = godotenv.Load()

	fmt.Println("config: Initializing Logging configs")
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat != "text_pretty" && logFormat != "text" && logFormat != "json" {
		fmt.Println("config: Invalid LOG_FORMAT - Continuing with text_pretty")
		logFormat = "text_pretty"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "DEBUG" && logLevel != "INFO" && logLevel != "WARNING" && logLevel != "ERROR" {
		fmt.Println("config: Invalid LOG_LEVEL - Continuing with INFO")
		logLevel = "INFO"
	}
	logLocation := os.Getenv("LOG_LOCATION")
	if logLocation != "console" && logLocation != "file" {
		fmt.Println("config: Invalid LOG_LOCATION - Continuing with console")
		logLocation = "console"
	}

	fmt.Printf("config: Initializing DB values\n")
	dbType := strings.ToLower(os.Getenv("DB_TYPE"))
	if dbType != "postgres" && dbType != "sqlite" {
		fmt.Println("config: DB_TYPE invalid or not set - Defaulting to sqlite")
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
			fmt.Println("config: Invalid DB_PORT - Continuing with Postgres default 5432")
			dbPort = 5432
			break
		default:
			dbPort = 0
		}
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" && dbType == "postgres" {
		fmt.Println("config: Invalid or missing DB_NAME - Continuing with default 'blockme'")
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
	resetMeStr := os.Getenv("RESET_ME")
	resetMe, err := strconv.ParseBool(resetMeStr)
	if err != nil {
		fmt.Println("config: Invalid or missing RESET_ME - Continuing with default false")
		resetMe = false
	}
	var (
		resetMeCron    string
		resetMeWebhook string
	)

	if resetMe {
		resetMeCron = os.Getenv("RESET_ME_CRON")
		if resetMeCron == "" {
			fmt.Println("config: Invalid or missing RESET_ME_CRON - Continuing with default '1 * * * *")
			resetMeCron = "1 * * * *"
		}
		resetMeWebhook = os.Getenv("RESET_ME_WEBHOOK")
		if resetMeWebhook == "" {
			fmt.Println("config: Missing RESET_ME_WEBHOOK while RESET_ME is enabled - Disabling RESET_ME")
			resetMe = false
		}
	}

	// Remove all chars not in the regex, split on commas, and parse each list item as an IP
	re := regexp.MustCompile("[^0-9.,]")
	ignoreListStr := os.Getenv("IGNORE_LIST")
	ignoreListStr = re.ReplaceAllString(ignoreListStr, "")
	var ignoreList []string
	for _, ignore := range strings.Split(ignoreListStr, ",") {
		addr := net.ParseIP(ignore)
		if addr != nil {
			ignoreList = append(ignoreList, addr.String())
		}
	}

	config := EnvConfig{
		LogFormat:   logFormat,
		LogLevel:    logLevel,
		LogLocation: logLocation,

		DbType: dbType,
		DbHost: dbHost,
		DbPort: dbPort,
		DbName: dbName,
		DbUser: dbUser,
		DbPass: dbPass,

		ResetMe:        resetMe,
		ResetMeCron:    resetMeCron,
		ResetMeWebhook: resetMeWebhook,

		IgnoreList: ignoreList,
	}

	Env = &config

	fmt.Printf("config: Successfully setup config\n")
	return nil
}
