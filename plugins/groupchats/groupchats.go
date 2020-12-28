package groupchats

import (
	"strconv"
	"strings"

	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	telegram "github.com/ad/corpobot/telegram"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Plugin ...
type Plugin struct {
}

func init() {
	plugins.RegisterPlugin(&Plugin{})
}

func (m *Plugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("groupchats.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("groupchatlist", "...", []string{"member", "admin", "owner"})
	plugins.RegisterCommand("groupchatinvitegenerate", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("groupchatuserban", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("groupchatuserunban", "...", []string{"admin", "owner"})
	plugins.RegisterCommand("groupchatmembers", "...", []string{"admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[groupchats.Plugin] Stopped")

	plugins.UnregisterCommand("groupchatlist")
	plugins.UnregisterCommand("groupchatinvitegenerate")
	plugins.UnregisterCommand("groupchatuserban")
	plugins.UnregisterCommand("groupchatuserunban")
	plugins.UnregisterCommand("groupchatmembers")
}

func (m *Plugin) Run(update *tgbotapi.Update, user *database.User) (bool, error) {
	args := strings.TrimSpace(update.Message.CommandArguments())

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "groupchatlist", user.Role) {
		return groupchatList(update, args)
	}

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "groupchatsinvitegenerate", user.Role) {
		return groupchatInviteGenerate(update, args)
	}

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "groupchatuserban", user.Role) {
		return groupchatUserBan(update, args)
	}

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "groupchatuserunban", user.Role) {
		return groupchatUserUnban(update, args)
	}

	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "groupchatmembers", user.Role) {
		return groupchatMembers(update, args)
	}

	return false, nil
}

func groupchatList(update *tgbotapi.Update, args string) (bool, error) {
	groupchats, err := database.GetGroupchats(plugins.DB, strings.Fields(args))
	if err != nil {
		return true, err
	}

	if len(groupchats) > 0 {
		var groupchatsList []string

		for _, u := range groupchats {
			groupchatsList = append(groupchatsList, "• "+u.String())
		}

		return true, telegram.Send(update.Message.Chat.ID, strings.Join(groupchatsList, "\n"))
	}

	return true, telegram.Send(update.Message.Chat.ID, "groupchat list is empty")
}

func groupchatInviteGenerate(update *tgbotapi.Update, args string) (bool, error) {
	if args == "" {
		return true, telegram.Send(update.Message.Chat.ID, "failed: empty groupchat ID")
	}

	telegramID, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
	}

	groupchat := &database.Groupchat{
		TelegramID: telegramID,
	}

	inviteLink, err := plugins.Bot.GetInviteLink(tgbotapi.ChatConfig{ChatID: groupchat.TelegramID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
	}

	groupchat.InviteLink = inviteLink
	if groupchat.InviteLink != "" {
		_, err := database.UpdateGroupChatInviteLink(plugins.DB, groupchat)
		if err != nil {
			return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
		}
	}

	return true, telegram.Send(update.Message.Chat.ID, "success")
}

func groupchatUserBan(update *tgbotapi.Update, args string) (bool, error) {
	errorString := "failed: you must provide the IDs of the user ans groupchat with a new line between them"

	params := strings.Split(args, "\n")

	if len(params) != 2 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	userIDstring, groupchatIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if userIDstring == "" || groupchatIDstring == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	var userID int
	if n, err := strconv.Atoi(userIDstring); err == nil {
		userID = n
	}

	var groupchatID int64
	n, err := strconv.ParseInt(groupchatIDstring, 10, 64)
	if err == nil {
		groupchatID = n
	}

	if userID == 0 || groupchatID == 0 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	_, err = plugins.Bot.KickChatMember(
		tgbotapi.KickChatMemberConfig{
			ChatMemberConfig: tgbotapi.ChatMemberConfig{
				ChatID: groupchatID,
				UserID: userID,
			},
		},
	)
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
	}

	return true, telegram.Send(update.Message.Chat.ID, "success")
}

func groupchatUserUnban(update *tgbotapi.Update, args string) (bool, error) {
	errorString := "failed: you must provide the IDs of the user ans groupchat with a new line between them"

	params := strings.Split(args, "\n")

	if len(params) != 2 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	userIDstring, groupchatIDstring := strings.TrimSpace(params[0]), strings.TrimSpace(params[1])

	if userIDstring == "" || groupchatIDstring == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	var userID int
	if n, err := strconv.Atoi(userIDstring); err == nil {
		userID = n
	}

	var groupchatID int64
	n, err := strconv.ParseInt(groupchatIDstring, 10, 64)
	if err == nil {
		groupchatID = n
	}

	if userID == 0 || groupchatID == 0 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	_, err = plugins.Bot.UnbanChatMember(tgbotapi.ChatMemberConfig{ChatID: groupchatID, UserID: userID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
	}

	return true, telegram.Send(update.Message.Chat.ID, "success")
}

func groupchatMembers(update *tgbotapi.Update, args string) (bool, error) {
	errorString := "failed: you must provide the groupchat ID"

	if args == "" {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	var groupchatID int64
	n, err := strconv.ParseInt(args, 10, 64)
	if err == nil {
		groupchatID = n
	}

	if groupchatID == 0 {
		return true, telegram.Send(update.Message.Chat.ID, errorString)
	}

	result, err := plugins.Bot.GetChatAdministrators(tgbotapi.ChatConfig{ChatID: groupchatID})
	if err != nil {
		return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
	}

	if len(result) > 0 {
		var usersList []string

		for _, u := range result {
			usersList = append(usersList, "@"+u.User.UserName+" "+u.User.FirstName+" "+u.User.LastName+" ["+strconv.Itoa(u.User.ID)+"] ("+u.Status+")")
		}
		return true, telegram.Send(update.Message.Chat.ID, strings.Join(usersList, "\n"))
	}

	return true, telegram.Send(update.Message.Chat.ID, "users not found")
}
