package start

import (
	// "strconv"

	"github.com/ad/corpobot/plugins"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type StartPlugin struct {
}

func init() {
	plugins.RegisterPlugin(&StartPlugin{})
}
func (m *StartPlugin) OnStart() {
	dlog.Debugln("[StartPlugin] Started")
	plugins.RegisterCommand("start", "...")
}
func (m *StartPlugin) OnStop() {
	dlog.Debugln("[StartPlugin] Stopped")
	plugins.UnregisterCommand("start")
}

func (m *StartPlugin) Run(update *tgbotapi.Update) {
	if update.Message.Command() == "start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"
		msg.Text = "Hello!"
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = update.Message.MessageID

		_, err11 := plugins.Bot.Send(msg)
		if err11 != nil {
			dlog.Errorln(err11)
		}
	}
}