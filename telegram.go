package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var authorizedUsers map[int]bool

func initTelegram() {
	authorizedUsers = make(map[int]bool)
	for _, userID := range strings.Split(os.Getenv("TG_AUTHORIZED_USERS"), ",") {
		id, err := strconv.Atoi(userID)
		if err == nil {
			authorizedUsers[id] = true
		}
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("[x] Telegram: authorized on account [%s], awaiting for updates...", bot.Self.UserName)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		if update.Message != nil {
			processUpdates(bot, &update)
		}
	}
}

func processUpdates(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !authorizedUsers[update.Message.From.ID] {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied.")
		_, _ = bot.Send(msg)
		return
	}
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "update":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Starting sequence...")
			sentMsg, _ := bot.Send(msg)
			editMsg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sentMsg.MessageID, "Processing, please wait...")
			_, _ = bot.Send(editMsg)
			data := parse(os.Getenv("BLUECOSMO_USERNAME"), os.Getenv("BLUECOSMO_PASSWORD"))
			if len(data) > 0 {
				log.Println("[x] Telegram: sending data...")
				pattern := `\d+`
				regex := regexp.MustCompile(pattern)
				totalMinutes, _ := strconv.Atoi(regex.FindString(data[1]))
				minutesUsed, _ := strconv.Atoi(regex.FindString(data[2]))
				leftMinutes := strconv.Itoa(totalMinutes - minutesUsed)
				newmsg := fmt.Sprintf(`<b>Account number</b>: %s
<b>Service number</b>: %s
<b>Plan name</b>: %s
<b>Minutes used:</b> %s
<b>Status:</b> %s

<i>Based on your plan name, it looks like you have <b>%s</b> minute(s). You spent <b>%s</b> minute(s), so you should have <b>%s</b> minute(s) left.</i>`,
					data[0], data[1], data[2], data[3], data[4], strconv.Itoa(totalMinutes), strconv.Itoa(minutesUsed), leftMinutes)
				editMsg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sentMsg.MessageID, newmsg)
				editMsg.ParseMode = tgbotapi.ModeHTML
				_, _ = bot.Send(editMsg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error while fetching data from BlueCosmo. Try again later.")
				_, _ = bot.Send(msg)
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command.")
			_, _ = bot.Send(msg)
		}
	}
}
