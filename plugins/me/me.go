package me

import (
	"fmt"

	"github.com/ad/corpobot/plugins"

	telegram "github.com/ad/corpobot/telegram"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type MePlugin struct {
}

func init() {
	plugins.RegisterPlugin(&MePlugin{})
}

func (m *MePlugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("me.MePlugin", "enabled") {
		return
	}

	plugins.RegisterCommand("me", "...")
}

func (m *MePlugin) OnStop() {
	dlog.Debugln("[MePlugin] Stopped")

	plugins.UnregisterCommand("me")
}

func (m *MePlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "me" {
		msg := fmt.Sprintf("Hello %s, your ID: %d", update.Message.From.UserName, update.Message.From.ID)

		return true, telegram.Send(update.Message.Chat.ID, msg)
	}

	return false, nil
}
