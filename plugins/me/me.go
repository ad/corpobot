package me

import (
	"fmt"

	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	"github.com/ad/corpobot/telegram"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Plugin struct{}

func init() {
	plugins.RegisterPlugin(&Plugin{})
}

func (m *Plugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("me.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("me", "Your ID/username", []string{"new", "member", "admin", "owner"}, me)
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[me.Plugin] Stopped")

	plugins.UnregisterCommand("me")
}

var me plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	msg := fmt.Sprintf("Hello %s, your ID: %d", user.UserName, user.TelegramID)
	return telegram.Send(user.TelegramID, msg)
}
