// db
package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func (bot *Bot) getQuote() string {
	length := len(bot.quotes)
	if length == 0 {
		return "No quotes stored!"
	}
	randomed := rand.Intn(length)
	if randomed == 0 {
		randomed = 1
	}
	tempInt := 1
	for quote, _ := range bot.quotes {
		if randomed == tempInt {
			return quote
		}
		tempInt++
	}
	return "Error!"
}

func (bot *Bot) writeQuoteDB() {
	dst, err := os.Create("quotes" + bot.channel + ".ini")
	defer dst.Close()
	if err != nil {
		fmt.Println("Can't write to QuoteDB from " + bot.channel)
		return
	}
	for split1, split2 := range bot.quotes {
		fmt.Fprintf(dst, split1+"|"+split2+"\n")
	}
}

func (bot *Bot) readQuoteDB() {
	quotes, err := ioutil.ReadFile("quotes" + bot.channel + ".ini")
	if err != nil {
		fmt.Println("Unable to read QuoteDB from " + bot.channel)
		return
	}
	split1 := strings.Split(string(quotes), "\n")
	for _, splitted1 := range split1 {
		split2 := strings.Split(splitted1, "|")
		bot.quotes[split2[0]] = split2[1]
	}

}

func (bot *Bot) readSettingsDB(channel string) bool {
	settings, err := ioutil.ReadFile("settings#" + channel + ".ini")
	bot.channel = "#" + channel
	if err != nil {
		fmt.Println("Unable to read SettingsDB from " + channel)
		return false
	}
	split1 := strings.Split(string(settings), "\n")
	for _, splitted1 := range split1 {
		split2 := strings.Split(splitted1, "|")
		if split2[0] == "nickname" {
			bot.nick = split2[1]
		} else if split2[0] == "timemsg" {
			bot.autoMSG1 = split2[1]
		} else if split2[0] == "linemsg" {
			bot.autoMSG2 = split2[1]
		} else if split2[0] == "timemsgminutes" {
			temp, _ := strconv.Atoi(split2[1])
			bot.autoMSG1Count = temp
		} else if split2[0] == "linemsgcount" {
			temp, _ := strconv.Atoi(split2[1])
			bot.autoMSG2Count = temp
		} else if split2[0] == "userspamcount" {
			temp, _ := strconv.Atoi(split2[1])
			bot.userMaxLastMsg = temp
		}

	}
	return true
}

func (bot *Bot) writeSettingsDB() {
	dst, err := os.Create("settings" + bot.channel + ".ini")
	defer dst.Close()
	if err != nil {
		fmt.Println("Can't write to SettingsDB from " + bot.channel)
		return
	}
	fmt.Fprintf(dst, "nickname|"+bot.nick+"\n")
	fmt.Fprintf(dst, "timemsg|"+bot.autoMSG1+"\n")
	fmt.Fprintf(dst, "linemsg|"+bot.autoMSG2+"\n")
	fmt.Fprintf(dst, "timemsgminutes|"+strconv.Itoa(bot.autoMSG1Count)+"\n")
	fmt.Fprintf(dst, "linemsgcount|"+strconv.Itoa(bot.autoMSG2Count)+"\n")
	fmt.Fprintf(dst, "userspamcount|"+strconv.Itoa(bot.userMaxLastMsg)+"\n")
}
