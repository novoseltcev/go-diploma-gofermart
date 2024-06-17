package gophermart

type Config struct {
	Address 	string		`env:"ADDRESS"`
	DatabaseDsn string		`env:"DATABASE_DSN"`
}
