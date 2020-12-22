package echo

import (
	// "strconv"

	"github.com/ad/corpobot/plugins"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type EchoPlugin struct {
}

func init() {
	plugins.RegisterPlugin(&EchoPlugin{})
}
func (m *EchoPlugin) OnStart() {
	dlog.Debugln("[EchoPlugin] Started")
	plugins.RegisterCommand("echo", "...")
}
func (m *EchoPlugin) OnStop() {
	dlog.Debugln("[EchoPlugin] Stopped")
	plugins.UnregisterCommand("echo")
}

func (m *EchoPlugin) Run(update *tgbotapi.Update) {
	if update.Message.Command() == "echo" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"
		msg.Text = update.Message.Text
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = update.Message.MessageID

		_, err11 := plugins.Bot.Send(msg)
		if err11 != nil {
			dlog.Errorln(err11)
		}
	}
}