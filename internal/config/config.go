package config

import (
	"log"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env              string `env:"ENV,required"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	LogFilePath      string `env:"LOG_FILE_PATH" envDefault:"./kladovkin-telegram-bot.log"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.Env),
		slog.String("telegram_bot_token", "********"),
		slog.String("log_file_path", c.LogFilePath),
	)
}

func MustLoad() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panicf("failed to load configuration: %v", err)
	}
	return &cfg
}
