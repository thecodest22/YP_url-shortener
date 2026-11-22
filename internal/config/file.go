package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// readFile читает значения из файла в структуру конфига.
func readFile(path string, raw *rawConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file %q: %w", path, err)
	}
	defer file.Close()

	if err := cleanenv.ParseYAML(file, raw); err != nil {
		return fmt.Errorf("parse config file %q: %w", path, err)
	}

	return nil
}
