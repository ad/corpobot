package admin

import (
	"bytes"

	"github.com/ad/corpobot/plugins"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	telegram "github.com/ad/corpobot/telegram"
)

type AdminPlugin struct {

}

func init() {
	plugins.RegisterPlugin(&AdminPlugin{})
}

func (m *AdminPlugin) OnStart() {
	dlog.Debugln("[AdminPlugin] Started")

	plugins.RegisterCommand("listplugins", "...")
	plugins.RegisterCommand("enableplugin", "...")
	plugins.RegisterCommand("disableplugin", "...")
}

func (m *AdminPlugin) OnStop() {
	dlog.Debugln("[AdminPlugin] Stopped")

	plugins.UnregisterCommand("listplugins")
	plugins.UnregisterCommand("enableplugin")
	plugins.UnregisterCommand("disableplugin")
}

func (m *AdminPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "listplugins" {
		return true, ListPlugins(update)
	}

	if update.Message.Command() == "enableplugin" {
		args := update.Message.CommandArguments()
		if plugins.EnablePlugin(args) {
			err := telegram.Send(update.Message.Chat.ID, args+" enabled")

			if err != nil {
				return true, err 
			}
		}

		return true, ListPlugins(update)
	}

	if update.Message.Command() == "disableplugin" {
		args := update.Message.CommandArguments()
		if plugins.DisablePlugin(args) {
			err := telegram.Send(update.Message.Chat.ID, args+" disabled")

			if err != nil {
				return true, err 
			}
		}

		return true, ListPlugins(update)
	}

	return false, nil
}

func ListPlugins(update *tgbotapi.Update) (error) {
	var loaded bytes.Buffer
	var unloaded bytes.Buffer

	for k := range plugins.Plugins {
		_, err := loaded.WriteString("\t" + k + "\n")
		if err != nil {
			return err
		}
	}

	for k := range plugins.DisabledPlugins {
		_, err := unloaded.WriteString("\t" + k + "\n")
		if err != nil {
			return err
		}
	}

	return telegram.Send(update.Message.Chat.ID, "Enabled plugins:\n" + loaded.String() + "\nDisabled plugins:\n" + unloaded.String())
}
