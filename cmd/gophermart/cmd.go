package main

import (
	"github.com/caarlos0/env/v10"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run gophermart",
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			address, _ := flags.GetString("a")
			databaseDsn, _ := flags.GetString("d")

			config := gophermart.Config{
				Address: address,
				DatabaseDsn: databaseDsn,
			}
			if err := env.Parse(&config); err != nil {
				log.Fatal(err)
			}
	
			app := gophermart.NewApp(config)
			if err := app.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringP("a", "a", ":8080", "Server address")
	flags.StringP("d", "d", "", "Database connection string")
	return cmd
}
