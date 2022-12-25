package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	commandStart           = "start"
	replyStart             = "Привет! Чтобы сохранять ссылки в своём Pocket аккаунте, тебе нужен доступ. Для этого перейди по ссылке:\n%s"
	replyAlreadyAuthorized = "Ты уже авторизирован. Отправляй мне ссылку, чтобы я её сохранил."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	default:
		return b.handleUnknownCommand(message)

	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена")

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Это невалидная ссылка"
		_, err = b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = "Ты не авторизирован! Используй команду /start"
		_, err = b.bot.Send(msg)
		return err
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		msg.Text = "Не удалось сохранить ссылку! Попробуй позже."
		_, err = b.bot.Send(msg)
		return err
	}

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "unknown command")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}
