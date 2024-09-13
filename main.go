package main

import (
	"log"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"github.com/joho/godotenv"
	"github.com/thrillee/triq/apps"
	"github.com/thrillee/triq/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apps.MountApps()

	cmd.Execute()
}
