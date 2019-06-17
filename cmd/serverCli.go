package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"secrets/config"
	_ "secrets/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a server http.",
	Long:  "Start the webserver for the secrets.",
	Run:   execute,
}

func init() {
	setupLog()
	serverCmd.PersistentFlags().StringVar(&config.BoxPath, "box-path", config.DefaultBoxPath, "Box path")
	serverCmd.PersistentFlags().StringVarP(&config.Server, "host", "s", "127.0.0.1", "Host to server.")
	serverCmd.PersistentFlags().StringVarP(&config.Port, "port", "p", "8000", "Port to listen on host")
	rootCmd.AddCommand(serverCmd)
}

func execute(cmd *cobra.Command, args []string) {
	setupLog()
	log.WithFields(log.Fields{"host": config.Server, "port": config.Port}).Info("Running server")
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
