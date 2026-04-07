package conf

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zeroibot/pack/str"
)

// LoadEnv loads environment variables from given path
func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		parts := str.CleanSplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}
		if key == "" || value == "" {
			continue
		}
		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

// LoadRequiredEnv loads environment variables from given path and checks that all required keys are set
func LoadRequiredEnv(path string, requiredKeys []string) error {
	err := LoadEnv(path)
	if err != nil {
		return err
	}
	for _, key := range requiredKeys {
		value, ok := os.LookupEnv(key)
		if value == "" || !ok {
			return fmt.Errorf("missing env variable: %s", key)
		}
	}
	return nil
}
