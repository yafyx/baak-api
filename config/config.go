package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	RootApp  string
	HTTPPort uint16
	BaseURL  string
}

var GlobalEnv Env

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
	}

	var ok bool

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	rootApp := strings.TrimSuffix(path, "/config")
	os.Setenv("APP_PATH", rootApp)
	GlobalEnv.RootApp = rootApp

	if port, err := strconv.Atoi(os.Getenv("PORT")); err != nil {
		panic("missing PORT environment")
	} else {
		GlobalEnv.HTTPPort = uint16(port)
	}

	GlobalEnv.BaseURL, ok = os.LookupEnv("BASE_URL")
	if !ok {
		panic("missing BASE_URL environment")
	}
}
