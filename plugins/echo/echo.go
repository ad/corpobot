package echo

import (
	"github.com/ad/corpobot/plugins"

	database "github.com/ad/corpobot/db"
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
	plugin := &database.Plugin{
		Name: "echo.EchoPlugin",
		State: "enabled",
	}

	plugin, err := database.AddPluginIfNotExist(plugins.DB, plugin)
	if err != nil {
		dlog.Errorln("failed: " + err.Error())
	}

	if plugin.State != "enabled" {
		dlog.Debugln("[EchoPlugin] Disabled")
		return
	}

	dlog.Debugln("[EchoPlugin] Started")

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