package configs

type Config struct {
	AllowedHosts []string
	FileDir      string
	MediaURL     string
	ServerCfg    *ServerConf
}

func NewConfig(allowedHosts []string, fileDIR, mediaURL string, sCfg *ServerConf) *Config {
	cfg := &Config{
		AllowedHosts: allowedHosts,
		FileDir:      fileDIR,
		ServerCfg:    sCfg,
		MediaURL:     mediaURL,
	}

	return cfg
}
