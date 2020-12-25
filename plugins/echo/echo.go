package echo

import (
	"github.com/ad/corpobot/plugins"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	telegram "github.com/ad/corpobot/telegram"
)

type EchoPlugin struct {

}

func init() {
	plugins.RegisterPlugin(&EchoPlugin{})
}

func (m *EchoPlugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("echo.EchoPlugin", "enabled") {
		return
	}

	plugins.RegisterCommand("echo", "example plugin")
}

func (m *EchoPlugin) OnStop() {
	dlog.Debugln("[EchoPlugin] Stopped")

	plugins.UnregisterCommand("echo")
}

func (m *EchoPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "echo" {
		return true, telegram.Send(update.Message.Chat.ID, update.Message.Text + " " + update.Message.CommandArguments())
	}

	return false, nil
}