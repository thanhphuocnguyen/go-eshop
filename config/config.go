package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Domain               string        `mapstructure:"DOMAIN"`
	Port                 string        `mapstructure:"PORT"`
	DbUrl                string        `mapstructure:"DB_URL"`
	MaxPoolSize          int           `mapstructure:"MAX_POOL_SIZE"`
	MigrationPath        string        `mapstructure:"MIGRATION_PATH"`
	RedisUrl             string        `mapstructure:"REDIS_URL"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Env                  string        `mapstructure:"ENV"`
	CloudinaryUrl        string        `mapstructure:"CLOUDINARY_URL"`
	CloudinaryFolder     string        `mapstructure:"CLOUDINARY_FOLDER"`
	StripeSecretKey      string        `mapstructure:"STRIPE_SECRET_KEY"`
	StripePublishableKey string        `mapstructure:"STRIPE_PUBLISHABLE_KEY"`
	StripeWebhookSecret  string        `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	SmtpUsername         string        `mapstructure:"SMTP_USERNAME"`
	SmtpPassword         string        `mapstructure:"SMTP_PASSWORD"`
	SymmetricKey         string        `mapstructure:"SYMMETRIC_KEY"`
}

func LoadConfig(path string) (cfg Config, err error) {
	// Load config file
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&cfg)
	return
}
