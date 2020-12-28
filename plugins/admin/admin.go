package admin

import (
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
	if !plugins.CheckIfPluginDisabled("admin.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("pluginlist", "...", []string{"owner"})
	plugins.RegisterCommand("pluginenable", "...", []string{"owner"})
	plugins.RegisterCommand("plugindisable", "...", []string{"owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[admin.Plugin] Stopped")

	plugins.UnregisterCommand("pluginlist")
	plugins.UnregisterCommand("pluginenable")
	plugins.UnregisterCommand("plugindisable")
}

func (m *Plugin) Run(update *tgbotapi.Update, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "pluginlist", user.Role) {
		return true, ListPlugins(update)
	}

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "pluginenable", user.Role) {
		args := update.Message.CommandArguments()
		if plugins.EnablePlugin(args) {
			plugin := &database.Plugin{
				Name:  args,
				State: "enabled",
			}

			_, err := database.UpdatePluginState(plugins.DB, plugin)
			if err != nil {
				return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
			}

			return true, telegram.Send(update.Message.Chat.ID, args+" enabled")
		}

		return true, ListPlugins(update)
	}

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "plugindisable", user.Role) {
		args := update.Message.CommandArguments()
		if plugins.DisablePlugin(args) {
			plugin := &database.Plugin{
				Name:  args,
				State: "disabled",
			}

			_, err := database.UpdatePluginState(plugins.DB, plugin)
			if err != nil {
				return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
			}

			return true, telegram.Send(update.Message.Chat.ID, args+" disabled")
		}

		return true, ListPlugins(update)
	}

	return false, nil
}

func ListPlugins(update *tgbotapi.Update) error {
	buttons := make([][]tgbotapi.KeyboardButton, 0)

	for k := range plugins.Plugins {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/plugindisable "+k)))
	}

	for k := range plugins.DisabledPlugins {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/pluginenable "+k)))
	}

	replyKeyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return telegram.SendCustom(update.Message.Chat.ID, 0, "test", false, &replyKeyboard)
}
