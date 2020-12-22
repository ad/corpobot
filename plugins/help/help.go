package help

import (
	"bytes"
	"sort"

	"github.com/ad/corpobot/plugins"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type HelpPlugin struct {
}

func init() {
	plugins.RegisterPlugin(&HelpPlugin{})
}

func (m *HelpPlugin) OnStart() {
	dlog.Debugln("[HelpPlugin] Started")
	plugins.RegisterCommand("help", "Display this help")
}

func (m *HelpPlugin) OnStop() {
	plugins.UnregisterCommand("help")

}

func (m *HelpPlugin) Run(update *tgbotapi.Update) {
	dlog.Debugln("[HelpPlugin] %s", update)

	if update.Message.Command() == "help" {
		mk := make([]string, len(plugins.Commands))
		i := 0
		for k, _ := range plugins.Commands {
			mk[i] = k
			i++
		}
		sort.Strings(mk)
		var buffer bytes.Buffer

		for _, v := range mk {
			buffer.WriteString("/" + v + " - " + plugins.Commands[v] + "\n")
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"
		msg.Text = "Those are my commands: \n"+buffer.String()
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = update.Message.MessageID

		_, err11 := plugins.Bot.Send(msg)
		if err11 != nil {
			dlog.Errorln(err11)
		}
	}
}