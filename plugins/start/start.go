package start

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
	if !plugins.CheckIfPluginDisabled("start.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("start", "...")
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[start.Plugin] Stopped")

	plugins.UnregisterCommand("start")
}

func (m *Plugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "start" {
		return true, telegram.Send(update.Message.Chat.ID, "Hello! Send /help")
	}

	return false, nil
}
