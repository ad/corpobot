package calendar

import (
	"fmt"
	"strconv"
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
			year, month, err := parseDate(date)
			if err == nil {
				replyKeyboard, _, _ = cal.HandlerPrevButton("/calendar", year, time.Month(month))
			}
		case strings.HasPrefix(args, ">"):
			date := strings.TrimLeft(args, ">")
			year, month, err := parseDate(date)
			if err == nil {
				replyKeyboard, _, _ = cal.HandlerNextButton("/calendar", year, time.Month(month))
			}
		default:
			currentTime := time.Now()
			year := currentTime.Year()
			month := currentTime.Month()

			year2, month2, err := parseDate(args)
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

func parseDate(date string) (int, int, error) {
	if date != "" {
		dateArray := strings.SplitN(date, ".", 2)
		if len(dateArray) == 2 {
			year, err1 := strconv.Atoi(dateArray[0])
			if err1 != nil {
				return 0, 0, err1
			}
			month, err2 := strconv.Atoi(dateArray[1])
			if err2 != nil {
				return 0, 0, err2
			}
			if year > 0 && month > 0 && month < 13 {
				return year, month, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("%s", "wrong date format")
}
