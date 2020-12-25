package groups

import (
	"strings"

	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	telegram "github.com/ad/corpobot/telegram"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type GroupsPlugin struct {
}

func init() {
	plugins.RegisterPlugin(&GroupsPlugin{})
}

func (m *GroupsPlugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("groups.GroupsPlugin", "enabled") {
		return
	}

	plugins.RegisterCommand("grouplist", "...")
	plugins.RegisterCommand("groupcreate", "...")
	plugins.RegisterCommand("grouprename", "...")
	plugins.RegisterCommand("groupdelete", "...")
	plugins.RegisterCommand("groupundelete", "...")
}

func (m *GroupsPlugin) OnStop() {
	dlog.Debugln("[GroupsPlugin] Stopped")

	plugins.UnregisterCommand("grouplist")
	plugins.UnregisterCommand("groupcreate")
	plugins.UnregisterCommand("grouprename")
	plugins.UnregisterCommand("groupdelete")
	plugins.UnregisterCommand("groupundelete")
}

func (m *GroupsPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "grouplist" {
		// TODO: check user rights

		args := update.Message.CommandArguments()

		groups, err := database.GetGroups(plugins.DB, strings.Fields(args))
		if err != nil {
			return true, err
		}

		if len(groups) > 0 {
			var groupsList []string

			for _, u := range groups {
				groupsList = append(groupsList, "• "+u.String())
			}

			return true, telegram.Send(update.Message.Chat.ID, strings.Join(groupsList, "\n"))
		}

		return true, telegram.Send(update.Message.Chat.ID, "group list is empty")
	}

	if update.Message.Command() == "groupcreate" {
		// TODO: check user rights

		args := strings.TrimSpace(update.Message.CommandArguments())

		if args == "" {
			return true, telegram.Send(update.Message.Chat.ID, "failed: empty group name")
		}

		group := &database.Group{
			Name: args,
		}

		_, err := database.AddGroupIfNotExist(plugins.DB, group)
		if err != nil {
			return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
		}

		return true, telegram.Send(update.Message.Chat.ID, "group created")
	}

	if update.Message.Command() == "grouprename" {
		// TODO: check user rights

		args := update.Message.CommandArguments()

		names := strings.Split(args, "\n")

		if len(names) != 2 {
			return true, telegram.Send(update.Message.Chat.ID, "failed: you must provide the names of the two groups with a new line between them")
		}

		oldName, newName := strings.TrimSpace(names[0]), strings.TrimSpace(names[1])

		if oldName == "" || newName == "" {
			return true, telegram.Send(update.Message.Chat.ID, "failed: you must provide the names of the two groups with a new line between them")
		}

		rows, err := database.UpdateGroupName(plugins.DB, oldName, newName)
		if err != nil {
			return true, err
		}

		if rows != 1 {
			return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" failed")
		}

		return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" success")
	}

	if update.Message.Command() == "groupdelete" || update.Message.Command() == "groupundelete" {
		// TODO: check user rights

		newState := "active"

		if update.Message.Command() == "groupdelete" {
			newState = "deleted"
		}

		args := update.Message.CommandArguments()

		group := &database.Group{
			Name:  args,
			State: newState,
		}

		rows, err := database.UpdateGroupState(plugins.DB, group)
		if err != nil {
			return true, err
		}

		if rows != 1 {
			return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" failed")
		}

		return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" success")
	}

	return false, nil
}
