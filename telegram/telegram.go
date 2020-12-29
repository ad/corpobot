package telegram

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	"golang.org/x/net/proxy"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
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

	dlog.Debugf("Authorized on account @%s [id %d]", bot.Self.UserName, bot.Self.ID)

	return bot, nil
}

// ProcessTelegramMessages ...
func ProcessTelegramMessages(db *sql.DB, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	plugins.Bot = bot

	for update := range updates {
		updateGroupChat(db, update.Message)

		var user *database.User

		if update.CallbackQuery != nil {
			user = &database.User{
				TelegramID: int64(update.CallbackQuery.From.ID),
				FirstName:  update.CallbackQuery.From.FirstName,
				LastName:   update.CallbackQuery.From.LastName,
				UserName:   update.CallbackQuery.From.UserName,
				IsBot:      update.CallbackQuery.From.IsBot,
			}
			dlog.Debugf(" <= %s [%d] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Data)
		}
		if update.Message != nil {
			user = &database.User{
				TelegramID: int64(update.Message.From.ID),
				FirstName:  update.Message.From.FirstName,
				LastName:   update.Message.From.LastName,
				UserName:   update.Message.From.UserName,
				IsBot:      update.Message.From.IsBot,
			}
		}
		if update.CallbackQuery == nil && update.Message == nil {
			continue
		}

		user, err := database.AddUserIfNotExist(db, user)
		if err != nil && err.Error() != database.UserAlreadyExists {
			dlog.Errorln(err.Error())
			continue
		}

		if update.Message != nil && update.Message.Text != "" {
			dlog.Debugf(" <= %s [%d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

			message := database.TelegramMessage{
				TelegramID: update.Message.From.ID,
				FirstName:  update.Message.From.FirstName,
				LastName:   update.Message.From.LastName,
				UserName:   update.Message.From.UserName,
				Message:    update.Message.Text,
				IsIncoming: true,
				Date:       time.Unix(int64(update.Message.Date), 0),
			}

			err2 := database.StoreTelegramMessage(db, &message)
			if err2 != nil {
				dlog.Errorf("store message from user @%s [%d] failed: %s", message.UserName, message.TelegramID, err2)
			}
		}

		ProcessTelegramCommand(&update, user)
	}
}

// ProcessTelegramCommand ...
func ProcessTelegramCommand(update *tgbotapi.Update, user *database.User) {
	// check if message from private chat
	if update.Message != nil && update.Message.Chat.Type != "private" {
		return
	}

	command := ""
	if update.CallbackQuery != nil {
		command = strings.TrimLeft(update.CallbackQuery.Data, "/")
		commands := strings.Split(command, " ")
		if len(commands) > 1 {
			command = commands[0]
		}
	}

	if command == "" && update.Message.Command() != "" {
		command = update.Message.Command()
	}

	if command != "" {
		if _, ok := plugins.Commands[command]; ok {
			for _, d := range plugins.Plugins {
				result, err := d.Run(update, command, user)
				if err != nil {
					dlog.Errorln(err)
				}

				if result {
					break
				}
			}
		} else {
			err := Send(user.TelegramID, "unknown command "+command+", use /help")
			if err != nil {
				dlog.Errorln(err)
			}
		}
	}
}

func updateGroupChat(db *sql.DB, message *tgbotapi.Message) {
	if message == nil {
		return
	}

	needCreate := false

	newChatMembers := message.NewChatMembers

	// TODO: сделать какую-то реакцию на удаление бота из чата *message.LeftChatMember
	if message.Chat.Type == "supergroup" && newChatMembers != nil && len(*newChatMembers) > 0 {
		// ищем себя в списке, чтобы определить что нас добавили в какой-то групчат
		for i := range *newChatMembers {
			user := (*newChatMembers)[i]
			if user.ID == plugins.Bot.Self.ID {
				needCreate = true
			}
		}
	}

	if message.Chat.Type == "supergroup" {
		needCreate = true
	}

	if needCreate {
		groupchat := database.Groupchat{
			Title:      message.Chat.Title,
			TelegramID: message.Chat.ID,
			InviteLink: message.Chat.InviteLink,
			State:      "active",
		}
		_, err := database.AddGroupChatIfNotExist(db, &groupchat)
		if err != nil && err.Error() != database.GroupChatAlreadyExists {
			dlog.Errorln(err)
		}
	}

	if message.NewChatTitle != "" {
		groupchat := database.Groupchat{
			Title:      message.NewChatTitle,
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
func Send(chatID int64, message string) error {
	return SendCustom(chatID, 0, message, false, nil)
}

// SendMarkdown ...
func SendMarkdown(chatID int64, replyTo int, message string) error {
	return SendCustom(chatID, replyTo, message, true, nil)
}

// SendPlain ...
func SendPlain(chatID int64, replyTo int, message string) error {
	return SendCustom(chatID, replyTo, message, false, nil)
}

// Send ...
func SendCustom(chatID int64, replyTo int, message string, isMarkdown bool, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, "")
	if isMarkdown {
		msg.ParseMode = "Markdown"
	}
	msg.Text = message
	msg.DisableWebPagePreview = true
	if replyTo != 0 {
		msg.ReplyToMessageID = replyTo
	}

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	} else {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}

	_, err := plugins.Bot.Send(msg)
	if err != nil {
		return err
	}

	dlog.Debugf(" => %s [%d] %s", plugins.Bot.Self.UserName, plugins.Bot.Self.ID, message)

	storeMessage := database.TelegramMessage{
		TelegramID: int(chatID),
		Message:    message,
		Date:       time.Unix(time.Now().Unix(), 0),
		IsIncoming: false,
	}

	err2 := database.StoreTelegramMessage(plugins.DB, &storeMessage)
	if err2 != nil {
		dlog.Errorf("store message for user [%d] failed: %s", chatID, err2)
	}

	return nil
}

func GetArguments(update *tgbotapi.Update) string {
	if update.CallbackQuery != nil {
		command := strings.TrimLeft(update.CallbackQuery.Data, "/")
		commands := strings.Split(command, " ")
		if len(commands) > 1 {
			return strings.TrimSpace(strings.Join(commands[1:], ""))
		}
	}

	if update.Message != nil && update.Message.CommandArguments() != "" {
		return strings.TrimSpace(update.Message.CommandArguments())
	}

	return ""
}
