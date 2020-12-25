package users

import (
	"strings"

	"github.com/ad/corpobot/plugins"
	database "github.com/ad/corpobot/db"
	telegram "github.com/ad/corpobot/telegram"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type UsersPlugin struct {

}

func init() {
	plugins.RegisterPlugin(&UsersPlugin{})
}

func (m *UsersPlugin) OnStart() {
	dlog.Debugln("[UsersPlugin] Started")

	plugins.RegisterCommand("userlist", "...")
	plugins.RegisterCommand("userpromote", "...")
	plugins.RegisterCommand("userblock", "...")
	plugins.RegisterCommand("userdelete", "...")
	plugins.RegisterCommand("userunblock", "...")
	plugins.RegisterCommand("userundelete", "...")
}

func (m *UsersPlugin) OnStop() {
	dlog.Debugln("[UsersPlugin] Stopped")

	plugins.UnregisterCommand("userlist")
	plugins.UnregisterCommand("userpromote")
	plugins.UnregisterCommand("userblock")
	plugins.UnregisterCommand("userdelete")
	plugins.UnregisterCommand("userunblock")
	plugins.UnregisterCommand("userundelete")
}

func (m *UsersPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "userlist" {
		// TODO: check user rights

		args := update.Message.CommandArguments()

		users, err := database.GetUsers(plugins.DB, strings.Fields(args))
		if err != nil {
			return true, err
		}

		if len(users) > 0 {
			var usersList []string

			for _, u := range users {
				usersList = append(usersList, "• " + u.String())
			}

			return true, telegram.Send(update.Message.Chat.ID, strings.Join(usersList, "\n"))
		}

		return true, telegram.Send(update.Message.Chat.ID, "user list is empty")
	}

	if update.Message.Command() == "userblock" || update.Message.Command() == "userunblock" {
		// TODO: check user rights

		var newRole = "member"

		if update.Message.Command() == "userblock" {
			newRole = "blocked"
		}

		args := strings.TrimLeft(update.Message.CommandArguments(), "@")

		user := &database.User{
			UserName: 	args,
			Role: 		newRole,
		}

		rows, err := database.UpdateUserRole(plugins.DB, user)
		if err != nil {
			return true, err
		}

		if rows != 1 {
			return true, telegram.Send(update.Message.Chat.ID, update.Message.Command() + " failed")
		}

		return true, telegram.Send(update.Message.Chat.ID, update.Message.Command() + " success")
	}

	if update.Message.Command() == "userdelete" || update.Message.Command() == "userundelete" {
		// TODO: check user rights

		var newRole = "member"

		if update.Message.Command() == "userdelete" {
			newRole = "deleted"
		}

		args := strings.TrimLeft(update.Message.CommandArguments(), "@")

		user := &database.User{
			UserName: 	args,
			Role: 		newRole,
		}

		rows, err := database.UpdateUserRole(plugins.DB, user)
		if err != nil {
			return true, err
		}

		if rows != 1 {
			return true, telegram.Send(update.Message.Chat.ID, update.Message.Command() + " failed")
		}

		return true, telegram.Send(update.Message.Chat.ID, update.Message.Command() + " success")
	}

	return false, nil
}
