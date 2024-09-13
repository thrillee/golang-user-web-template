package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thrillee/triq/internals/schemas"
)

var dbMigrationCMD = &cobra.Command{
	Use:   "makemigrations",
	Short: "Migrate DB changes",
	Long:  `Migrate DB changes`,
	Run:   makeMigration,
}

func init() {
	rootCmd.AddCommand(dbMigrationCMD)
}

func makeMigration(cmd *cobra.Command, args []string) {
	schemas.LoadModels()
}
