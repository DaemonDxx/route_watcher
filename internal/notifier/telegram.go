package notifier

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGNotifier struct {
	bot *tgbotapi.BotAPI
}

func NewTelegram(token string) (*TGNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot api client: %s", err.Error())
	}
	return &TGNotifier{
		bot: bot,
	}, nil
}

func (t *TGNotifier) Notify(chatID int64, message string) error {
	_, err := t.bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		return fmt.Errorf("error sending message: %s", err.Error())
	}
	return nil
}
