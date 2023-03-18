package config

import (
	"github.com/spf13/viper"
)

type Env interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
}

type ViperConfig struct {
	cfg *viper.Viper
}

func NewViperConfig() Env {
	v := viper.New()
	c := &ViperConfig{cfg: v}
	c.Load()
	return c
}

func (c *ViperConfig) Load() {
	c.cfg.SetConfigFile(".env")

	c.cfg.AutomaticEnv()

	if err := c.cfg.ReadInConfig(); err != nil {
		panic(err)
	}
}

func (c *ViperConfig) Get(key string) interface{} {
	val := c.cfg.Get(key)
	return val
}

func (c *ViperConfig) GetString(key string) string {
	c.cfg.AutomaticEnv()
	val := c.cfg.GetString(key)
	return val
}

func (c *ViperConfig) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

func (c *ViperConfig) GetBool(key string) bool {
	return c.cfg.GetBool(key)
}

func (c *ViperConfig) GetFloat64(key string) float64 {
	return c.cfg.GetFloat64(key)
}
