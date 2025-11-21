package config

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" yaml:"host" validate:"required"`
	Port     int    `env:"DATABASE_PORT" yaml:"port" validate:"gt=0,lte=65535"`
	User     string `env:"DATABASE_USER" yaml:"user" validate:"required"`
	Password string `env:"DATABASE_PASSWORD" yaml:"password" validate:"required" sensitive:"true"`
	Name     string `env:"DATABASE_NAME" yaml:"db_name" validate:"required"`
	SSLMode  string `env:"DATABASE_SSL_MODE" yaml:"ssl_mode" validate:"oneof=disable allow prefer require verify-ca verify-full"`
}

type FiberConfig struct {
	Port int `env:"FIBER_PORT" yaml:"port" validate:"gt=0,lte=65535"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Fiber    FiberConfig    `yaml:"fiber"`
}
