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

func (m *Plugin) Run(update *tgbotapi.Update, command string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "pluginlist", user.Role) {
		replyKeyboard := ListPlugins(update, user)
		return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
	}

	if plugins.CheckIfCommandIsAllowed(command, "pluginenable", user.Role) {
		args := telegram.GetArguments(update)
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

				replyKeyboard := ListPlugins(update, user)
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

		replyKeyboard := ListPlugins(update, user)
		return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
	}

	if plugins.CheckIfCommandIsAllowed(command, "plugindisable", user.Role) {
		args := telegram.GetArguments(update)
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
				replyKeyboard := ListPlugins(update, user)
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
		replyKeyboard := ListPlugins(update, user)
		return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
	}

	return false, nil
}

func ListPlugins(update *tgbotapi.Update, user *database.User) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	plugins, err := database.GetPlugins(plugins.DB)
	if err != nil {
		dlog.Errorln(err.Error())
		return tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	for _, plugin := range plugins {
		if plugin.State == "enabled" {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("disable "+plugin.Name, "/plugindisable "+plugin.Name)))
		} else {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("enable "+plugin.Name, "/pluginenable "+plugin.Name)))
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
