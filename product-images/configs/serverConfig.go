package configs

import (
	"time"
)

type ServerConf struct {
	Addr         string
	IdleTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	AllowedHosts []string
}

// NewServerConf returns initialized pointer of ServerConf
func NewServerConf(addr string, allowedHosts []string, idleTO, readTO, wTO time.Duration) *ServerConf {
	return &ServerConf{
		Addr:         addr,
		IdleTimeOut:  idleTO,
		ReadTimeOut:  readTO,
		WriteTimeOut: wTO,
		AllowedHosts: allowedHosts,
	}
}
