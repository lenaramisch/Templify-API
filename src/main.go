package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.SMSService.com/pkg/db"
	"example.SMSService.com/pkg/domain"
	"example.SMSService.com/pkg/handler"
	"example.SMSService.com/pkg/router"
	emailservice "example.SMSService.com/pkg/service/email"
	mjmlservice "example.SMSService.com/pkg/service/mjml"
	smsservice "example.SMSService.com/pkg/service/sms"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SendgridConfig  emailservice.SendgridConfig
	SMSTwilioConfig smsservice.TwilioSMSSenderConfig
	MJMLConfig      mjmlservice.MJMLConfig
	DBConfig        db.RepositoryConfig
}

func loadConfig() AppConfig {
	godotenv.Load()
	dbPortString := os.Getenv("DB_PORT")
	dbPortInt, err := strconv.Atoi(dbPortString)
	if err != nil {
		log.Fatal("Converting db port to int failed")
	}

	return AppConfig{
		SendgridConfig: emailservice.SendgridConfig{
			ApiKey:       os.Getenv("API_KEY"),
			FromEmail:    os.Getenv("FROM_EMAIL"),
			FromName:     os.Getenv("FROM_NAME"),
			ReplyToEmail: os.Getenv("REPLY_TO_EMAIL"),
			ReplyToName:  os.Getenv("REPLY_TO_NAME"),
		},
		SMSTwilioConfig: smsservice.TwilioSMSSenderConfig{
			AccountSID: os.Getenv("ACCOUNT_SID"),
			AuthToken:  os.Getenv("AUTH_TOKEN"),
			FromNumber: os.Getenv("FROM_NUMBER"),
		},
		DBConfig: db.RepositoryConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPortInt,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		MJMLConfig: mjmlservice.MJMLConfig{},
	}
}

func main() {
	// prepare required services for usecase
	appConfig := loadConfig()
	sendgridEmailService := emailservice.NewSendGridService(appConfig.SendgridConfig)
	smsTwilioService := smsservice.NewTwilioSMSSender(appConfig.SMSTwilioConfig)
	mjmlService := mjmlservice.NewMJMLService(appConfig.MJMLConfig)

	// create usecase
	usecase := domain.NewUsecase(sendgridEmailService, smsTwilioService, mjmlService)

	//prepare handler
	handler := handler.NewAPIHandler(usecase)
	router := router.CreateRouter(handler)
	fmt.Println("Starting the server on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting the server")
	}
}
