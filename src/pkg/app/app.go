package app

import (
	"context"
	"os"
	"os/signal"
	"templify/pkg/db"
	"templify/pkg/domain/usecase"
	"templify/pkg/logging"
	"templify/pkg/router"
	"templify/pkg/server"
	generatedAPI "templify/pkg/server/generated"
	"templify/pkg/server/handler/apihandler"
	emailservice "templify/pkg/service/email"
	filemanager "templify/pkg/service/filemanager"
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

	sendgridEmailService := emailservice.NewSendGridService(cfg.SendgridConfig, logger)
	smsTwilioService := smsservice.NewTwilioSMSSender(cfg.SMSTwilioConfig, logger)
	mjmlService := mjmlservice.NewMJMLService(cfg.MJMLConfig, logger)
	repository := db.NewRepository(cfg.DBConfig, logger)
	typstService := typstservice.NewTypstService(cfg.TypstConfig, logger)

	var filemanagerService usecase.FileManagerService
	if cfg.FileManagerConfig.IsAWS {
		filemanagerService = filemanager.NewFileManagerAWSService(cfg.FileManagerConfig, logger)
	} else {
		filemanagerService = filemanager.NewFileManagerMinioService(cfg.FileManagerConfig, logger)
	}
	// ===== App Logic =====
	appLogic := usecase.NewUsecase(sendgridEmailService, smsTwilioService, mjmlService, repository, typstService, filemanagerService, logger)

	// ===== Handlers =====
	apiHandler := apihandler.NewAPIHandler(appLogic, cfg.Info, logger, cfg.Server.BaseURL)

	// ===== Router =====
	handler := generatedAPI.HandlerFromMux(apiHandler, nil)
	swagger, err := generatedAPI.GetSwagger()
	if err != nil {
		logger.Error("failed to get swagger", "error", err)
		return err
	}
	r := router.New(handler, cfg.Router, logger, swagger)

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
	err = srv.Shutdown(context.Background())
	if err != nil {
		logger.Error("server shutdown error", "error", err)
		return err
	}
	stop()

	return nil
}
