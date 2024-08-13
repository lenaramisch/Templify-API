package app

import (
	"templify/pkg/db"
	domain "templify/pkg/domain/model"
	"templify/pkg/server"
	"templify/pkg/server/router"
	emailservice "templify/pkg/service/email"
	mjmlservice "templify/pkg/service/mjml"
	smsservice "templify/pkg/service/sms"
	typstservice "templify/pkg/service/typst"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Info   *domain.Info
	Router *router.Config
	Server *server.Config
	// Custom configs below
	SendgridConfig  *emailservice.SendgridConfig
	SMSTwilioConfig *smsservice.TwilioSMSSenderConfig
	MJMLConfig      *mjmlservice.MJMLConfig
	DBConfig        *db.RepositoryConfig
	TypstConfig     *typstservice.TypstConfig
}

func SetDefaults() {
	viper.SetDefault("APP_SERVER_PORT", "80")
	viper.SetDefault("APP_SERVER_TIMEOUT", 60*time.Second)
	viper.SetDefault("APP_SERVER_CORS_HEADERS", []string{"*"})
	viper.SetDefault("APP_SERVER_CORS_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("APP_SERVER_CORS_ORIGINS", []string{"*"})
}

func LoadConfig(
	version string,
	buildDate string,
	details string,
) (*Config, error) {
	SetDefaults()
	viper.AutomaticEnv()

	infoConfig := &domain.Info{
		Version:   version,
		BuildDate: buildDate,
		Details:   details,
	}

	routerConfig := &router.Config{
		Timeout: viper.GetDuration("APP_SERVER_TIMEOUT"),
		CORS: router.CORSConfig{
			AllowCredentials: viper.GetBool("APP_SERVER_CORS_ALLOW_CREDENTIALS"),
			Headers:          viper.GetStringSlice("APP_SERVER_CORS_HEADERS"),
			Methods:          viper.GetStringSlice("APP_SERVER_CORS_METHODS"),
			Origins:          viper.GetStringSlice("APP_SERVER_CORS_ORIGINS"),
		},
	}

	serverConfig := &server.Config{
		Address: "0.0.0.0:" + viper.GetString("APP_SERVER_PORT"),
		BaseURL: viper.GetString("APP_SERVER_BASE_URL"),
	}

	// Custom configs below
	mjmlConfig := &mjmlservice.MJMLConfig{
		Host: viper.GetString("MJML_HOST"),
		Port: viper.GetInt("MJML_PORT"),
	}

	dbConfig := &db.RepositoryConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
	}

	sendgridConfig := &emailservice.SendgridConfig{
		ApiKey:       viper.GetString("SENDGRID_API_KEY"),
		FromEmail:    viper.GetString("SENDGRID_FROM_EMAIL"),
		FromName:     viper.GetString("SENDGRID_FROM_NAME"),
		ReplyToEmail: viper.GetString("SENDGRID_REPLY_TO_EMAIL"),
		ReplyToName:  viper.GetString("SENDGRID_REPLY_TO_NAME"),
	}

	smsTwilioConfig := &smsservice.TwilioSMSSenderConfig{
		AccountSID: viper.GetString("TWILIO_ACCOUNT_SID"),
		AuthToken:  viper.GetString("TWILIO_AUTH_TOKEN"),
		FromNumber: viper.GetString("TWILIO_FROM_NUMBER"),
	}

	typstConfig := &typstservice.TypstConfig{}

	return &Config{
		Info:   infoConfig,
		Router: routerConfig,
		Server: serverConfig,
		// Custom configs below
		SendgridConfig:  sendgridConfig,
		SMSTwilioConfig: smsTwilioConfig,
		MJMLConfig:      mjmlConfig,
		DBConfig:        dbConfig,
		TypstConfig:     typstConfig,
	}, nil
}