package config

import (
	"flag"
	"os"
)

// Config ...
type Config struct {
	TelegramToken         string
	TelegramProxyHost     string
	TelegramProxyPort     string
	TelegramProxyUser     string
	TelegramProxyPassword string
	TelegramDebug         bool
}

// InitConfig ...
func InitConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.TelegramToken, "telegram_token", lookupEnvOrString("CORPOBOT_TELEGRAM_TOKEN", config.TelegramToken), "telegramToken")
	flag.StringVar(&config.TelegramProxyHost, "telegram_proxy_host", lookupEnvOrString("CORPOBOT_TELEGRAM_PROXY_HOST", config.TelegramProxyHost), "telegramProxyHost")
	flag.StringVar(&config.TelegramProxyPort, "telegram_proxy_port", lookupEnvOrString("CORPOBOT_TELEGRAM_PROXY_PORT", config.TelegramProxyPort), "telegramProxyPort")
	flag.StringVar(&config.TelegramProxyUser, "telegram_proxy_user", lookupEnvOrString("CORPOBOT_TELEGRAM_PROXY_USER", config.TelegramProxyUser), "telegramProxyUser")
	flag.StringVar(&config.TelegramProxyPassword, "telegram_proxy_password", lookupEnvOrString("CORPOBOT_TELEGRAM_PROXY_PASSWORD", config.TelegramProxyPassword), "telegramProxyPassword")
	flag.BoolVar(&config.TelegramDebug, "telegram_debug", lookupEnvOrBool("CORPOBOT_TELEGRAM_DEBUG", config.TelegramDebug), "telegramDebug")

	flag.Parse()

	return config
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// func lookupEnvOrInt(key string, defaultVal int) int {
// 	if val, ok := os.LookupEnv(key); ok {
// 		if x, err := strconv.Atoi(val); err == nil {
// 			return x
// 		}
// 	}
// 	return defaultVal
// }

func lookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if val == "true" {
			return true
		}
		if val == "false" {
			return false
		}
	}
	return defaultVal
}
