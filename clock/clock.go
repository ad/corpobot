package clock

import (
	"fmt"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func GenerateClock(command string, hour, minute int, lang string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	if hour < 0 {
		keyboard = addHours(command, keyboard, lang)
	} else {
		keyboard = addMinutes(command, hour, keyboard, lang)
	}
	return keyboard
}

func addHours(command string, keyboard tgbotapi.InlineKeyboardMarkup, lang string) tgbotapi.InlineKeyboardMarkup {
	var rowHours []tgbotapi.InlineKeyboardButton
	pos := 0
	for i := 0; i < 24; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%02d:...", i), fmt.Sprintf("%s %02d", command, i))
		rowHours = append(rowHours, btn)
		if (pos+1)%6 == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowHours)
			rowHours = []tgbotapi.InlineKeyboardButton{}
		}
		pos++
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowHours)
	return keyboard
}

func addMinutes(command string, hour int, keyboard tgbotapi.InlineKeyboardMarkup, lang string) tgbotapi.InlineKeyboardMarkup {
	var rowHours []tgbotapi.InlineKeyboardButton
	pos := 0
	for i := 0; i < 60; i += 5 {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%02d:%02d", hour, i), fmt.Sprintf("%s %02d:%02d", command, hour, i))
		rowHours = append(rowHours, btn)
		if (pos+1)%4 == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowHours)
			rowHours = []tgbotapi.InlineKeyboardButton{}
		}
		pos++
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowHours)
	return keyboard
}
