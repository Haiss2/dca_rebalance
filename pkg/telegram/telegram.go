package telegram

import (
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	botAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
)

const (
	telegramTokenFlag  = "telegram-token"
	telegramChatIDFlag = "telegram-chat-id"
)

type TelegramBot struct {
	api    *botAPI.BotAPI
	chatId int64
}

func NewTelegramBotFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   telegramTokenFlag,
			Usage:  "token for telegram bot",
			EnvVar: "TELEGRAM_BOT_TOKEN",
		},
		cli.Int64Flag{
			Name:   telegramChatIDFlag,
			Usage:  "telegram chat ID to notify",
			EnvVar: "TELEGRAM_CHAT_ID",
		},
	}
}

func NewTelegramBot(c *cli.Context) (*TelegramBot, error) {
	api, err := botAPI.NewBotAPI(c.GlobalString(telegramTokenFlag))

	if err != nil {
		return nil, err
	}
	chatId := c.Int64(telegramChatIDFlag)
	err = validation.Validate(chatId, validation.Required)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{
		api:    api,
		chatId: chatId,
	}, nil
}

func (t *TelegramBot) Notify(msg string) error {
	sendMsg := botAPI.NewMessage(t.chatId, msg)
	_, err := t.api.Send(sendMsg)
	if err != nil {
		log.Print(err)
	}
	return err
}
