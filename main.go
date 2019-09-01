package main

import (
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	ping "github.com/digineo/go-ping"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	addr := os.Getenv("POSIE_ADDR")
	token := os.Getenv("POSIE_TG_TOKEN")
	chat := os.Getenv("POSIE_TG_CHAT")
	textOk := os.Getenv("POSIE_TEXT_OK")
	textFail := os.Getenv("POSIE_TEXT_FAIL")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	st := make(chan os.Signal, 1)
	signal.Notify(st, os.Interrupt)

	tick := time.NewTicker(time.Second*10 + time.Millisecond*time.Duration(rand.Intn(10)))

	chatID, err := strconv.Atoi(chat)
	if err != nil {
		panic(err)
	}

	pinger, err := ping.New("0.0.0.0", "")
	if err != nil {
		panic(err)
	}

	dest, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		panic(err)
	}

	var unaccCnt, accCnt int
	var changed bool

ST:
	for {
		select {
		case <-st:
			tick.Stop()
			break ST
		case <-tick.C:
			_, err := pinger.Ping(dest, time.Second*10)

			if err != nil {
				unaccCnt++
				accCnt = 0
				changed = true
			} else {
				accCnt++
				unaccCnt = 0
			}

			var txt string

			if accCnt == 2 && changed {
				txt = textOk
				println("Ok detected")
			} else if unaccCnt == 2 {
				txt = textFail
				println("Fail detected")
			}

			if txt != "" {
				msg := tgbotapi.NewMessage(int64(chatID), txt)
				_, err := bot.Send(msg)
				if err != nil {
					println(err.Error())
				}
			}
		}
	}
}
