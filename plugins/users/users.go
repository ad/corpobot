package users

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cal "github.com/ad/corpobot/calendar"
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
	if !plugins.CheckIfPluginDisabled("users.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("userlist", "User list", []string{database.Member, database.Admin, database.Owner}, userList)
	plugins.RegisterCommand("user", "User actions", []string{database.Admin, database.Owner}, user)
	plugins.RegisterCommand("userpromote", "Change user role", []string{database.Admin, database.Owner}, userPromote)
	plugins.RegisterCommand("userblock", "Block user", []string{database.Admin, database.Owner}, userBlockUnblock)
	plugins.RegisterCommand("userdelete", "Delete user", []string{database.Admin, database.Owner}, userDeleteUndelete)
	plugins.RegisterCommand("userunblock", "Unblock user", []string{database.Admin, database.Owner}, userBlockUnblock)
	plugins.RegisterCommand("userundelete", "Undelete user", []string{database.Admin, database.Owner}, userDeleteUndelete)
	plugins.RegisterCommand("userbirthday", "Set user birthday", []string{database.Member, database.Admin, database.Owner}, userBirthday)
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
	plugins.UnregisterCommand("userbirthday")
}

var userList plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if !plugins.CheckIfCommandIsAllowed(command, "userlist", user.Role) {
		return false, nil
	}

	replyKeyboard := listUsers(args)
	return true, telegram.SendCustom(user.TelegramID, 0, "Choose user", false, &replyKeyboard)
}

var user plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if !plugins.CheckIfCommandIsAllowed(command, "user", user.Role) {
		return false, nil
	}

	userID, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return true, telegram.Send(user.TelegramID, "wrong telegramID provided")
	}

	replyKeyboard := userActionsList(userID)
	if update.CallbackQuery != nil {
		_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, args+" success"))
		if err != nil {
			dlog.Errorln(err.Error())
		}

		editKeyboard := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      update.CallbackQuery.Message.Chat.ID,
				MessageID:   update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &replyKeyboard,
			},
			Text: user.Paragraph(),
		}

		_, err = plugins.Bot.Send(editKeyboard)
		return true, err
	}

	return true, telegram.SendCustom(user.TelegramID, 0, user.Paragraph(), false, &replyKeyboard)
}

var userPromote plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if !plugins.CheckIfCommandIsAllowed(command, "userpromote", user.Role) {
		return false, nil
	}

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

		replyKeyboard := userActionsList(telegramID)
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

	errNotifyUser := telegram.Send(telegramID, "you were assigned the role \""+newRole+"\", user /help for command list")
	if errNotifyUser != nil {
		dlog.Errorln(errNotifyUser.Error())
	}

	return true, telegram.Send(user.TelegramID, "success")
}

var userBlockUnblock plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if !plugins.CheckIfCommandIsAllowed(command, "userblock", user.Role) && !plugins.CheckIfCommandIsAllowed(command, "userunblock", user.Role) {
		return false, nil
	}
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

		replyKeyboard := userActionsList(telegramID)
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

	errNotifyUser := telegram.Send(telegramID, "you were assigned the role \""+newRole+"\", user /help for command list")
	if errNotifyUser != nil {
		dlog.Errorln(errNotifyUser.Error())
	}

	return true, telegram.Send(user.TelegramID, "success")

}

var userDeleteUndelete plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if !plugins.CheckIfCommandIsAllowed(command, "userdelete", user.Role) && !plugins.CheckIfCommandIsAllowed(command, "userundelete", user.Role) {
		return false, nil
	}
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

		replyKeyboard := userActionsList(telegramID)
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

	errNotifyUser := telegram.Send(telegramID, "you were assigned the role \""+newRole+"\", user /help for command list")
	if errNotifyUser != nil {
		dlog.Errorln(errNotifyUser.Error())
	}

	return true, telegram.Send(user.TelegramID, "success")

}

var userBirthday plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if !plugins.CheckIfCommandIsAllowed(command, "userbirthday", user.Role) {
		return false, nil
	}

	if result, err := storeBirthday(update, user, args); result {
		return true, err
	}

	lang := telegram.GetLanguage(update)

	replyKeyboard := tgbotapi.InlineKeyboardMarkup{}

	switch {
	case strings.HasPrefix(args, "<"):
		date := strings.TrimLeft(args, "<")
		year, month, _, err := cal.ParseDate(date)
		if err == nil {
			replyKeyboard, _, _ = cal.HandlerPrevMonth("/userbirthday", year, time.Month(month), lang)
		}
	case strings.HasPrefix(args, ">"):
		date := strings.TrimLeft(args, ">")
		year, month, _, err := cal.ParseDate(date)
		if err == nil {
			replyKeyboard, _, _ = cal.HandlerNextMonth("/userbirthday", year, time.Month(month), lang)
		}
	case strings.HasPrefix(args, "«"):
		date := strings.TrimLeft(args, "«")
		year, month, _, err := cal.ParseDate(date)
		if err == nil {
			replyKeyboard, _, _ = cal.HandlerPrevYear("/userbirthday", year, time.Month(month), lang)
		}
	case strings.HasPrefix(args, "»"):
		date := strings.TrimLeft(args, "»")
		year, month, _, err := cal.ParseDate(date)
		if err == nil {
			replyKeyboard, _, _ = cal.HandlerNextYear("/userbirthday", year, time.Month(month), lang)
		}
	case strings.HasPrefix(args, "m"):
		currentTime := time.Now()
		year := currentTime.Year()
		month := currentTime.Month()

		date := strings.TrimLeft(args, "m")
		year2, month2, _, err := cal.ParseDate(date)
		if err == nil {
			year = year2
			month = time.Month(month2)
		}
		replyKeyboard = cal.GenerateMonths("/userbirthday", year, month, lang)
	case strings.HasPrefix(args, "y"):
		currentTime := time.Now()
		year := currentTime.Year()
		month := currentTime.Month()

		date := strings.TrimLeft(args, "y")
		year2, month2, _, err := cal.ParseDate(date)
		if err == nil {
			year = year2
			month = time.Month(month2)
		}
		replyKeyboard = cal.GenerateYears("/userbirthday", year, month, lang)
	default:
		currentTime := time.Now()
		year := currentTime.Year()
		month := currentTime.Month()

		year2, month2, _, err := cal.ParseDate(args)
		if err == nil {
			year = year2
			month = time.Month(month2)
		}
		replyKeyboard = cal.GenerateCalendar("/userbirthday", year, month, lang)
	}

	if update.CallbackQuery != nil {
		_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
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

	return true, telegram.SendCustom(user.TelegramID, 0, "Choose your birth date", false, &replyKeyboard)
}

func storeBirthday(update *tgbotapi.Update, user *database.User, args string) (bool, error) {
	year, month, day, err := cal.ParseDate(args)
	if err == nil && year != 0 && month != 0 && day != 0 {
		layout := "2006-01-02"
		t, err := time.Parse(layout, fmt.Sprintf("%d-%02d-%02d", year, month, day))
		if err != nil {
			return true, telegram.Send(user.TelegramID, "failed: "+err.Error())
		}

		u := &database.User{
			TelegramID: user.TelegramID,
			Birthday:   t,
		}

		_, err = database.UpdateUserBirthday(plugins.DB, u)
		if err != nil {
			return true, telegram.Send(user.TelegramID, "failed: "+err.Error())
		}

		var deleteMessage tgbotapi.DeleteMessageConfig
		if update.CallbackQuery != nil {
			deleteMessage = tgbotapi.DeleteMessageConfig{ChatID: update.CallbackQuery.Message.Chat.ID, MessageID: update.CallbackQuery.Message.MessageID}
		} else if update.Message != nil {
			deleteMessage = tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: update.Message.MessageID}
		}
		_, err = plugins.Bot.DeleteMessage(deleteMessage)
		if err != nil {
			dlog.Errorln(err.Error())
		}

		return true, telegram.Send(user.TelegramID, "Your birth date ("+args+") saved")
	}

	return false, nil
}

func listUsers(args string) tgbotapi.InlineKeyboardMarkup {
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

func userActionsList(userID int64) tgbotapi.InlineKeyboardMarkup {
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
