package database

// Config содержит параметры подключения к PostgreSQL
type Config struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:""`
	DBName   string `env:"DB_NAME" envDefault:"myapp"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}
