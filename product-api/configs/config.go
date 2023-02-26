package configs

type Config struct {
	AllowedHosts       []string
	ImageDIR           string
	MediaURL           string
	ServerCfg          *ServerConf
	CurrencyServerBase string
}

func NewConfig(allowedHosts []string, imageDIR, mediaURL, currencyServerBase string, sCfg *ServerConf) *Config {
	cfg := &Config{
		AllowedHosts:       allowedHosts,
		ImageDIR:           imageDIR,
		ServerCfg:          sCfg,
		MediaURL:           mediaURL,
		CurrencyServerBase: currencyServerBase,
	}

	return cfg
}
