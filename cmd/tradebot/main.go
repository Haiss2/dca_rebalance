package main

import (
	"fmt"
	"log"
	"os"

	libapp "github.com/Haiss2/dca/pkg/app"
	"github.com/Haiss2/dca/pkg/hunter"
	"github.com/Haiss2/dca/pkg/pricing"
	"github.com/Haiss2/dca/pkg/server"
	"github.com/Haiss2/dca/pkg/storage"
	"github.com/Haiss2/dca/pkg/telegram"
	"github.com/Haiss2/dca/pkg/trade"

	futu "github.com/adshao/go-binance/v2/futures"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var Version string = "0.1.0"

func main() {
	app := libapp.NewApp()
	app.Name = "HaiLe bot"
	app.Version = Version
	app.Action = run
	app.Flags = append(app.Flags, telegram.NewTelegramBotFlag()...)

	if err := app.Run(os.Args); err != nil {
		log.Panic(err)
	}
}

func run(c *cli.Context) error {
	// initiate logger
	logger, _, flush, err := libapp.NewLogger(c)
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	defer flush()
	zap.ReplaceGlobals(logger)
	l := logger.Sugar()
	l.Infow("app starting ..")

	// initiate telegram bot
	tele, err := telegram.NewTelegramBot(c)
	if err != nil {
		return fmt.Errorf("failed to initiate telegram bot: %w", err)
	}
	tele.Notify("DCA_MM Bot Hello World!")

	// initiate ram storage
	db := storage.NewRamStorage()

	// initiate client
	longClientId := c.String(libapp.LongClientIdFlag)
	longClientSecret := c.String(libapp.LongClientSecretFlag)
	longClient := futu.NewClient(longClientId, longClientSecret)

	shortClientId := c.String(libapp.ShortClientIdFlag)
	shortClientSecret := c.String(libapp.ShortClientSecretFlag)
	shortClient := futu.NewClient(shortClientId, shortClientSecret)

	// initiate pricing
	symbol := c.String(libapp.SymbolFlag)
	duration := c.Duration(libapp.StoreDurationFlag)
	synth := pricing.NewSynthetic(symbol, longClient, db, duration)

	// initiate trade module
	longTrade := trade.NewTradeModule(longClient)
	shortTrade := trade.NewTradeModule(shortClient)

	// initiate hunter and hunt
	hunter := hunter.NewHunter(c, synth, longTrade, shortTrade, tele)
	go hunter.Hunt()

	// initiate and run server
	bindAddress := c.String(libapp.BindAddressFlag)
	s := server.New(bindAddress)

	return s.Run()
}
