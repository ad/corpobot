package me

import (
	"fmt"

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
	if !plugins.CheckIfPluginDisabled("me.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("me", "...", []string{"new", "member", "admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[me.Plugin] Stopped")

	plugins.UnregisterCommand("me")
}

func (m *Plugin) Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "me", user.Role) {
		msg := fmt.Sprintf("Hello %s, your ID: %d", user.UserName, user.TelegramID)

		return true, telegram.Send(user.TelegramID, msg)
	}

	return false, nil
}
