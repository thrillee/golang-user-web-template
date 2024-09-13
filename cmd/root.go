/*
Copyright © 2023 Bellotobiloba01@gmail.com
*/
package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   `faruk - Web Framework The Django of Go`,
	Short: "This is a go web framework",
	Long:  `This is a go web framework`,
}

func Execute() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
