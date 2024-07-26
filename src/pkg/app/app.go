package app

import (
	"context"
	"os"
	"os/signal"
	"templify/pkg/db"
	"templify/pkg/domain/usecase"
	"templify/pkg/logging"
	"templify/pkg/server"
	"templify/pkg/server/handler/apihandler"
	"templify/pkg/server/router"
	emailservice "templify/pkg/service/email"
	mjmlservice "templify/pkg/service/mjml"
	smsservice "templify/pkg/service/sms"
	typstservice "templify/pkg/service/typst"
)

// Run runs the app
// nolint: funlen
func Run(cfg *Config, shutdownChannel chan os.Signal) error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// ===== Logger =====
	logger := logging.SetLogger()

	sendgridEmailService := emailservice.NewSendGridService(cfg.SendgridConfig)
	smsTwilioService := smsservice.NewTwilioSMSSender(cfg.SMSTwilioConfig)
	mjmlService := mjmlservice.NewMJMLService(cfg.MJMLConfig)
	repository := db.NewRepository(cfg.DBConfig)
	typstService := typstservice.NewTypstService(cfg.TypstConfig)

	// ===== App Logic =====
	appLogic := usecase.NewUsecase(sendgridEmailService, smsTwilioService, mjmlService, repository, typstService)

	// ===== Handlers =====
	apiHandler := apihandler.NewAPIHandler(appLogic, cfg.Info, logger, cfg.Server.BaseURL)

	// ===== Router =====
	r := router.New(apiHandler, cfg.Router)

	// ===== Server =====
	srv := server.NewServer(cfg.Server, r)

	srvErr := make(chan error, 1)
	go func() {
		logger.Info("server started!", "address", cfg.Server.Address)
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	<-ctx.Done()

	// Stop receiving signal notifications as soon as possible.
	err := srv.Shutdown(context.Background())
	if err != nil {
		logger.Error("server shutdown error", "error", err)
		return err
	}
	stop()

	return nil
}
