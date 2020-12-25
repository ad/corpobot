package echo

import (
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
	if !plugins.CheckIfPluginDisabled("echo.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("echo", "example plugin")
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[echo.Plugin] Stopped")

	plugins.UnregisterCommand("echo")
}

func (m *Plugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "echo" {
		return true, telegram.Send(update.Message.Chat.ID, update.Message.Text+" "+update.Message.CommandArguments())
	}

	return false, nil
}
