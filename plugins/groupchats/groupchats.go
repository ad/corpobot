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

// GroupchatsPlugin ...
type GroupchatsPlugin struct {
}

func init() {
	plugins.RegisterPlugin(&GroupchatsPlugin{})
}

func (m *GroupchatsPlugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("groupchats.GroupchatsPlugin", "enabled") {
		return
	}

	plugins.RegisterCommand("groupchatlist", "...")
	plugins.RegisterCommand("groupchatsinvitegenerate", "...")
}

func (m *GroupchatsPlugin) OnStop() {
	dlog.Debugln("[GroupchatsPlugin] Stopped")

	plugins.UnregisterCommand("groupchatlist")
	plugins.UnregisterCommand("groupchatsinvitegenerate")
}

func (m *GroupchatsPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "groupchatlist" {
		// TODO: check user rights

		args := update.Message.CommandArguments()

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

	if update.Message.Command() == "groupchatsinvitegenerate" {
		// TODO: check user rights

		args := strings.TrimSpace(update.Message.CommandArguments())

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

		// if update.Message.Chat.InviteLink == "" {
		inviteLink, err := plugins.Bot.GetInviteLink(tgbotapi.ChatConfig{ChatID: groupchat.TelegramID})
		if err != nil {
			return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
		} else {
			groupchat.InviteLink = inviteLink
		}

		// }
		if groupchat.InviteLink != "" {
			_, err := database.UpdateGroupChatInviteLink(plugins.DB, groupchat)
			if err != nil {
				return true, telegram.Send(update.Message.Chat.ID, "failed: "+err.Error())
			}
		}

		return true, telegram.Send(update.Message.Chat.ID, "success")
	}

	return false, nil
}
