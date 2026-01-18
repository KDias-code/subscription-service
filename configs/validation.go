package configs

import "fmt"

func (c *Config) Validate() error {
	if c.Postgres.Dsn == "" {
		return fmt.Errorf("postgres.dsn is required")
	}

	return nil
}
