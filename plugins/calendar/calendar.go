package calendar

import (
	"strings"
	"time"

	cal "github.com/ad/corpobot/calendar"
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
	if !plugins.CheckIfPluginDisabled("calendar.Plugin", "enabled") {
		return
	}

	plugins.RegisterCommand("calendar", "Show calendar", []string{"member", "admin", "owner"})
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[calendar.Plugin] Stopped")

	plugins.UnregisterCommand("calendar")
}

func (m *Plugin) Run(update *tgbotapi.Update, command, args string, user *database.User) (bool, error) {
	if plugins.CheckIfCommandIsAllowed(command, "calendar", user.Role) {
		replyKeyboard := tgbotapi.InlineKeyboardMarkup{}

		switch {
		case strings.HasPrefix(args, "<"):
			date := strings.TrimLeft(args, "<")
			year, month, _, err := cal.ParseDate(date)
			if err == nil {
				replyKeyboard, _, _ = cal.HandlerPrevMonth("/calendar", year, time.Month(month))
			}
		case strings.HasPrefix(args, ">"):
			date := strings.TrimLeft(args, ">")
			year, month, _, err := cal.ParseDate(date)
			if err == nil {
				replyKeyboard, _, _ = cal.HandlerNextMonth("/calendar", year, time.Month(month))
			}
		case strings.HasPrefix(args, "«"):
			date := strings.TrimLeft(args, "«")
			year, month, _, err := cal.ParseDate(date)
			if err == nil {
				replyKeyboard, _, _ = cal.HandlerPrevYear("/calendar", year, time.Month(month))
			}
		case strings.HasPrefix(args, "»"):
			date := strings.TrimLeft(args, "»")
			year, month, _, err := cal.ParseDate(date)
			if err == nil {
				replyKeyboard, _, _ = cal.HandlerNextYear("/calendar", year, time.Month(month))
			}
		case strings.HasPrefix(args, "m"):
			currentTime := time.Now()
			year := currentTime.Year()
			month := currentTime.Month()

			date := strings.TrimLeft(args, "m")
			year2, month2, _, err := cal.ParseDate(date)
			if err == nil {
				year = year2
				month = time.Month(month2)
			}
			replyKeyboard = cal.GenerateMonths("/calendar", year, month)
		case strings.HasPrefix(args, "y"):
			currentTime := time.Now()
			year := currentTime.Year()
			month := currentTime.Month()

			date := strings.TrimLeft(args, "y")
			year2, month2, _, err := cal.ParseDate(date)
			if err == nil {
				year = year2
				month = time.Month(month2)
			}
			replyKeyboard = cal.GenerateYears("/calendar", year, month)
		default:
			currentTime := time.Now()
			year := currentTime.Year()
			month := currentTime.Month()

			year2, month2, _, err := cal.ParseDate(args)
			if err == nil {
				year = year2
				month = time.Month(month2)
			}
			replyKeyboard = cal.GenerateCalendar("/calendar", year, month)
		}

		if update.CallbackQuery != nil {
			_, err := plugins.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			if err != nil {
				dlog.Errorln(err.Error())
			}

			edit := tgbotapi.EditMessageReplyMarkupConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      update.CallbackQuery.Message.Chat.ID,
					MessageID:   update.CallbackQuery.Message.MessageID,
					ReplyMarkup: &replyKeyboard,
				},
			}

			_, err = plugins.Bot.Send(edit)
			return true, err
		}

		return true, telegram.SendCustom(user.TelegramID, 0, "Calendar", false, &replyKeyboard)
	}

	return false, nil
}
