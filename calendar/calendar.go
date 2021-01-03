package calendar

import (
	"fmt"
	"strconv"
	"strings"
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

func GenerateCalendar(command string, year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	keyboard = addMonthYearRow(command, year, month, keyboard)
	keyboard = addDaysNamesRow(keyboard)
	keyboard = generateMonth(command, year, int(month), keyboard)
	// keyboard = addSpecialButtons(command, year, month, keyboard)
	return keyboard
}

func HandlerPrevMonth(command string, year int, month time.Month) (tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	if month != time.January {
		month--
	} else {
		month = 12
		year--
	}
	return GenerateCalendar(command, year, month), year, month
}

func HandlerNextMonth(command string, year int, month time.Month) (tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	if month != time.December {
		month++
	} else {
		month = 1
		year++
	}
	return GenerateCalendar(command, year, month), year, month
}

func HandlerPrevYear(command string, year int, month time.Month) (tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	year--
	return GenerateCalendar(command, year, month), year, month
}

func HandlerNextYear(command string, year int, month time.Month) (tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	year++
	return GenerateCalendar(command, year, month), year, month
}

func ParseDate(date string) (int, int, int, error) {
	if date != "" {
		dateArray := strings.SplitN(date, ".", 3)
		if len(dateArray) >= 2 {
			year, err1 := strconv.Atoi(dateArray[0])
			if err1 != nil {
				return 0, 0, 0, err1
			}
			month, err2 := strconv.Atoi(dateArray[1])
			if err2 != nil {
				return 0, 0, 0, err2
			}
			day := 0
			if len(dateArray) == 3 {
				var err3 error
				day, err3 = strconv.Atoi(dateArray[2])
				if err3 != nil {
					return 0, 0, 0, err3
				}
			}
			if year > 0 && month > 0 && month < 13 {
				return year, month, day, nil
			}
		}
	}
	return 0, 0, 0, fmt.Errorf("%s", "wrong date format")
}

func addMonthYearRow(command string, year int, month time.Month, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton
	btnPrevYear := tgbotapi.NewInlineKeyboardButtonData("«", command+" « "+fmt.Sprintf("%v.%d", year, month))
	btnPrevMonth := tgbotapi.NewInlineKeyboardButtonData("<", command+" < "+fmt.Sprintf("%v.%d", year, month))
	btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %v", month, year), command)
	btnNextMonth := tgbotapi.NewInlineKeyboardButtonData(">", command+" > "+fmt.Sprintf("%v.%d", year, month))
	btnNextYear := tgbotapi.NewInlineKeyboardButtonData("»", command+" » "+fmt.Sprintf("%v.%d", year, month))
	row = append(row, btnPrevYear, btnPrevMonth, btn, btnNextMonth, btnNextYear)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}

func addDaysNamesRow(keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	lang := "ru"
	days := getWeekdayNames(lang)
	var rowDays []tgbotapi.InlineKeyboardButton
	for _, day := range days {
		btn := tgbotapi.NewInlineKeyboardButtonData(day, " ")
		rowDays = append(rowDays, btn)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}

func generateMonth(command string, year int, month int, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	firstDay := date(year, month, 0)
	amountDaysInMonth := date(year, month+1, 0).Day()

	weekday := int(firstDay.Weekday())
	rowDays := []tgbotapi.InlineKeyboardButton{}
	for i := 1; i <= weekday; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(" ", " ")
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
		btn := tgbotapi.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("%s %v.%v.%v", command, year, monthStr, day))
		rowDays = append(rowDays, btn)
		amountWeek++
	}
	for i := 1; i <= 7-amountWeek; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(" ", " ")
		rowDays = append(rowDays, btn)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)

	return keyboard
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// func addSpecialButtons(command string, year int, month time.Month, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
// 	var rowDays = []tgbotapi.InlineKeyboardButton{}
// 	btnPrev := tgbotapi.NewInlineKeyboardButtonData("<", command+" < "+fmt.Sprintf("%v.%d", year, month))
// 	btnNext := tgbotapi.NewInlineKeyboardButtonData(">", command+" > "+fmt.Sprintf("%v.%d", year, month))
// 	rowDays = append(rowDays, btnPrev, btnNext)
// 	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
// 	return keyboard
// }
