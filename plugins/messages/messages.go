package messages

// https://core.telegram.org/bots/api#deletemessage

import (
	"strconv"
	"strings"

	"github.com/ad/corpobot/plugins"

	database "github.com/ad/corpobot/db"
	telegram "github.com/ad/corpobot/telegram"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Plugin struct {
}

func init() {
	plugins.RegisterPlugin(&Plugin{})
}

func (m *Plugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("messages.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("broadcast", "Send message to all users", []string{"admin", "owner"})
	plugins.RegisterCommand("message", "Send message to user", []string{"admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[messages.Plugin] Stopped")

	plugins.UnregisterCommand("broadcast")
	plugins.UnregisterCommand("message")
}

func (m *Plugin) Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "broadcast", user.Role) {
		if args == "" {
			return true, telegram.Send(user.TelegramID, "failed: empty message")
		}

		users, err := database.GetUsers(plugins.DB, []string{})
		if err != nil {
			return true, telegram.Send(user.TelegramID, err.Error())
		}

		if len(users) > 0 {
			var usersList []string

			for _, u := range users {
				err = telegram.Send(u.TelegramID, args)
				if err != nil {
					usersList = append(usersList, "* "+u.String()+" — failed: "+err.Error())
				} else {
					usersList = append(usersList, "* "+u.String()+" — success")
				}
			}

			return true, telegram.Send(user.TelegramID, strings.Join(usersList, "\n"))
		}

		return true, telegram.Send(user.TelegramID, args+" broadcast")
	}

	if plugins.CheckIfCommandIsAllowed(command, "message", user.Role) {
		errorString := "failed: you must provide user id message with a new line between them"
		params := strings.Split(args, "\n")

		if len(params) != 2 {
			return true, telegram.Send(user.TelegramID, "failed: empty message")
		}

		userIDstring, message := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

		if userIDstring == "" || message == "" {
			return true, telegram.Send(user.TelegramID, errorString)
		}

		userID, err := strconv.ParseInt(userIDstring, 10, 64)
		if err != nil {
			return true, telegram.Send(user.TelegramID, errorString)
		}

		err = telegram.Send(userID, message)
		if err != nil {
			return true, telegram.Send(user.TelegramID, err.Error())
		}

		return true, telegram.Send(user.TelegramID, "message sent")
	}

	return false, nil
}
