// commands
package main

import (
	"fmt"
	"github.com/shkh/lastfm-go/lastfm"
	"io/ioutil"
	"net/http"
	"strings"
)

func (bot *Bot) CmdInterpreter(username string, usermessage string) {
	message := strings.ToLower(usermessage)
	tempstr := strings.Split(message, " ")

	for _, str := range tempstr {
		if strings.HasPrefix(str, "https://") || strings.HasPrefix(str, "http://") {
			bot.Message("^ " + webTitle(str))
		} else if isWebsite(str) {
			bot.Message("^ " + webTitle("http://"+str))
		}
	}

	if strings.HasPrefix(message, "!help") {
		bot.Message("For help on the bot please go to http://commandanddemand.com/bot.html")
	} else if strings.HasPrefix(message, "!quote") {
		bot.Message(bot.getQuote())
	} else if strings.HasPrefix(message, "!addquote ") {
		stringpls := strings.Replace(message, "!addquote ", "", 1)
		if bot.isMod(username) {
			bot.quotes[stringpls] = username
			bot.writeQuoteDB()
		} else {
			bot.Message(username + " you are not a mod!")
		}
	} else if strings.HasPrefix(message, "!timeout ") {
		stringpls := strings.Replace(message, "!timeout ", "", 1)
		temp1 := strings.Split(stringpls, " ")
		temp2 := strings.Replace(stringpls, temp1[0], "", 1)
		if temp2 == "" {
			temp2 = "no reason"
		}
		if bot.isMod(username) {
			bot.timeout(temp1[0], temp2)
		} else {
			bot.Message(username + " you are not a mod!")
		}
	} else if strings.HasPrefix(message, "!ban ") {
		stringpls := strings.Replace(message, "!ban ", "", 1)
		temp1 := strings.Split(stringpls, " ")
		temp2 := strings.Replace(stringpls, temp1[0], "", 1)
		if temp2 == "" {
			temp2 = "no reason"
		}
		if bot.isMod(username) {
			bot.ban(temp1[0], temp2)
		} else {
			bot.Message(username + " you are not a mod!")
		}
	} else if message == "!song" {
		api := lastfm.New("e6563970017df6d5966edfa836e12835", "dcc462ffd8a371fee5a5b49c248a2371")
		temp, _ := api.User.GetRecentTracks(lastfm.P{"user": bot.lastfm})
		var inserthere string
		if temp.Tracks[0].Date.Date != "" {
			inserthere = ". It was played on: " + temp.Tracks[0].Date.Date
		}
		bot.Message("Song: " + temp.Tracks[0].Artist.Name + " - " + temp.Tracks[0].Name + inserthere)
	}
}

//Website stuff
func webTitle(website string) string {
	response, err := http.Get(website)
	if err != nil {
		return "Error reading website"
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "Error reading website"
		}
		derp := strings.Split(string(contents), "<title>")
		derpz := strings.Split(derp[1], "</title>")
		return derpz[0]
	}
}

func isWebsite(website string) bool {
	domains := []string{".com", ".net", ".org", ".info"}
	for _, domain := range domains {
		if strings.Contains(website, domain) {
			return true
		}
	}
	return false
}

//End website stuff

//Mod stuff
func (bot *Bot) isMod(username string) bool {
	temp := strings.Replace(bot.channel, "#", "", 1)
	if bot.mods[username] == true || temp == username {
		return true
	}
	return false
}

func (bot *Bot) timeout(username string, reason string) {
	fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :/timeout "+username+"\r\n")
	bot.Message(username + " was timed out(" + reason + ")!")
}

func (bot *Bot) ban(username string, reason string) {
	fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :/ban "+username+"\r\n")
	bot.Message(username + " was banned(" + reason + ")!")
}

//End mod stuff
