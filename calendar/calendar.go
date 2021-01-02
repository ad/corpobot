package calendar

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// // getMonthNames ...
// func getMonthNames(lang string) [12]string {
// 	switch lang {
// 	case "ru":
// 		return [12]string{"янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"}
// 	default:
// 		return [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
// 	}
// }

// getWeekdayNames ...
func getWeekdayNames(lang string) [7]string {
	switch lang {
	case "ru":
		return [7]string{"П", "В", "С", "Ч", "П", "С", "В"}
	default:
		return [7]string{"M", "T", "W", "T", "F", "S", "S"}
	}
}

func GenerateCalendar(year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	keyboard = addMonthYearRow(year, month, keyboard)
	keyboard = addDaysNamesRow(keyboard)
	keyboard = generateMonth(year, int(month), keyboard)
	keyboard = addSpecialButtons(year, month, keyboard)
	return keyboard
}

func HandlerPrevButton(year int, month time.Month) (tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	if month != time.January {
		month--
	} else {
		month = 12
		year--
	}
	return GenerateCalendar(year, month), year, month
}

func HandlerNextButton(year int, month time.Month) (tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	if month != time.December {
		month++
	} else {
		month = 1
		year++
	}
	return GenerateCalendar(year, month), year, month
}

func addMonthYearRow(year int, month time.Month, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %v", month, year), "1")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}

func addDaysNamesRow(keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	lang := "ru"
	days := getWeekdayNames(lang)
	var rowDays []tgbotapi.InlineKeyboardButton
	for _, day := range days {
		btn := tgbotapi.NewInlineKeyboardButtonData(day, day)
		rowDays = append(rowDays, btn)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}

func generateMonth(year int, month int, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	firstDay := date(year, month, 0)
	amountDaysInMonth := date(year, month+1, 0).Day()

	weekday := int(firstDay.Weekday())
	rowDays := []tgbotapi.InlineKeyboardButton{}
	for i := 1; i <= weekday; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(" ", strconv.Itoa(i))
		rowDays = append(rowDays, btn)
	}

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
			amountWeek = 0
			rowDays = []tgbotapi.InlineKeyboardButton{}
		}

		day := strconv.Itoa(i)
		if len(day) == 1 {
			day = fmt.Sprintf("0%v", day)
		}
		monthStr := strconv.Itoa(month)
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%v", monthStr)
		}

		btnText := fmt.Sprintf("%v", i)
		if time.Now().Day() == i && time.Now().Month() == time.Month(month) && time.Now().Year() == year {
			btnText = fmt.Sprintf("[%v]", i)
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("%v.%v.%v", year, monthStr, day))
		rowDays = append(rowDays, btn)
		amountWeek++
	}
	for i := 1; i <= 7-amountWeek; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(" ", strconv.Itoa(i))
		rowDays = append(rowDays, btn)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)

	return keyboard
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func addSpecialButtons(year int, month time.Month, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	var rowDays = []tgbotapi.InlineKeyboardButton{}
	btnPrev := tgbotapi.NewInlineKeyboardButtonData("<", "/calendar < "+fmt.Sprintf("%v.%d", year, month))
	btnNext := tgbotapi.NewInlineKeyboardButtonData(">", "/calendar > "+fmt.Sprintf("%v.%d", year, month))
	rowDays = append(rowDays, btnPrev, btnNext)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}
