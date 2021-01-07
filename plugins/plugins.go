package plugins

import (
	"reflect"
	"strings"
	"sync"

	"github.com/ad/corpobot/config"
	database "github.com/ad/corpobot/db"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TelegramPlugin ...
type TelegramPlugin interface {
	OnStart()
	OnStop()
}

// Command ...
type Command struct {
	Description string          `sql:"description"`
	Roles       map[string]bool `sql:"roles"`
	Callback    CommandCallback
}

type CommandCallback func(update *tgbotapi.Update, command, args string, user *database.User) (bool, error)

var (
	Plugins         sync.Map
	DisabledPlugins sync.Map
	Commands        sync.Map
	Bot             *tgbotapi.BotAPI
	DB              *sql.DB
	Config          *config.Config
)

// Register a Plugin
func RegisterPlugin(p TelegramPlugin) {
	Plugins.Store(KeyOf(p), p)
}

// Disable a plugin
func DisablePlugin(plugin string) bool {
	plugin = strings.TrimSpace(plugin)

	v, exists := Plugins.Load(plugin)
	if exists {
		DisabledPlugins.Store(plugin, v.(TelegramPlugin))
		_, disabled := DisabledPlugins.Load(plugin)
		if disabled {
			Plugins.Delete(plugin)
			v.(TelegramPlugin).OnStop()

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

	_, PluginExists := Plugins.Load(plugin)
	if PluginExists {
		return true
	}

	PluginInstance, InstanceExists := DisabledPlugins.Load(plugin)
	Plugins.Store(plugin, PluginInstance)
	if InstanceExists {
		DisabledPlugins.Delete(plugin)
		dlog.Debugln("[" + plugin + "] enabled ")

		PluginInstance.(TelegramPlugin).OnStart()
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
		return false
	}

	if !plugin.IsEnabled() {
		if v, ok := Plugins.Load(name); ok {
			DisabledPlugins.Store(name, v.(TelegramPlugin))
			Plugins.Delete(name)
			dlog.Debugln("[" + name + "] Disabled")
			return false
		}
	}

	dlog.Debugln("[" + name + "] Started")

	return true
}

func CheckIfCommandIsAllowed(command, role string) bool {
	if v, ok := Commands.Load(command); ok {
		roles := v.(Command).Roles
		if _, ok2 := roles[role]; ok2 {
			return true
		}
	}

	return false
}

// KeyOf ...
func KeyOf(p TelegramPlugin) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}

// Register a Command exported by a plugin
func RegisterCommand(command string, description string, roles []string, callback CommandCallback) {
	r := make(map[string]bool)
	for _, v := range roles {
		r[v] = true
	}
	Commands.Store(command, Command{Description: description, Roles: r, Callback: callback})
}

// UnRegister a Command exported by a plugin
func UnregisterCommand(command string) {
	Commands.Delete(command)
}

func (cmd Command) IsAllowedForRole(role string) bool {
	_, ok := cmd.Roles[role]
	return ok
}
