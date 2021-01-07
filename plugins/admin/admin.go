package admin

import (
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
	if !plugins.CheckIfPluginDisabled("admin.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("pluginlist", "List of plugins", []string{"owner"}, pluginList)
	plugins.RegisterCommand("pluginenable", "Enable plugin", []string{"owner"}, pluginEnable)
	plugins.RegisterCommand("plugindisable", "Disable plugin", []string{"owner"}, pluginDisable)
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[admin.Plugin] Stopped")

	plugins.UnregisterCommand("pluginlist")
	plugins.UnregisterCommand("pluginenable")
	plugins.UnregisterCommand("plugindisable")
}

var pluginList plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	replyKeyboard := listPlugins()
	return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
}

var pluginEnable plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	plugin := &database.Plugin{
		Name:  args,
		State: "enabled",
	}

	_, err := database.UpdatePluginState(plugins.DB, plugin)
	if err != nil {
		return true, telegram.Send(user.TelegramID, "failed: "+err.Error())
	}
	if plugins.EnablePlugin(args) {
		if update.CallbackQuery != nil {
			_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" enabled"))
			if err != nil {
				dlog.Errorln(err.Error())
			}

			replyKeyboard := listPlugins()
			edit := tgbotapi.EditMessageReplyMarkupConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      update.CallbackQuery.Message.Chat.ID,
					MessageID:   update.CallbackQuery.Message.MessageID,
					ReplyMarkup: &replyKeyboard,
				},
			}

			_, err = plugins.Bot.Send(edit)
			return true, err
		}

		return true, telegram.Send(user.TelegramID, args+" enabled")
	}

	replyKeyboard := listPlugins()
	return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
}

var pluginDisable plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	plugin := &database.Plugin{
		Name:  args,
		State: "disabled",
	}

	_, err := database.UpdatePluginState(plugins.DB, plugin)
	if err != nil {
		return true, telegram.Send(user.TelegramID, "failed: "+err.Error())
	}
	if plugins.DisablePlugin(args) {
		if update.CallbackQuery != nil {
			_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" disabled"))
			if err != nil {
				dlog.Errorln(err.Error())
			}
			replyKeyboard := listPlugins()
			edit := tgbotapi.EditMessageReplyMarkupConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      update.CallbackQuery.Message.Chat.ID,
					MessageID:   update.CallbackQuery.Message.MessageID,
					ReplyMarkup: &replyKeyboard,
				},
			}

			_, err = plugins.Bot.Send(edit)
			return true, err
		}

		return true, telegram.Send(user.TelegramID, args+" disabled")
	}
	replyKeyboard := listPlugins()
	return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
}

func listPlugins() tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	plugins, err := database.GetPlugins(plugins.DB)
	if err != nil {
		dlog.Errorln(err.Error())
		return tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	for _, plugin := range plugins {
		if plugin.IsEnabled() {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("disable "+plugin.Name, "/plugindisable "+plugin.Name)))
		} else {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("enable "+plugin.Name, "/pluginenable "+plugin.Name)))
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
