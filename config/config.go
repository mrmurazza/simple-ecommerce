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

	DBHost string `envconfig:"DB_HOST" default:"/Users/m.r.murazza/Project/personal/simple-ecommerce/ecommerce.db"`
}

func Get() *Config {
	c := Config{}
	once.Do(func() {
		envconfig.MustProcess("", &c)
		cfg = &c
	})

	return cfg
}
