package configs

type Config struct {
	AllowedHosts []string
	ImageDIR     string
	MediaURL     string
	ServerCfg    *ServerConf
}

func NewConfig(allowedHosts []string, imageDIR, mediaURL string, sCfg *ServerConf) *Config {
	cfg := &Config{
		AllowedHosts: allowedHosts,
		ImageDIR:     imageDIR,
		ServerCfg:    sCfg,
		MediaURL:     mediaURL,
	}

	return cfg
}
