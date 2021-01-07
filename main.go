package main

import (
	"log"
	"sync"
	"time"

	config "github.com/ad/corpobot/config"
	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	_ "github.com/ad/corpobot/plugins/admin"
	_ "github.com/ad/corpobot/plugins/echo"
	_ "github.com/ad/corpobot/plugins/groupchats"
	_ "github.com/ad/corpobot/plugins/groups"
	_ "github.com/ad/corpobot/plugins/me"
	_ "github.com/ad/corpobot/plugins/messages"
	_ "github.com/ad/corpobot/plugins/starthelp"
	_ "github.com/ad/corpobot/plugins/users"
	telegram "github.com/ad/corpobot/telegram"
	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const version = "0.2.8"

var (
	err error

	bot *tgbotapi.BotAPI
	db  *sql.DB
)

func main() {
	dlog.Infof("Started version %s", version)

	// Init Config
	config := config.InitConfig()
	plugins.Config = config

	log.SetFlags(0)

	// Init DB
	db, err = database.InitDB()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
		return
	}
	defer func() { _ = db.Close() }()
	plugins.DB = db

	// Init Telegram
	bot, err = telegram.InitTelegram(config.TelegramToken, config.TelegramProxyHost, config.TelegramProxyPort, config.TelegramProxyUser, config.TelegramProxyPassword, config.TelegramDebug)
	if err != nil {
		log.Printf("fail on telegram login: %v", err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("[INIT] [Failed to init Telegram updates chan: %v]", err)
		return
	}

	dlog.Debugln("Waiting for plugins...")
	for {
		if plugins.DB != nil && plugins.DB.Ping() == nil {
			var wg sync.WaitGroup

			// Bootstrapper for plugins
			plugins.Plugins.Range(func(k, v interface{}) bool {
				wg.Add(1)
				go func() {
					defer wg.Done()
					v.(plugins.TelegramPlugin).OnStart()
				}()
				return true
			})
			wg.Wait()
			break
		}
		time.Sleep(time.Second)
	}

	telegram.ProcessTelegramMessages(db, bot, updates)
}
