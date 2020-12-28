package start

import (
	"github.com/ad/corpobot/plugins"

	database "github.com/ad/corpobot/db"
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
	if !plugins.CheckIfPluginDisabled("start.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("start", "...", []string{"new", "member", "admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[start.Plugin] Stopped")

	plugins.UnregisterCommand("start")
}

func (m *Plugin) Run(update *tgbotapi.Update, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(update.Message.Command(), "start", user.Role) {
		return true, telegram.Send(update.Message.Chat.ID, "Hello! Send /help")
	}

	return false, nil
}
