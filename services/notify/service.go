package notify

import (
	"context"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/forkpoons/cleanserver/core"
	"github.com/gaarx/gaarx"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"strings"
)

type Service struct {
	log  func() *logrus.Entry
	app  *gaarx.App
	ctx  context.Context
	bot  *tgbotapi.BotAPI
	stop chan bool
}

func Create(ctx context.Context) *Service {
	return &Service{
		ctx:  ctx,
		stop: make(chan bool),
	}
}

func (m *Service) Start(app *gaarx.App) error {
	m.app = app
	m.log = func() *logrus.Entry {
		return app.GetLog().WithField("service", "Notify service")
	}

	bot, err := tgbotapi.NewBotAPI("1139723632:AAHIz3-mA_KCdIpFxLsCIeAkBDqnpEEfW-Y")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	m.bot = bot
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if updates == nil {
		return err
	}
	for {
		select {
		case update := <-updates:
			if err := m.ProcessMessageFromTelegram(update); err != nil {
				return err
			}
		case data := <-m.app.Event(core.MessageEvent).Listen():
			if err := m.ProcessMessageFromWeb(data); err != nil {
				return err
			}
		case _ = <-m.stop:
			break
		}
	}
	return nil
}

func (m *Service) Stop() {
	m.stop <- true
}

func (m *Service) GetName() string {
	return "Notify service"
}

func (m *Service) ProcessMessageFromTelegram(update tgbotapi.Update) error {
	err := m.app.Storage().Set(core.TelegramUsersScope, strconv.FormatInt(update.Message.Chat.ID, 10), update.Message.From.UserName)
	if err != nil {
		return nil
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = m.bot.Send(msg)
	return nil
}

func (m *Service) ProcessMessageFromWeb(data interface{}) error {

	messages, _ := m.app.Storage().GetAll(core.TelegramUsersScope)
	str := fmt.Sprintf("%v", data)
	str = strings.TrimPrefix(str, "map[")
	str = strings.TrimSuffix(str, "]")
	for message := range messages {
		chatID, _ := strconv.ParseInt(message, 10, 64)
		msg := tgbotapi.NewMessage(chatID, str)
		_, _ = m.bot.Send(msg)
	}
	return nil
}
