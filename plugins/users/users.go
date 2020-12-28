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

func (m *Plugin) Run(update *tgbotapi.Update, command string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "userlist", user.Role) {
		args := telegram.GetArguments(update)

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

	if plugins.CheckIfCommandIsAllowed(command, "userblock", user.Role) || plugins.CheckIfCommandIsAllowed(command, "userunblock", user.Role) {
		newRole := "member"

		if command == "userblock" {
			newRole = "blocked"
		}

		args := strings.TrimLeft(update.Message.CommandArguments(), "@")

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

	if plugins.CheckIfCommandIsAllowed(command, "userdelete", user.Role) || plugins.CheckIfCommandIsAllowed(command, "userundelete", user.Role) {
		newRole := "member"

		if command == "userdelete" {
			newRole = "deleted"
		}

		args := strings.TrimLeft(update.Message.CommandArguments(), "@")

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

	return false, nil
}
