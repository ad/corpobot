package plugins

import (
	"reflect"
	"strings"

	database "github.com/ad/corpobot/db"
	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TelegramPlugin ...
type TelegramPlugin interface {
	Run(update *tgbotapi.Update) (bool, error)
	OnStart()
	OnStop()
}

// These are are registered plugins
var (
	Plugins         = map[string]TelegramPlugin{}
	DisabledPlugins = map[string]TelegramPlugin{}
	Commands        = make(map[string]string)
	Bot             *tgbotapi.BotAPI
	DB              *sql.DB
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
	} else {
		dlog.Debugln("Plugin '" + plugin + "' does not exist or is not loaded")
	}
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
		PluginInstance.OnStart()

		dlog.Debugln(plugin + " enabled ")
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
		dlog.Debugln("[" + name + "] Disabled")
		return false
	}

	dlog.Debugln("[" + name + "] Started")

	return true
}

// KeyOf ...
func KeyOf(p TelegramPlugin) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}

// Register a Command exported by a plugin
func RegisterCommand(command string, description string) {
	Commands[command] = description
}

// UnRegister a Command exported by a plugin
func UnregisterCommand(command string) {
	delete(Commands, command)
}
