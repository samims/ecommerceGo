package configs

import (
	"github.com/spf13/viper"
)

// AppConfig ...
type AppConfig interface {
	GetAllowedHosts() []string
	GetImageDir() string
	GetMediaURL() string
	GetCurrencyServerBase() string
}

type appConfig struct {
	AllowedHosts       []string `mapstructure:"bind_address"`
	ImageDIR           string   `mapstructure:"image_dir"`
	MediaURL           string   `mapstructure:"media_url"`
	CurrencyServerBase string   `mapstructure:"currency_server_base"`
}

func NewAppConfig(allowedHosts []string, imageDIR, mediaURL, currencyServerBase string) AppConfig {
	appCfg := &appConfig{
		AllowedHosts:       allowedHosts,
		ImageDIR:           imageDIR,
		MediaURL:           mediaURL,
		CurrencyServerBase: currencyServerBase,
	}
	return appCfg
}

func (a appConfig) GetAllowedHosts() []string {
	return viper.GetStringSlice("ALLOWED_HOSTS")
}

func (a appConfig) GetImageDir() string {
	return viper.GetString("IMAGE_DIR")
}

func (a appConfig) GetMediaURL() string {
	return viper.GetString("MEDIA_URL")
}

func (a appConfig) GetCurrencyServerBase() string {
	return viper.GetString("CURRENCY_SERVER_BASE")
}
