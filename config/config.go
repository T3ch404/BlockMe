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
	DbType string
	DbHost string
	DbPort int
	DbName string
	DbUser string
	DbPass string

	IAmALittleBitch     bool
	IAmALittleBitchCron string
	IAmALittleBitchUrl  string

	IgnoreList []string
}

var Env *EnvConfig
var ErrMissingRequiredEnvVar = errors.New("missing required config")

func InitConfig() error {
	fmt.Printf("\nInitializing config...\n")
	_ = godotenv.Load()

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
	iAmALittleBitchStr := os.Getenv("I_AM_A_LITTLE_BITCH")
	iAmALittleBitch, err := strconv.ParseBool(iAmALittleBitchStr)
	if err != nil {
		fmt.Println("config: Invalid or missing I_AM_A_LITTLE_BITCH - Continuing with default false")
		iAmALittleBitch = false
	}
	var (
		iAmALittleBitchCron string
		iAmALittleBitchUrl  string
	)

	if iAmALittleBitch {
		iAmALittleBitchCron = os.Getenv("I_AM_A_LITTLE_BITCH_CRON")
		if iAmALittleBitchCron == "" {
			fmt.Println("config: Invalid or missing I_AM_A_LITTLE_BITCH_CRON - Continuing with default '1 * * * *")
			iAmALittleBitchCron = "1 * * * *"
		}
		iAmALittleBitchUrl = os.Getenv("I_AM_A_LITTLE_BITCH_WEBHOOK_URL")
		if iAmALittleBitchUrl == "" {
			fmt.Println("config: Missing I_AM_A_LITTLE_BITCH_WEBHOOK_URL while I_AM_A_LITTLE_BITCH is enabled - Disabling I_AM_A_LITTLE_BITCH")
			iAmALittleBitch = false
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
		DbType: dbType,
		DbHost: dbHost,
		DbPort: dbPort,
		DbName: dbName,
		DbUser: dbUser,
		DbPass: dbPass,

		IAmALittleBitch:     iAmALittleBitch,
		IAmALittleBitchCron: iAmALittleBitchCron,
		IAmALittleBitchUrl:  iAmALittleBitchUrl,

		IgnoreList: ignoreList,
	}

	Env = &config

	fmt.Printf("config: Successfully setup config\n")
	return nil
}
