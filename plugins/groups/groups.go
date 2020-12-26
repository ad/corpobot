package groups

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
	if !plugins.CheckIfPluginDisabled("groups.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("grouplist", "...")
	plugins.RegisterCommand("groupcreate", "...")
	plugins.RegisterCommand("grouprename", "...")
	plugins.RegisterCommand("groupdelete", "...")
	plugins.RegisterCommand("groupundelete", "...")
	plugins.RegisterCommand("groupaddgroupchat", "...")
	plugins.RegisterCommand("groupdeletegroupchat", "...")
	plugins.RegisterCommand("groupadduser", "...")
	plugins.RegisterCommand("groupdeleteuser", "...")
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

func (m *Plugin) Run(update *tgbotapi.Update) (bool, error) {
	args := strings.TrimSpace(update.Message.CommandArguments())

	if update.Message.Command() == "grouplist" {
		return groupList(update, args)
	}

	if update.Message.Command() == "groupcreate" {
		return groupCreate(update, args)
	}

	if update.Message.Command() == "grouprename" {
		return groupRename(update, args)
	}

	if update.Message.Command() == "groupdelete" || update.Message.Command() == "groupundelete" {
		return groupDeleteUndelete(update, args)
	}

	if update.Message.Command() == "groupaddgroupchat" {
		return groupAddGroupChat(update, args)
	}

	if update.Message.Command() == "groupdeletegroupchat" {
		return groupDeleteGroupChat(update, args)
	}

	if update.Message.Command() == "groupadduser" {
		return groupAddUser(update, args)
	}

	if update.Message.Command() == "groupdeleteuser" {
		return groupDeleteUser(update, args)
	}

	return false, nil
}

func groupList(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

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

func groupCreate(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

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

func groupRename(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

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

func groupDeleteUndelete(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

	newState := "active"

	if update.Message.Command() == "groupdelete" {
		newState = "deleted"
	}

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

func groupAddGroupChat(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and groupchat id) with a new line between them"

	if len(params) != 2 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	groupName, groupchatIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || groupchatIDstring == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	groupchatID, err := strconv.ParseInt(groupchatIDstring, 10, 64)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	groupchat, err := database.GetGroupChatByTelegramID(plugins.DB, &database.Groupchat{TelegramID: groupchatID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	_, err = database.AddGroupGroupChatIfNotExist(plugins.DB, group, groupchat)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" success")
}

func groupDeleteGroupChat(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and groupchat id) with a new line between them"

	if len(params) != 2 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	groupName, groupchatIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || groupchatIDstring == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	groupchatID, err := strconv.ParseInt(groupchatIDstring, 10, 64)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	groupchat, err := database.GetGroupChatByTelegramID(plugins.DB, &database.Groupchat{TelegramID: groupchatID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	_, err = database.DeleteGroupGroupChat(plugins.DB, group, groupchat)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" success")
}

func groupAddUser(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and user id) with a new line between them"

	if len(params) != 2 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	groupName, userIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || userIDstring == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	userID, err := strconv.Atoi(userIDstring)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	user, err := database.GetUserByTelegramID(plugins.DB, &database.User{TelegramID: userID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	_, err = database.AddGroupUserIfNotExist(plugins.DB, group, user)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" success")
}

func groupDeleteUser(update *tgbotapi.Update, args string) (bool, error) {
	// TODO: check user rights

	params := strings.Split(args, "\n")

	errorString := "failed: you must provide two lines (group name and user id) with a new line between them"

	if len(params) != 2 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	groupName, userIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if groupName == "" || userIDstring == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	userID, err := strconv.Atoi(userIDstring)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	group, err := database.GetGroupByName(plugins.DB, &database.Group{Name: groupName})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	user, err := database.GetUserByTelegramID(plugins.DB, &database.User{TelegramID: userID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	_, err = database.DeleteGroupUser(plugins.DB, group, user)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, err.Error())
	}

	return true, telegram.Send(update.Message.Chat.ID, update.Message.Command()+" success")
}
