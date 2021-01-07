package echo

import (
	"github.com/ad/corpobot/plugins"
	"github.com/ad/corpobot/telegram"

	database "github.com/ad/corpobot/db"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Plugin struct{}

func init() {
	plugins.RegisterPlugin(&Plugin{})
}

func (m *Plugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("echo.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("echo", "example plugin", []string{"new", "member", "admin", "owner"}, echo)
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[echo.Plugin] Stopped")

	plugins.UnregisterCommand("echo")
}

var echo plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "echo", user.Role) {
		return true, telegram.Send(user.TelegramID, command+" "+args)
	}

	return false, nil
}
