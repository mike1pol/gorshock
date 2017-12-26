package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const Token = "470061611:AAG_TOk_1lIRTwOLvlff6ZcFRKnf3MPobIA"
const MaxMessageLength = 4096
const MessageWasTruncatedText = "\n\n(message was truncated)"

var CommandRe = regexp.MustCompile("^/([a-z]+)(@gorshock_bot)? *(.*)$")
var VolumeRe = regexp.MustCompile("\\d+%")

func execCmd(arg string, args ...string) string {
	cmd := exec.Command(arg, args...)

	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprint(err)
	}

	return string(stdoutStderr)
}

func handleCommand(msg string) string {
	matches := CommandRe.FindStringSubmatch(msg)

	if len(matches) == 0 {
		return "Не понял команды"
	} else {
		var cmd = matches[1]
		var rest = strings.Trim(matches[len(matches)-1], " \t\n")

		if cmd == "volume" {
			vols := strings.Split(rest, " ")

			if len(vols) == 1 && len(vols[0]) > 0 {
				vol, _ := strconv.ParseInt(vols[0], 10, 32)

				if vol < 0 {
					vol = 0
				} else if vol > 100 {
					vol = 100
				}

				execCmd("amixer", "set", "PCM", fmt.Sprintf("%d%%", vol))

				return fmt.Sprintf("Volume set to %d%%", vol)
			} else {
				vol := execCmd("amixer", "get", "PCM")

				matches := VolumeRe.FindStringSubmatch(vol)

				if matches == nil {
					return "Unknown error happened"
				} else {
					return matches[0]
				}
			}
		} else if cmd == "mpc" {
			args := strings.Split(rest, " ")
			return execCmd("mpc", args...)
		} else if cmd == "ip" {
			return execCmd("/bin/bash", "-c", "/sbin/ifconfig wlan0 | awk '/inet /{print substr($2,1)}'")
		}
	}

	return "ok"
}

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		response := handleCommand(msg.Text)

		if len(response) > MaxMessageLength {
			response = response[0:MaxMessageLength-len(MessageWasTruncatedText)] + MessageWasTruncatedText
		}

		msg.Text = response
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
