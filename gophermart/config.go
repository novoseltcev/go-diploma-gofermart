package gophermart


type Config struct {
	Address			string	`env:"RUN_ADDRESS"`
	DatabaseDsn		string	`env:"DATABASE_URI"`
	AccrualAddress	string	`env:"ACCRUAL_SYSTEM_ADDRESS"`
	JwtSecret		string	`env:"JWT_SECRET"`
	JwtLifetime		int8	`env:"JWT_LIFETIME"`
}
