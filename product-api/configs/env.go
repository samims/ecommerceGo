package configs

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

type env struct {
	cfg *viper.Viper
}

func NewEnv() Env {
	v := viper.New()
	c := &env{cfg: v}
	c.Load()
	return c
}

func (e *env) Load() {
	e.cfg.SetConfigFile("product-api/.env")
	e.cfg.AutomaticEnv()

	if err := e.cfg.ReadInConfig(); err != nil {
		panic(err)
	}
}

func (e *env) Get(key string) interface{} {
	return e.cfg.Get(key)
}

func (e *env) GetString(key string) string {
	return e.cfg.GetString(key)
}

func (e *env) GetInt(key string) int {
	return e.cfg.GetInt(key)
}

func (e *env) GetBool(key string) bool {
	return e.cfg.GetBool(key)
}

func (e *env) GetFloat64(key string) float64 {
	return e.cfg.GetFloat64(key)
}

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
