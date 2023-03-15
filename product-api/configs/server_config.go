package configs

import (
	"time"
)

type ServerConfig interface {
	GetBindingAddr() string
	GetIdleTimeOut() time.Duration
	GetReadTimeOut() time.Duration
	GetWriteTimeOut() time.Duration
}

type serverConfig struct {
	Addr         string
	IdleTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

// NewServerConf returns initialized pointer of ServerConf
func NewServerConf(addr string, idleTO, readTO, wTO time.Duration) ServerConfig {
	return &serverConfig{
		Addr:         addr,
		IdleTimeOut:  idleTO,
		ReadTimeOut:  readTO,
		WriteTimeOut: wTO,
	}
}

func (s serverConfig) GetBindingAddr() string {
	return s.Addr
}

func (s serverConfig) GetIdleTimeOut() time.Duration {
	return s.IdleTimeOut
}

func (s serverConfig) GetReadTimeOut() time.Duration {
	return s.ReadTimeOut
}

func (s serverConfig) GetWriteTimeOut() time.Duration {
	return s.WriteTimeOut
}
