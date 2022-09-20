package app

import (
	"github.com/urfave/cli"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Flags = NewAppFlags()

	return app
}

const (
	SymbolFlag            = "symbol"
	BindAddressFlag       = "bind-address"
	StoreDurationFlag     = "store-duration"
	LongClientIdFlag      = "long-client-id"
	LongClientSecretFlag  = "long-client-secret"
	ShortClientIdFlag     = "short-client-id"
	ShortClientSecretFlag = "short-client-secret"
)

func NewAppFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   BindAddressFlag,
			Usage:  "provide application host and port",
			EnvVar: "BIND_ADDRESS",
		},

		cli.StringFlag{
			Name:   SymbolFlag,
			Usage:  "provide symbol for trading",
			EnvVar: "SYMBOL",
		},

		cli.DurationFlag{
			Name:   StoreDurationFlag,
			Usage:  "time of the pricing data are stored",
			EnvVar: "STORE_DURATION",
		},

		cli.StringFlag{
			Name:   LongClientIdFlag,
			Usage:  "provide binance client id for long futures",
			EnvVar: "LONG_CLIENT_ID",
		},

		cli.StringFlag{
			Name:   LongClientSecretFlag,
			Usage:  "provide binance client secret for long futures",
			EnvVar: "LONG_CLIENT_SECRET",
		},

		cli.StringFlag{
			Name:   ShortClientIdFlag,
			Usage:  "provide binance client id for short futures",
			EnvVar: "SHORT_CLIENT_ID",
		},

		cli.StringFlag{
			Name:   ShortClientSecretFlag,
			Usage:  "provide binance client secret for short futures",
			EnvVar: "SHORT_CLIENT_SECRET",
		},
	}
}
