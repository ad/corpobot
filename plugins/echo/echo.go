package echo

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
	if !plugins.CheckIfPluginDisabled("echo.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("echo", "example plugin", []string{"new", "member", "admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[echo.Plugin] Stopped")

	plugins.UnregisterCommand("echo")
}

func (m *Plugin) Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "echo", user.Role) {
		return true, telegram.Send(user.TelegramID, command+" "+args)
	}

	return false, nil
}
