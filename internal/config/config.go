package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Load() (*AppConfiguration, error) {
	return loadConfig(viper.New(), pflag.CommandLine)
}

func loadCommandLineArguments(v *viper.Viper, flags *pflag.FlagSet) error {
	flags.StringP("config", "c", "", "Specifies path for config file")

	if err := v.BindPFlags(flags); err != nil {
		return err
	}

	pflag.Parse()
	return nil
}

func loadConfig(v *viper.Viper, flags *pflag.FlagSet) (*AppConfiguration, error) {
	config := &AppConfiguration{}

	if err := loadCommandLineArguments(v, flags); err != nil {
		return nil, err
	}

	if err := loadConfigFile(v); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfigFile(v *viper.Viper) error {
	configFile := v.GetString("config")
	if configFile == "" {
		v.AddConfigPath("./app/config")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	} else {
		v.SetConfigFile(configFile)
	}

	return v.ReadInConfig()
}
