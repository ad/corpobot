package me

import (
	"fmt"

	"github.com/ad/corpobot/plugins"

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

	plugins.RegisterCommand("me", "...")
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[me.Plugin] Stopped")

	plugins.UnregisterCommand("me")
}

func (m *Plugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "me" {
		msg := fmt.Sprintf("Hello %s, your ID: %d", update.Message.From.UserName, update.Message.From.ID)

		return true, telegram.Send(update.Message.Chat.ID, msg)
	}

	return false, nil
}
