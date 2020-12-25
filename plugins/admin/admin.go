package admin

import (
	"bytes"

	"github.com/ad/corpobot/plugins"

	database "github.com/ad/corpobot/db"
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
	if !plugins.CheckIfPluginDisabled("admin.AdminPlugin", "enabled") {
		return
	}

	plugins.RegisterCommand("pluginlist", "...")
	plugins.RegisterCommand("pluginenable", "...")
	plugins.RegisterCommand("plugindisable", "...")
}

func (m *AdminPlugin) OnStop() {
	dlog.Debugln("[AdminPlugin] Stopped")

	plugins.UnregisterCommand("pluginlist")
	plugins.UnregisterCommand("pluginenable")
	plugins.UnregisterCommand("plugindisable")
}

func (m *AdminPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "pluginlist" {
		return true, ListPlugins(update)
	}

	if update.Message.Command() == "pluginenable" {
		args := update.Message.CommandArguments()
		if plugins.EnablePlugin(args) {
			plugin := &database.Plugin{
				Name: args,
				State: "enabled",
			}

			_, err := database.UpdatePluginState(plugins.DB, plugin)
			if err != nil {
				return true, telegram.Send(update.Message.Chat.ID, "failed: " + err.Error())
			}

			return true, telegram.Send(update.Message.Chat.ID, args+" enabled")
		}

		return true, ListPlugins(update)
	}

	if update.Message.Command() == "plugindisable" {
		args := update.Message.CommandArguments()
		if plugins.DisablePlugin(args) {
			plugin := &database.Plugin{
				Name: args,
				State: "disabled",
			}

			_, err := database.UpdatePluginState(plugins.DB, plugin)
			if err != nil {
				return true, telegram.Send(update.Message.Chat.ID, "failed: " + err.Error())
			}

			return true, telegram.Send(update.Message.Chat.ID, args+" disabled")
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
