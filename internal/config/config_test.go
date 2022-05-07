package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFileSucceeds(t *testing.T) {
	v := viper.New()
	v.Set("config", "../../config/local.yaml")

	config, err := loadConfig(v, pflag.NewFlagSet("flagset", pflag.ContinueOnError))

	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, "wss://ws-feed.exchange.coinbase.com", config.Spec.Websocket)
}

func TestLoadCommandLineArguments(t *testing.T) {
	v := viper.New()
	flagset := pflag.NewFlagSet("flagset", pflag.ContinueOnError)

	err := loadCommandLineArguments(v, flagset)

	require.NoError(t, err)

	_ = flagset.Set("config", "./path/config")

	require.Equal(t, "./path/config", v.GetString("config"))
}
