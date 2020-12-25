package help

import (
	"bytes"
	"sort"

	"github.com/ad/corpobot/plugins"

	database "github.com/ad/corpobot/db"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	telegram "github.com/ad/corpobot/telegram"
)

type HelpPlugin struct {

}

func init() {
	plugins.RegisterPlugin(&HelpPlugin{})
}

func (m *HelpPlugin) OnStart() {
	plugin := &database.Plugin{
		Name: "help.HelpPlugin",
		State: "enabled",
	}

	plugin, err := database.AddPluginIfNotExist(plugins.DB, plugin)
	if err != nil {
		dlog.Errorln("failed: " + err.Error())
	}

	if plugin.State != "enabled" {
		dlog.Debugln("[HelpPlugin] Disabled")
		return
	}

	
	dlog.Debugln("[HelpPlugin] Started")

	plugins.RegisterCommand("help", "Display this help")
}

func (m *HelpPlugin) OnStop() {
	dlog.Debugln("[HelpPlugin] Stopped")

	plugins.UnregisterCommand("help")
}

func (m *HelpPlugin) Run(update *tgbotapi.Update) (bool, error) {
	if update.Message.Command() == "help" {
		var mk []string
		
		for k := range plugins.Commands {
			mk = append(mk, k)
		}

		sort.Strings(mk)
		var buffer bytes.Buffer

		for _, v := range mk {
			_, err := buffer.WriteString("/" + v + " - " + plugins.Commands[v] + "\n")
			if err != nil {
				return true, err
			}
		}

		return true, telegram.Send(update.Message.Chat.ID, "Those are my commands: \n"+buffer.String())
	}

	return false, nil
}