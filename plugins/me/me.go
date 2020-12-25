package me

import (
	"fmt"

	"github.com/ad/corpobot/plugins"

	database "github.com/ad/corpobot/db"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	telegram "github.com/ad/corpobot/telegram"
)

type MePlugin struct {

}

func init() {
	plugins.RegisterPlugin(&MePlugin{})
}
func (m *MePlugin) OnStart() {
	plugin := &database.Plugin{
		Name: "me.MePlugin",
		State: "enabled",
	}

	plugin, err := database.AddPluginIfNotExist(plugins.DB, plugin)
	if err != nil {
		dlog.Errorln("failed: " + err.Error())
	}

	if plugin.State != "enabled" {
		dlog.Debugln("[MePlugin] Disabled")
		return
	}

	
	dlog.Debugln("[MePlugin] Started")

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