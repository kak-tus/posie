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

const switchLimit = time.Minute * 1

const (
	_ = iota
	stateOK
	stateFail
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

	lastFail := time.Now()
	lastOk := time.Now()
	state := stateOK

ST:
	for {
		select {
		case <-st:
			tick.Stop()
			break ST
		case <-tick.C:
			_, err := pinger.Ping(dest, time.Second*10)

			if err != nil {
				lastFail = time.Now()
			} else {
				lastOk = time.Now()
			}

			var txt string

			if lastOk.After(lastFail) && lastOk.Sub(lastFail) > switchLimit && state == stateFail {
				txt = textOk
				state = stateOK
				println("Ok detected")
			} else if lastFail.After(lastOk) && lastFail.Sub(lastOk) > switchLimit && state == stateOK {
				txt = textFail
				state = stateFail
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
