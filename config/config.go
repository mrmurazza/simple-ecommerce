package config

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	ApplicationName string `envconfig:"APP_NAME" default:"simple-eccommerce"`

	JWTExpiryDuration time.Duration `envconfig:"JWT_EXPIRY" default:"1h"`
	JWTSignatureKey   string        `envconfig:"JWT_SIGN_KEY" default:"this is a secret"`

	DBHost     string `envconfig:"DB_HOST" default:"ecommerce.db"`
	DBDriver   string `envconfig:"DB_DRIVER" default:"sqlite3"`
	DBUser     string `envconfig:"DB_USER" default:""`
	DBPassword string `envconfig:"DB_PASSWORD" default:""`
	DBPort     string `envconfig:"DB_PORT" default:""`
	DBName     string `envconfig:"DB_NAME" default:"simple-ecommerce"`
}

func Get() *Config {
	c := Config{}
	once.Do(func() {
		envconfig.MustProcess("", &c)
		cfg = &c
	})

	return cfg
}
