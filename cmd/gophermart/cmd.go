package main

import (
	"github.com/caarlos0/env/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Use this command to run gophermart",
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			address, _ := flags.GetString("a")
			databaseDsn, _ := flags.GetString("d")
			accrualAddress, _ := flags.GetString("r")
			jwtSecret, _ := flags.GetString("s")
			jwtLifetime, _ := flags.GetInt("l")

			config := gophermart.Config{
				Address: address,
				DatabaseDsn: databaseDsn,
				AccrualAddress: accrualAddress,
				JwtSecret: jwtSecret,
				JwtLifetime: int8(jwtLifetime),
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
	flags.StringP("d", "d", "", "Database URI")
	flags.StringP("r", "r", "", "Accrual system address")
	flags.StringP("secret", "s", "", "JWT secret key")
	flags.IntP("l", "l", 0, "JWT lifetime (in days)")
	return cmd
}
