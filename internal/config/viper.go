package config

import (
	"errors"
	"github.com/mgfeller/common-adapter-library/config"
	"github.com/spf13/viper"
)

const (
	ServerKey       = "server"
	MeshSpecKey     = "mesh"
	MeshInstanceKey = "instance"
)

type Viper struct {
	instance *viper.Viper
}

func NewViper(serverConfig map[string]string, meshConfig map[string]string, providerConfig map[string]string) (config.Handler, error) {
	v := viper.New()
	v.AddConfigPath(providerConfig["filepath"])
	v.SetConfigType(providerConfig["filetype"])
	v.SetConfigName(providerConfig["filename"])
	v.AutomaticEnv()

	for key, value := range serverConfig {
		v.SetDefault(ServerKey+"."+key, value)
	}

	for key, value := range meshConfig {
		v.SetDefault(MeshSpecKey+"."+key, value)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			// Config file was found but another error was produced
			return nil, config.ErrViper(err)
		}
	}
	return &Viper{
		instance: v,
	}, nil
}

func (v *Viper) SetKey(key string, value string) {
	v.instance.Set(key, value)
}

func (v *Viper) GetKey(key string) string {
	return v.instance.Get(key).(string)
}

func (v *Viper) Server(result interface{}) error {
	return v.instance.Sub(ServerKey).Unmarshal(result)
}

func (v *Viper) MeshSpec(result interface{}) error {
	return v.instance.Sub(MeshSpecKey).Unmarshal(result)
}

func (v *Viper) MeshInstance(result interface{}) error {
	return v.instance.Sub(MeshInstanceKey).Unmarshal(result)
}

func (v *Viper) Operations(result interface{}) error {
	return errors.New("config 'operations' not implemented")
}
