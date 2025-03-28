package configs

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// Postgres sozlamalari
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`

	// Umumiy xizmat sozlamalari
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	LoggerLevel string `mapstructure:"LOGGER_LEVEL"`

	// HTTP server sozlamalari
	HTTPHost string `mapstructure:"HTTP_HOST"`
	HTTPPort int    `mapstructure:"HTTP_PORT"`

	// Email sozlamalari
	From       string `mapstructure:"EMAIL_SENDER"`
	SenderPass string `mapstructure:"EMAIL_SENDER_PASS"`
	SMTPHost   string `mapstructure:"EMAIL_SMTP_HOST"`
	SMTPPort   string `mapstructure:"EMAIL_SMTP_PORT"`

	// Eski Config dan qo‘shilgan maydonlar
	ServerPort       string `mapstructure:"SERVER_PORT"`
	WebappBaseURL    string `mapstructure:"WEBAPP_BASE_URL"`
	EnvType          string `mapstructure:"ENV_TYPE"`
	DBSource         string `mapstructure:"DB_SOURCE"`
	ChatAPIHost      string `mapstructure:"CHAT_API_HOST"`
	ChatBaseURL      string `mapstructure:"CHAT_BASE_URL"`
	SecretSessionKey string `mapstructure:"SECRET_SESSION_KEY"`
}

var prodRequiredVariables = []string{
	"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
	"SERVICE_NAME", "ENVIRONMENT", "LOGGER_LEVEL",
	"REDIS_HOST", "REDIS_PORT",
	"HTTP_PORT", "HTTP_HOST",
	"EMAIL_SENDER", "EMAIL_SENDER_PASS", "EMAIL_SMTP_HOST", "EMAIL_SMTP_PORT",
	"SERVER_PORT", "WEBAPP_BASE_URL", "ENV_TYPE", "DB_SOURCE", "CHAT_BASE_URL", "SECRET_SESSION_KEY",
}

var devRequiredVariables = []string{
	"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
	"SERVICE_NAME", "ENVIRONMENT", "LOGGER_LEVEL",
	"REDIS_HOST", "REDIS_PORT",
	"HTTP_PORT", "HTTP_HOST",
	"EMAIL_SENDER", "EMAIL_SENDER_PASS", "EMAIL_SMTP_HOST", "EMAIL_SMTP_PORT",
	"SERVER_PORT", "WEBAPP_BASE_URL", "ENV_TYPE", "DB_SOURCE", "CHAT_BASE_URL", "SECRET_SESSION_KEY",
}

func LoadConfig(path string) (Config, error) {
	err := godotenv.Load(path + ".env")
	if err != nil {
		log.Printf(".env fayl topilmadi, standart qiymatlardan foydalaniladi: %v", err)
	}

	viper.SetConfigFile(path + ".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	viper.SetDefault("POSTGRES_HOST", "postgres")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "1111")
	viper.SetDefault("POSTGRES_DB", "med")

	viper.SetDefault("SERVICE_NAME", "med-service")
	viper.SetDefault("ENVIRONMENT", "dev")
	viper.SetDefault("LOGGER_LEVEL", "debug")

	viper.SetDefault("HTTP_HOST", "localhost")
	viper.SetDefault("HTTP_PORT", 8085)

	viper.SetDefault("EMAIL_SENDER", "hatamovjamshid47@gmail.com")
	viper.SetDefault("EMAIL_SENDER_PASS", "yfvo hhhq xurb mind")
	viper.SetDefault("EMAIL_SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("EMAIL_SMTP_PORT", "587")

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("WEBAPP_BASE_URL", "http://localhost:8080")
	viper.SetDefault("ENV_TYPE", "dev")
	viper.SetDefault("DB_SOURCE", "postgres://postgres:1111@localhost:5432/med")
	viper.SetDefault("CHAT_API_HOST", "localhost")
	viper.SetDefault("CHAT_BASE_URL", "http://localhost:3000")
	viper.SetDefault("SECRET_SESSION_KEY", "your_secret_key")

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Konfiguratsiyani unmarshal qilishda xatolik: %v", err)
	}

	const devEnvType = "dev"
	const prodEnvType = "prod"
	if cfg.EnvType != devEnvType && cfg.EnvType != prodEnvType {
		fmt.Println("!!!!!!!!!!!!!!!!!!!")
		log.Fatalf(`ENV_TYPE o‘zgaruvchisi "dev" yoki "prod" bo‘lishi kerak`)
	}

	if cfg.EnvType == devEnvType {
		for _, key := range devRequiredVariables {
			if !viper.IsSet(key) || viper.GetString(key) == "" {
				log.Printf("Ogohlantirish: majburiy o‘zgaruvchi %s o‘rnatilmagan, standart qiymat ishlatiladi", key)
			}
		}
	}

	if cfg.EnvType == prodEnvType {
		for _, key := range prodRequiredVariables {
			if !viper.IsSet(key) || viper.GetString(key) == "" {
				fmt.Println("!!!!!!!!!!!!!!!!!!!")
				log.Fatalf("Majburiy o‘zgaruvchi %s o‘rnatilmagan", key)
			}
		}
	}

	return cfg, nil
}
