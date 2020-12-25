package help

import (
	"bytes"
	"sort"

	"github.com/ad/corpobot/plugins"

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
	if !plugins.CheckIfPluginDisabled("help.HelpPlugin", "enabled") {
		return
	}

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