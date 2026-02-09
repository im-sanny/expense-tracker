package config

import "fmt"

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}
