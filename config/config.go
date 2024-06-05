package config

import (
	"fmt"
	"github.com/libcodebr/didyoumean/infra/mongodb"
	v "github.com/libcodebr/didyoumean/pkg/verifier"
	"github.com/spf13/viper"
	"os"
)

type ManagerConfig struct {
	//Logger     *l.Config          `mapstructure:"logger"`
	Mongo *mongodb.Config `mapstructure:"mongo" validate:"required"`
}

// LoadConfig loads the configuration from the file
func LoadConfig(vp *viper.Viper, configFile string) (*ManagerConfig, error) {
	var cfg = new(ManagerConfig)

	vp.SetConfigFile(configFile)

	vp.SetEnvPrefix("didyoumean")
	vp.AutomaticEnv()

	if err := vp.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "using config file:", vp.ConfigFileUsed())
	}

	if err := vp.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := v.Verifier.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
