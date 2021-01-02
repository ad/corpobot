package starthelp

import (
	"bytes"
	"sort"

	database "github.com/ad/corpobot/db"
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
	if !plugins.CheckIfPluginDisabled("starthelp.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("start", "Bot /start command", []string{"new", "member", "admin", "owner"})
	plugins.RegisterCommand("help", "Display this help", []string{"new", "member", "admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[starthelp.Plugin] Stopped")

	plugins.UnregisterCommand("start")
	plugins.UnregisterCommand("help")
}

func (m *Plugin) Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "start", user.Role) {
		return true, telegram.Send(user.TelegramID, "Hello! Send /help")
	}

	if plugins.CheckIfCommandIsAllowed(command, "help", user.Role) {
		var mk []string

		for k := range plugins.Commands {
			mk = append(mk, k)
		}

		sort.Strings(mk)
		var buffer bytes.Buffer

		for _, v := range mk {
			if plugins.CheckIfCommandIsAllowed(v, v, user.Role) {
				_, err := buffer.WriteString("/" + v + " - " + plugins.Commands[v].Description + "\n")
				if err != nil {
					return true, err
				}
			}
		}

		return true, telegram.Send(user.TelegramID, "Those are my commands: \n"+buffer.String())
	}

	return false, nil
}
