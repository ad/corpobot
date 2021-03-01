package groups

import (
	"strconv"
	"strings"

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
	if !plugins.CheckIfPluginDisabled("groups.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("grouplist", "Group list", []string{database.Member, database.Admin, database.Owner}, groupList)
	plugins.RegisterCommand("groupcreate", "Create group", []string{database.Admin, database.Owner}, groupCreate)
	plugins.RegisterCommand("grouprename", "Rename group", []string{database.Admin, database.Owner}, groupRename)
	plugins.RegisterCommand("groupdelete", "Delete group", []string{database.Admin, database.Owner}, groupDeleteUndelete)
	plugins.RegisterCommand("groupundelete", "Undelete group", []string{database.Admin, database.Owner}, groupDeleteUndelete)
	plugins.RegisterCommand("groupaddgroupchat", "Add groupchat to group", []string{database.Admin, database.Owner}, groupAddGroupChat)
	plugins.RegisterCommand("groupdeletegroupchat", "Delete groupchat from group", []string{database.Admin, database.Owner}, groupDeleteGroupChat)
	plugins.RegisterCommand("groupadduser", "Add user to group", []string{database.Admin, database.Owner}, groupAddUser)
	plugins.RegisterCommand("groupdeleteuser", "Delete user from group", []string{database.Admin, database.Owner}, groupDeleteUser)
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[groups.Plugin] Stopped")

	plugins.UnregisterCommand("grouplist")
	plugins.UnregisterCommand("groupcreate")
	plugins.UnregisterCommand("grouprename")
	plugins.UnregisterCommand("groupdelete")
	plugins.UnregisterCommand("groupundelete")
	plugins.UnregisterCommand("groupaddgroupchat")
	plugins.UnregisterCommand("groupdeletegroupchat")
	plugins.UnregisterCommand("groupadduser")
	plugins.UnregisterCommand("groupdeleteuser")
}

var groupList plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	groups, err := database.GetGroups(plugins.DB, strings.Fields(args))
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		return telegram.Send(user.TelegramID, "group list is empty")
	}

	var groupsList []string

	for _, u := range groups {
		groupsList = append(groupsList, "* "+u.String())

		groupchats, err := database.GetGroupchatsByGroupID(plugins.DB, u.ID)
		if err != nil {
			return err
		}
		for _, c := range groupchats {
			groupsList = append(groupsList, "    * "+c.String())
		}
	}

	return telegram.Send(user.TelegramID, strings.Join(groupsList, "\n"))
}

var groupCreate plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	if args == "" {
		return telegram.Send(user.TelegramID, "failed: empty group name")
	}

	group := &database.Group{
		Name: args,
	}

	_, err := database.AddGroupIfNotExist(plugins.DB, group)
	if err != nil {
		return telegram.Send(user.TelegramID, "failed: "+err.Error())
	}

	return telegram.Send(user.TelegramID, "group created")
}

var groupRename plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	names := strings.Split(args, "\n")

	if len(names) != 2 {
		return telegram.Send(user.TelegramID, "failed: you must provide the names of the two groups with a new line between them")
	}

	oldName, newName := strings.TrimSpace(names[0]), strings.TrimSpace(names[1])

	if oldName == "" || newName == "" {
		return telegram.Send(user.TelegramID, "failed: you must provide the names of the two groups with a new line between them")
	}

	rows, err := database.UpdateGroupName(plugins.DB, oldName, newName)
	if err != nil {
		return err
	}

	if rows != 1 {
		return telegram.Send(user.TelegramID, "failed")
	}

	return telegram.Send(user.TelegramID, "success")
}

var groupDeleteUndelete plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	newState := "active"

	if command == "groupdelete" {
		newState = "deleted"
	}

	group := &database.Group{
		Name:  args,
		State: newState,
	}

	rows, err := database.UpdateGroupState(plugins.DB, group)
	if err != nil {
		return err
	}

	if rows != 1 {
		return telegram.Send(user.TelegramID, "failed")
	}

	return telegram.Send(user.TelegramID, "success")
}

var groupAddGroupChat plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and groupchat id) with a new line between them"

	if len(params) != 2 {
		return telegram.Send(user.TelegramID, errorString)
	}

	groupName, groupchatIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || groupchatIDstring == "" {
		return telegram.Send(user.TelegramID, errorString)
	}

	groupchatID, err := strconv.ParseInt(groupchatIDstring, 10, 64)
	if err != nil {
		return telegram.Send(user.TelegramID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	groupchat, err := database.GetGroupChatByTelegramID(plugins.DB, &database.Groupchat{TelegramID: groupchatID})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	_, err = database.AddGroupGroupChatIfNotExist(plugins.DB, group, groupchat)
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	return telegram.Send(user.TelegramID, "success")
}

var groupDeleteGroupChat plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and groupchat id) with a new line between them"

	if len(params) != 2 {
		return telegram.Send(user.TelegramID, errorString)
	}

	groupName, groupchatIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || groupchatIDstring == "" {
		return telegram.Send(user.TelegramID, errorString)
	}

	groupchatID, err := strconv.ParseInt(groupchatIDstring, 10, 64)
	if err != nil {
		return telegram.Send(user.TelegramID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	groupchat, err := database.GetGroupChatByTelegramID(plugins.DB, &database.Groupchat{TelegramID: groupchatID})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	_, err = database.DeleteGroupGroupChat(plugins.DB, group, groupchat)
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	return telegram.Send(user.TelegramID, "success")
}

var groupAddUser plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and user id) with a new line between them"

	if len(params) != 2 {
		return telegram.Send(user.TelegramID, errorString)
	}

	groupName, userIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || userIDstring == "" {
		return telegram.Send(user.TelegramID, errorString)
	}

	userID, err := strconv.ParseInt(userIDstring, 10, 64)
	if err != nil {
		return telegram.Send(user.TelegramID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	userFromDB, err := database.GetUserByTelegramID(plugins.DB, &database.User{TelegramID: userID})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	_, err = database.AddGroupUserIfNotExist(plugins.DB, group, userFromDB)
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	return telegram.Send(user.TelegramID, "success")
}

var groupDeleteUser plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and user id) with a new line between them"

	if len(params) != 2 {
		return telegram.Send(user.TelegramID, errorString)
	}

	groupName, userIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || userIDstring == "" {
		return telegram.Send(user.TelegramID, errorString)
	}

	userID, err := strconv.ParseInt(userIDstring, 10, 64)
	if err != nil {
		return telegram.Send(user.TelegramID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	userFromDB, err := database.GetUserByTelegramID(plugins.DB, &database.User{TelegramID: userID})
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	_, err = database.DeleteGroupUser(plugins.DB, group, userFromDB)
	if err != nil {
		return telegram.Send(user.TelegramID, err.Error())
	}

	return telegram.Send(user.TelegramID, "success")
}
