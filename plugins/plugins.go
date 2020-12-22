package plugins

// type Plugin interface {
// 	GetCommand() string
// 	RunCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool
// }

import (
	"reflect"
	"strings"

	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type TelegramPlugin interface {
	Run(update *tgbotapi.Update)
	OnStart()
	OnStop()
}

// These are are registered plugins
var Plugins = map[string]TelegramPlugin{}
var DisabledPlugins = map[string]TelegramPlugin{}
var Commands = make(map[string]string)
var Bot *tgbotapi.BotAPI

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