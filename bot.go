package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	logfile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logfile)
	} else {
		log.Error("Failed to open log file for writing")
	}
}

func generateRandomNumber() int {
	return rand.Intn(900) + 100
}

func main() {
	viper.SetConfigFile("config.ini")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	botToken := viper.GetString("Telegram.token")
	text := viper.GetString("Telegram.text")
	chatIDsStr := viper.GetString("Telegram.chat_id")

	chatIDs := parseChatIDs(chatIDsStr)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	num := generateRandomNumber()

	for _, chatID := range chatIDs {

		msg := tgbotapi.NewMessage(chatID, text)
		_, err := bot.Send(msg)
		if err != nil {
			errMsg := "Error sending message to chat ID " + strconv.FormatInt(chatID, 10) + ": " + err.Error()
			log.Error(errMsg)
			continue
		}

		numMsg := tgbotapi.NewMessage(chatID, strconv.Itoa(num))
		_, err = bot.Send(numMsg)
		if err != nil {
			errMsg := "Error sending number message to chat ID " + strconv.FormatInt(chatID, 10) + ": " + err.Error()
			log.Error(errMsg)
			continue
		}
	}

	log.Info("Finished sending messages")

	os.Exit(0)
}

func parseChatIDs(chatIDsStr string) []int64 {
	chatIDsStr = strings.Trim(chatIDsStr, "[]")
	ids := strings.Split(chatIDsStr, " ")
	var result []int64
	for _, idStr := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err == nil {
			result = append(result, id)
		}
	}
	return result
}
