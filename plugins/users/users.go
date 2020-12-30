package users

import (
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

	plugins.RegisterCommand("userlist", "...", []string{"member", "admin", "owner"})
	plugins.RegisterCommand("userpromote", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("userblock", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("userdelete", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("userunblock", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("userundelete", "...", []string{"admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[users.Plugin] Stopped")

	plugins.UnregisterCommand("userlist")
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
	users, err := database.GetUsers(plugins.DB, strings.Fields(args))
	if err != nil {
		return true, err
	}

	if len(users) > 0 {
		var usersList []string

		for _, u := range users {
			usersList = append(usersList, "• "+u.String())
		}

		return true, telegram.Send(user.TelegramID, strings.Join(usersList, "\n"))
	}

	return true, telegram.Send(user.TelegramID, "user list is empty")
}

func userBlockUnblock(update *tgbotapi.Update, user *database.User, command, args string) (bool, error) {
	newRole := "member"

	if command == "userblock" {
		newRole = "blocked"
	}

	args = strings.TrimLeft(args, "@")

	u := &database.User{
		UserName: args,
		Role:     newRole,
	}

	rows, err := database.UpdateUserRole(plugins.DB, u)
	if err != nil {
		return true, err
	}

	if rows != 1 {
		return true, telegram.Send(user.TelegramID, "failed")
	}

	return true, telegram.Send(user.TelegramID, "success")
}

func userDeleteUndelete(update *tgbotapi.Update, user *database.User, command, args string) (bool, error) {
	newRole := "member"

	if command == "userdelete" {
		newRole = "deleted"
	}

	args = strings.TrimLeft(args, "@")

	u := &database.User{
		UserName: args,
		Role:     newRole,
	}

	rows, err := database.UpdateUserRole(plugins.DB, u)
	if err != nil {
		return true, err
	}

	if rows != 1 {
		return true, telegram.Send(user.TelegramID, "failed")
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

	if telegramIDstring == "" || newRole == "" || newRole == "owner" {
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
		Role:     newRole,
	}

	rows, err := database.UpdateUserRole(plugins.DB, u)
	if err != nil {
		return true, err
	}

	if rows != 1 {
		return true, telegram.Send(user.TelegramID, "failed")
	}

	return true, telegram.Send(user.TelegramID, "success")
}
