package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Load(path, profileOverride string) (Config, error) {
	v := viper.New()
	v.SetEnvPrefix("BIZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	applyDefaults(v)

	if err := readConfig(v, path); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	if strings.TrimSpace(profileOverride) != "" {
		cfg.Profile = strings.TrimSpace(profileOverride)
	}
	return cfg, nil
}

func readConfig(v *viper.Viper, path string) error {
	if path != "" {
		v.SetConfigFile(path)
		return v.ReadInConfig()
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "biz"))
	}
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return err
		}
	}
	return nil
}
