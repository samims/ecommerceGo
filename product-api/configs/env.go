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
	GetStringSlice(key string) []string
}

type env struct {
	viperLib *viper.Viper
}

func NewEnv() Env {
	v := viper.New()
	c := &env{viperLib: v}
	c.Load()
	return c
}

func (e *env) GetStringSlice(key string) []string {
	return e.viperLib.GetStringSlice(key)
}

func (e *env) Load() {
	e.viperLib.SetConfigFile(".env")
	e.viperLib.AutomaticEnv()

	if err := e.viperLib.ReadInConfig(); err != nil {
		panic(err)
	}
}

func (e *env) Get(key string) interface{} {
	return e.viperLib.Get(key)
}

func (e *env) GetString(key string) string {
	return e.viperLib.GetString(key)
}

func (e *env) GetInt(key string) int {
	return e.viperLib.GetInt(key)
}

func (e *env) GetBool(key string) bool {
	return e.viperLib.GetBool(key)
}

func (e *env) GetFloat64(key string) float64 {
	return e.viperLib.GetFloat64(key)
}

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
