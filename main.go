package main

import (
	"log"
	// "strconv"

	config "github.com/ad/corpobot/config"
	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	telegram "github.com/ad/corpobot/telegram"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	_ "github.com/ad/corpobot/plugins/help"
	_ "github.com/ad/corpobot/plugins/start"
	_ "github.com/ad/corpobot/plugins/echo"
)

const version = "0.0.1"

var (
	err error

	bot    *tgbotapi.BotAPI
	db     *sql.DB
)

func main() {
	dlog.Infof("Started version %s", version)

	// Init Config
	config := config.InitConfig()

	log.SetFlags(0)

	// Init DB
	db, err = database.InitDB()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
		return
	}
	defer func() { _ = db.Close() }()

	// Init Telegram
	bot, err = telegram.InitTelegram(config.TelegramToken, config.TelegramProxyHost, config.TelegramProxyPort, config.TelegramProxyUser, config.TelegramProxyPassword, config.TelegramDebug)
	if err != nil {
		log.Fatalf("fail on telegram login: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("[INIT] [Failed to init Telegram updates chan: %v]", err)
	}

	// Bootstrapper for plugins
	for _, d := range plugins.Plugins {
		go d.OnStart()
	}
	// log.Println(strconv.Itoa(len(plugins.Plugins)) + " plugins loaded")

	telegram.ProcessTelegramMessages(db, bot, updates)
}
