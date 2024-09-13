package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thrillee/triq/internals/servers"
)

var httpServerCMD = &cobra.Command{
	Use:   "http",
	Short: "Start Http Server",
	Long:  `Start Http Server`,
	Run:   startServer,
}

func init() {
	rootCmd.AddCommand(httpServerCMD)

	httpServerCMD.Flags().IntP("port", "p", 8000, "Server Port")
	httpServerCMD.Flags().StringP("host", "o", "0.0.0.0", "Host")
}

func startServer(cmd *cobra.Command, args []string) {
	host, err := cmd.Flags().GetString("host")
	if err != nil {
		log.Fatal(err)
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		log.Fatal(err)
	}

	serverProps := servers.HttpServerProps{
		Host: host,
		Port: port,
	}
	// apps.MountApps()
	servers.ListenAndServe(&serverProps)
}
