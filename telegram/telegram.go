package telegram

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	dlog "github.com/amoghe/distillog"
	"golang.org/x/net/proxy"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	database "github.com/ad/corpobot/db"
	sql "github.com/lazada/sqle"
	"github.com/ad/corpobot/plugins"
)

// InitTelegram ...
func InitTelegram(token, proxyHost, proxyPort, proxyUser, proxyPassword string, debug bool) (bot *tgbotapi.BotAPI, err error) {
	var tr http.Transport

	if proxyHost != "" {
		tr = http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				socksDialer, err2 := proxy.SOCKS5(
					"tcp",
					fmt.Sprintf("%s:%s", proxyHost, proxyPort),
					&proxy.Auth{User: proxyUser, Password: proxyPassword},
					proxy.Direct,
				)
				if err2 != nil {
					return nil, err2
				}

				return socksDialer.Dial(network, addr)
			},
		}
	}

	bot, err = tgbotapi.NewBotAPIWithClient(token, &http.Client{Transport: &tr})
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	dlog.Debugf("Authorized on account @%s", bot.Self.UserName)

	return bot, nil
}

// ProcessTelegramMessages ...
func ProcessTelegramMessages(db *sql.DB, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	plugins.Bot = bot
	plugins.DB = db

	for update := range updates {
		updateGroupChat(db, update.Message)

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Text != "" {
			dlog.Debugf("%s [%d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

			message := database.TelegramMessage{
				TelegramID: update.Message.From.ID,
				FirstName: 	update.Message.From.FirstName,
				LastName: 	update.Message.From.LastName,
				UserName: 	update.Message.From.UserName,
				Message:  	update.Message.Text,
				Date:     	time.Unix(int64(update.Message.Date), 0),
			}


			err2 := database.StoreTelegramMessage(db, &message)
			if err2 != nil {
				dlog.Errorf("store message from user @%s [%d] failed: %s", message.UserName, message.TelegramID, err2)
			}
		}

		if update.Message.Command() != "" {
			if _, ok := plugins.Commands[update.Message.Command()]; ok {
				for _, d := range plugins.Plugins {
					upd := update
					result, err := d.Run(&upd)
					if err != nil {
						dlog.Errorln(err)
					}

					if result {
						break
					}
				}
			} else {
				err := Send(update.Message.Chat.ID, "unknown command, use /help")
				if err != nil {
					dlog.Errorln(err)
				}
			}
		}
	}
}

func updateGroupChat(db *sql.DB, message *tgbotapi.Message) {
	if (message == nil) {
		return
	}

	needCreate := false

	newChatMembers := message.NewChatMembers

	// TODO: сделать какую-то реакцию на удаление бота из чата *message.LeftChatMember
	if (message.Chat.Type == "supergroup" && newChatMembers != nil && len(*newChatMembers) > 0) {
		// ищем себя в списке, чтобы определить что нас добавили в какой-то групчат
		for i := range *newChatMembers {
			user := (*newChatMembers)[i]
			if user.ID == plugins.Bot.Self.ID {
				needCreate = true
			}
		}
	}

	if (message.Chat.Type == "supergroup") {
		needCreate = true
	}

	if needCreate {
		groupchat := database.Groupchat{
			Title: 		message.Chat.Title,
			TelegramID: message.Chat.ID,
			InviteLink: message.Chat.InviteLink,
			State: 		"active",
		}
		_, err := database.AddGroupChatIfNotExist(db, &groupchat)
		if err != nil && err.Error() != database.GroupChatAlreadyExists {
			dlog.Errorln(err)
		}
	}

	if message.NewChatTitle != "" {
		groupchat := database.Groupchat{
			Title: 		message.NewChatTitle,
			TelegramID: message.Chat.ID,
			InviteLink: message.Chat.InviteLink,
		}
		_, err := database.UpdateGroupChatTitle(db, &groupchat)
		if err != nil {
			dlog.Errorln(err)
		}
	}
}

// SendPlain ...
func Send(chatID int64, message string) (error) {
	return SendCustom(chatID, 0, message, false)
}

// SendMarkdown ...
func SendMarkdown(chatID int64, replyTo int, message string, isMarkdown bool) (error) {
	return SendCustom(chatID, replyTo, message, true)
}

// SendPlain ...
func SendPlain(chatID int64, replyTo int, message string) (error) {
	return SendCustom(chatID, replyTo, message, false)
}

// Send ...
func SendCustom(chatID int64, replyTo int, message string, isMarkdown bool) (error) {
	msg := tgbotapi.NewMessage(chatID, "")
	if isMarkdown {
		msg.ParseMode = "Markdown"
	}
	msg.Text = message
	msg.DisableWebPagePreview = true
	if replyTo != 0 {
		msg.ReplyToMessageID = replyTo
	}

	_, err := plugins.Bot.Send(msg)
	return err
}
