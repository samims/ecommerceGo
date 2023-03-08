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

type viperConfig struct {
	cfg *viper.Viper
}

func NewViperConfig() Env {
	v := viper.New()
	//envPAth := filepath.Join("..", ".env")
	//v.SetConfigFile(envPAth)
	c := &viperConfig{cfg: v}
	c.Load()
	return c
}

func (c *viperConfig) Load() {
	c.cfg.SetConfigFile(".env")

	c.cfg.AutomaticEnv()

	if err := c.cfg.ReadInConfig(); err != nil {
		panic(err)
	}
}

func (c *viperConfig) Get(key string) interface{} {
	val := c.cfg.Get(key)
	return val
}

func (c *viperConfig) GetString(key string) string {
	c.cfg.AutomaticEnv()
	val := c.cfg.GetString(key)
	return val
}

func (c *viperConfig) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

func (c *viperConfig) GetBool(key string) bool {
	return c.cfg.GetBool(key)
}

func (c *viperConfig) GetFloat64(key string) float64 {
	return c.cfg.GetFloat64(key)
}
