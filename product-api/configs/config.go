package configs

// Config defines the main configuration interface.
type Config interface {
	ServerConfig() ServerConfig
	AppConfig() AppConfig
}

type config struct {
	serverConfig *serverConfig
	appConfig    *appConfig
}

func (c *config) ServerConfig() ServerConfig {
	return c.serverConfig
}

func (c *config) AppConfig() AppConfig {
	return c.appConfig
}

func NewConfig(sConf ServerConfig, aConf AppConfig) Config {
	return &config{
		serverConfig: sConf.(*serverConfig),
		appConfig:    aConf.(*appConfig),
	}
}
