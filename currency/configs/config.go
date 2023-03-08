package config


type AppConfig struct {
	// Add any configuration fields here
	Environment string
	Port        int
}

func LoadConfig(env Env) *AppConfig {
	config := &AppConfig{}

	// Load configuration values into AppConfig struct
	config.Environment = env.GetString("ENVIRONMENT")
	config.Port = env.GetInt("PORT")

	// Initialize logrus logger
//	logLevel, err := logger.Ne
//	if err != nil {
//		logLevel = logrus.InfoLevel
//	}
//	logrus.SetLevel(logLevel)
//	logrus.SetOutput(os.Stdout)
//	logrus.SetFormatter(&logrus.JSONFormatter{
//		TimestampFormat: "2006-01-02 15:04:05.000",
//	})

	return config
}
