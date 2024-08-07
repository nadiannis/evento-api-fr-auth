package config

type Config struct {
	Port int
	DB   struct {
		DSN string
	}
	JWT struct {
		Secret string
	}
}
