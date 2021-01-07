package starthelp

import (
	"bytes"
	"sort"

	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	"github.com/ad/corpobot/telegram"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Plugin struct{}

func init() {
	plugins.RegisterPlugin(&Plugin{})
}

func (m *Plugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("starthelp.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("start", "Bot /start command", []string{"new", "member", "admin", "owner"}, start)
	plugins.RegisterCommand("help", "Display this help", []string{"new", "member", "admin", "owner"}, help)
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[starthelp.Plugin] Stopped")

	plugins.UnregisterCommand("start")
	plugins.UnregisterCommand("help")
}

var start plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	return telegram.Send(user.TelegramID, "Hello! Send /help")
}

var help plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	mk := make(map[string]string)
	var keys []string

	plugins.Commands.Range(func(k, v interface{}) bool {
		if plugins.CheckIfCommandIsAllowed(k.(string), user.Role) {
			mk[k.(string)] = v.(plugins.Command).Description
			keys = append(keys, k.(string))
		}
		return true
	})

	sort.Strings(keys)
	var buffer bytes.Buffer

	for _, k := range keys {
		_, err := buffer.WriteString("/" + k + " - " + mk[k] + "\n")
		if err != nil {
			return err
		}
	}

	return telegram.Send(user.TelegramID, "Those are my commands: \n"+buffer.String())
}
