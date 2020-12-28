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

func (m *Plugin) Run(update *tgbotapi.Update, command string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "start", user.Role) {
		return true, telegram.Send(user.TelegramID, "Hello! Send /help")
	}

	return false, nil
}
