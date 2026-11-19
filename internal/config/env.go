package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// readEnv читает значения из переменных окружения и пишет в структуру.
func readEnv(raw *rawConfig) error {
	if err := cleanenv.ReadEnv(raw); err != nil {
		return fmt.Errorf("read environment config: %w", err)
	}

	return nil
}
