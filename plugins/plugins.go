package plugins

import (
	"reflect"
	"strings"

	config "github.com/ad/corpobot/config"
	database "github.com/ad/corpobot/db"
	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TelegramPlugin ...
type TelegramPlugin interface {
	Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error)
	OnStart()
	OnStop()
}

// Command ...
type Command struct {
	Description string          `sql:"description"`
	Roles       map[string]bool `sql:"roles"`
}

var (
	Plugins         = map[string]TelegramPlugin{}
	DisabledPlugins = map[string]TelegramPlugin{}
	Commands        = make(map[string]Command)
	Bot             *tgbotapi.BotAPI
	DB              *sql.DB
	Config          *config.Config
)

// Register a Plugin
func RegisterPlugin(p TelegramPlugin) {
	Plugins[KeyOf(p)] = p
}

// Disable a plugin
func DisablePlugin(plugin string) bool {
	plugin = strings.TrimSpace(plugin)

	_, exists := Plugins[plugin]
	if exists {
		DisabledPlugins[plugin] = Plugins[plugin]
		_, disabled := DisabledPlugins[plugin]
		if disabled {
			delete(Plugins, plugin)
			DisabledPlugins[plugin].OnStop()

			dlog.Debugln(plugin + " removed from running plugins")
		} else {
			dlog.Debugln("Can't disable " + plugin + ", odd")
		}
		return disabled
	}

	dlog.Debugln("Plugin '" + plugin + "' does not exist or is not loaded")

	return exists
}

// Enable a plugin
func EnablePlugin(plugin string) bool {
	plugin = strings.TrimSpace(plugin)

	_, PluginExists := Plugins[plugin]
	if PluginExists {
		return true
	}

	PluginInstance, InstanceExists := DisabledPlugins[plugin]
	Plugins[plugin] = PluginInstance
	if InstanceExists {
		delete(DisabledPlugins, plugin)
		dlog.Debugln("[" + plugin + "] enabled ")

		PluginInstance.OnStart()
		return true
	}
	return false
}

func CheckIfPluginDisabled(name, state string) bool {
	plugin := &database.Plugin{
		Name:  name,
		State: state,
	}

	plugin, err := database.AddPluginIfNotExist(DB, plugin)
	if err != nil {
		dlog.Errorln("failed: " + err.Error())
	}

	if plugin.State != "enabled" {
		DisabledPlugins[name] = Plugins[name]
		delete(Plugins, name)
		dlog.Debugln("[" + name + "] Disabled")
		return false
	}

	dlog.Debugln("[" + name + "] Started")

	return true
}

func CheckIfCommandIsAllowed(command, command2, role string) bool {
	if command == command2 {
		if _, ok := Commands[command]; ok {
			roles := Commands[command].Roles
			if _, ok2 := roles[role]; ok2 {
				return true
			}
		}
	}

	return false
}

// KeyOf ...
func KeyOf(p TelegramPlugin) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}

// Register a Command exported by a plugin
func RegisterCommand(command string, description string, roles []string) {
	r := make(map[string]bool)
	for _, v := range roles {
		r[v] = true
	}

	Commands[command] = Command{Description: description, Roles: r}
}

// UnRegister a Command exported by a plugin
func UnregisterCommand(command string) {
	delete(Commands, command)
}
