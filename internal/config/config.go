package config

type Config struct {
	HTTPPort string
	PostgresURL string
	RedisAddr string
}

func Load() *Config {
	return &Config{
		HTTPPort: ":8080",
		PostgresURL: "postgres://admin:admin@localhost:5432/distributed-task-queue",
		RedisAddr: "localhost:6380",
	}
}