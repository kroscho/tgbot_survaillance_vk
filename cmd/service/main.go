package main

import (
	"context"
	"os"
	"tgbot_surveillance/config"
	"tgbot_surveillance/internal/domain/tracked"
	postgresTracked "tgbot_surveillance/internal/domain/tracked/postgres"
	"tgbot_surveillance/internal/domain/user"
	postgresUser "tgbot_surveillance/internal/domain/user/postgres"
	"tgbot_surveillance/internal/domain/userVk"
	postgresUserVk "tgbot_surveillance/internal/domain/userVk/postgres"
	"tgbot_surveillance/pkg/clock"
	"tgbot_surveillance/pkg/database/psql"
	"tgbot_surveillance/pkg/logger"
	"tgbot_surveillance/transport/telegram"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	logger := logger.New()

	if err := run(logger); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func run(logger logrus.FieldLogger) error {
	defer logger.Info("graceful shutdown successfully finished")

	clk := clock.Real{}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found")
	}

	cfg := config.New()

	logger.Infof("starting service")

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Authorized on account %s", bot.Self.UserName)

	dbClient, err := psql.New(cfg.StorageDSN, logger)
	if err != nil {
		return errors.Wrap(err, "create db client")
	}
	defer dbClient.Close()

	userService := user.NewService(
		postgresUser.NewStore(dbClient.GetConnection(), clk),
	)

	trackedSevice := tracked.NewService(
		userService,
		postgresTracked.NewStore(dbClient.GetConnection(), clk),
	)

	userVkService := userVk.NewService(
		userService,
		postgresUserVk.NewStore(dbClient.GetConnection(), clk),
	)

	services := telegram.Services{
		UserService:    userService,
		TrackedService: trackedSevice,
		UserVkService:  userVkService,
	}

	server := telegram.NewServer(bot, logger, services)
	err = server.Run(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	return nil
}
