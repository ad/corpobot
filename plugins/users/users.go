package users

import (
	"strconv"
	"strings"

	database "github.com/ad/corpobot/db"
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
	if !plugins.CheckIfPluginDisabled("users.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("userlist", "User list", []string{database.Member, database.Admin, database.Owner})
	plugins.RegisterCommand("user", "User actions", []string{database.Admin, database.Owner})
	plugins.RegisterCommand("userpromote", "Change user role", []string{database.Admin, database.Owner})
	plugins.RegisterCommand("userblock", "Block user", []string{database.Admin, database.Owner})
	plugins.RegisterCommand("userdelete", "Delete user", []string{database.Admin, database.Owner})
	plugins.RegisterCommand("userunblock", "Unblock user", []string{database.Admin, database.Owner})
	plugins.RegisterCommand("userundelete", "Undelete user", []string{database.Admin, database.Owner})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[users.Plugin] Stopped")

	plugins.UnregisterCommand("userlist")
	plugins.UnregisterCommand("user")
	plugins.UnregisterCommand("userpromote")
	plugins.UnregisterCommand("userblock")
	plugins.UnregisterCommand("userdelete")
	plugins.UnregisterCommand("userunblock")
	plugins.UnregisterCommand("userundelete")
}

func (m *Plugin) Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "userlist", user.Role) {
		return userList(update, user, args)
	}

	if plugins.CheckIfCommandIsAllowed(command, "user", user.Role) {
		return userActions(update, user, args)
	}

	if plugins.CheckIfCommandIsAllowed(command, "userblock", user.Role) || plugins.CheckIfCommandIsAllowed(command, "userunblock", user.Role) {
		return userBlockUnblock(update, user, command, args)
	}

	if plugins.CheckIfCommandIsAllowed(command, "userdelete", user.Role) || plugins.CheckIfCommandIsAllowed(command, "userundelete", user.Role) {
		return userDeleteUndelete(update, user, command, args)
	}

	if plugins.CheckIfCommandIsAllowed(command, "userpromote", user.Role) {
		return userPromote(update, user, args)
	}

	return false, nil
}

func userList(update *tgbotapi.Update, user *database.User, args string) (bool, error) {
	replyKeyboard := ListUsers(update, user, args)
	return true, telegram.SendCustom(user.TelegramID, 0, "Choose user", false, &replyKeyboard)
}

func userActions(update *tgbotapi.Update, user *database.User, args string) (bool, error) {
	userID, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return true, telegram.Send(user.TelegramID, "wrong telegramID provided")
	}

	replyKeyboard := userActionsList(update, user, userID)
	if update.CallbackQuery != nil {
		_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" success"))
		if err != nil {
			dlog.Errorln(err.Error())
		}

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

	return true, telegram.SendCustom(user.TelegramID, 0, "Choose action", false, &replyKeyboard)
}

func userBlockUnblock(update *tgbotapi.Update, user *database.User, command, args string) (bool, error) {
	newRole := database.Member

	if command == "userblock" {
		newRole = database.Blocked
	}

	var telegramID int64
	n, err := strconv.ParseInt(args, 10, 64)
	if err == nil {
		telegramID = n
	}

	if telegramID == 0 {
		return true, telegram.Send(user.TelegramID, "please provide user telegram ID")
	}

	u := &database.User{
		TelegramID: telegramID,
		Role:       newRole,
	}

	rows, err := database.UpdateUserRole(plugins.DB, u)
	if err != nil {
		return true, err
	}

	if rows != 1 {
		return true, telegram.Send(user.TelegramID, "failed")
	}

	if update.CallbackQuery != nil {
		_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" success"))
		if err != nil {
			dlog.Errorln(err.Error())
		}

		replyKeyboard := userActionsList(update, user, telegramID)
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

	return true, telegram.Send(user.TelegramID, "success")
}

func userDeleteUndelete(update *tgbotapi.Update, user *database.User, command, args string) (bool, error) {
	newRole := database.Member

	if command == "userdelete" {
		newRole = database.Deleted
	}

	var telegramID int64
	n, err := strconv.ParseInt(args, 10, 64)
	if err == nil {
		telegramID = n
	}

	if telegramID == 0 {
		return true, telegram.Send(user.TelegramID, "please provide user telegram ID")
	}

	u := &database.User{
		TelegramID: telegramID,
		Role:       newRole,
	}

	rows, err := database.UpdateUserRole(plugins.DB, u)
	if err != nil {
		return true, err
	}

	if rows != 1 {
		return true, telegram.Send(user.TelegramID, "failed")
	}

	if update.CallbackQuery != nil {
		_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" success"))
		if err != nil {
			dlog.Errorln(err.Error())
		}

		replyKeyboard := userActionsList(update, user, telegramID)
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

	return true, telegram.Send(user.TelegramID, "success")
}

func userPromote(update *tgbotapi.Update, user *database.User, args string) (bool, error) {
	errorString := "failed: you must provide TelegramID and new role with a new line between them"

	params := strings.Split(args, "\n")

	if len(params) != 2 {
		return true, telegram.Send(user.TelegramID, errorString)
	}

	telegramIDstring, newRole := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if telegramIDstring == "" || newRole == "" || newRole == database.Owner {
		return true, telegram.Send(user.TelegramID, errorString)
	}

	var telegramID int64
	n, err := strconv.ParseInt(telegramIDstring, 10, 64)
	if err == nil {
		telegramID = n
	}

	if telegramID == 0 {
		return true, telegram.Send(user.TelegramID, errorString)
	}

	u := &database.User{
		TelegramID: telegramID,
		Role:       newRole,
	}

	rows, err := database.UpdateUserRole(plugins.DB, u)
	if err != nil {
		return true, err
	}

	if rows != 1 {
		return true, telegram.Send(user.TelegramID, "failed")
	}

	if update.CallbackQuery != nil {
		_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" enabled"))
		if err != nil {
			dlog.Errorln(err.Error())
		}

		replyKeyboard := userActionsList(update, user, telegramID)
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

	return true, telegram.Send(user.TelegramID, "success")
}

func ListUsers(update *tgbotapi.Update, user *database.User, args string) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	users, err := database.GetUsers(plugins.DB, strings.Fields(args))
	if err != nil {
		dlog.Errorln(err.Error())
		return tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	for _, u := range users {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(u.String(), "/user "+strconv.FormatInt(u.TelegramID, 10))))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func userActionsList(update *tgbotapi.Update, user *database.User, userID int64) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	userFromDB, err := database.GetUserByTelegramID(plugins.DB, &database.User{TelegramID: userID})
	if err != nil {
		dlog.Errorln(err.Error())
		return tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	switch userFromDB.Role {
	case database.Deleted:
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("undelete user", "/userundelete "+strconv.FormatInt(userFromDB.TelegramID, 10))))
	case database.Blocked:
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("unblock user", "/userunblock "+strconv.FormatInt(userFromDB.TelegramID, 10))))
	case database.New:
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("block user", "/userblock "+strconv.FormatInt(userFromDB.TelegramID, 10))))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("delete user", "/userdelete "+strconv.FormatInt(userFromDB.TelegramID, 10))))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("make member", "/userpromote "+strconv.FormatInt(userFromDB.TelegramID, 10)+"\nmember")))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("make admin", "/userpromote "+strconv.FormatInt(userFromDB.TelegramID, 10)+"\nadmin")))
	case database.Member:
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("block user", "/userblock "+strconv.FormatInt(userFromDB.TelegramID, 10))))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("delete user", "/userdelete "+strconv.FormatInt(userFromDB.TelegramID, 10))))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("make admin", "/userpromote "+strconv.FormatInt(userFromDB.TelegramID, 10)+"\nadmin")))
	case database.Admin:
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("block user", "/userblock "+strconv.FormatInt(userFromDB.TelegramID, 10))))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("delete user", "/userdelete "+strconv.FormatInt(userFromDB.TelegramID, 10))))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("make member", "/userpromote "+strconv.FormatInt(userFromDB.TelegramID, 10)+"\nmember")))
	default:
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("user actions", "/user "+strconv.FormatInt(userFromDB.TelegramID, 10))))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
