package start

import (
	"github.com/ad/corpobot/plugins"
	
	database "github.com/ad/corpobot/db"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	telegram "github.com/ad/corpobot/telegram"
)

type StartPlugin struct {

}

func init() {
	plugins.RegisterPlugin(&StartPlugin{})
}

func (m *StartPlugin) OnStart() {
	plugin := &database.Plugin{
		Name: "start.StartPlugin",
		State: "enabled",
	}

	plugin, err := database.AddPluginIfNotExist(plugins.DB, plugin)
	if err != nil {
		dlog.Errorln("failed: " + err.Error())
	}

	if plugin.State != "enabled" {
		dlog.Debugln("[StartPlugin] Disabled")
		return
	}

	
	dlog.Debugln("[StartPlugin] Started")

	plugins.RegisterCommand("start", "...")
}

func (m *StartPlugin) OnStop() {
	dlog.Debugln("[StartPlugin] Stopped")

	plugins.UnregisterCommand("start")
}

func (m *StartPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "start" {
		return true, telegram.Send(update.Message.Chat.ID, "Hello! Send /help")
	}

	return false, nil
}