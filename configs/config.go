package configs

type Config struct {
	AllowedHosts []string
	FileDir      string
	ServerCfg    *ServerConf
}

func NewConfig(allowedHosts []string, fileDIR string, sCfg *ServerConf) *Config {
	cfg := &Config{
		AllowedHosts: allowedHosts,
		FileDir:      fileDIR,
		ServerCfg:    sCfg,
	}

	return cfg
}
