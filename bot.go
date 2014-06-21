package main

import (
	"bufio"
	"fmt"
	"net/http"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type Bot struct {
	server		string
	port		string
	nick		string
	channel		string
	autoMSG1	string
	autoMSG2	string
	autoMSG2Count	int
	conn		net.Conn
}

func NewBot() *Bot {
	return &Bot{
		server:		"irc.twitch.tv",
		port:		"6667",
		nick:		"quanticbot", //Change to your Twitch username
		channel:	"#vaultpls", //Change to your channel
		autoMSG1:	"I am a timed auto-message!",
		autoMSG2:	"I am a line counting auto-message!",
		autoMSG2Count:	50,
		conn:		nil, //Don't change this
	}
}
func (bot *Bot) Connect() {
	var err error
	bot.conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		fmt.Printf("Unable to connect to Twitch IRC server! Reconnecting in 10 seconds...\n")
		time.Sleep(10 * time.Second)
                bot.Connect()
	}
	fmt.Printf("Connected to IRC server %s (%s)\n", bot.server, bot.conn.RemoteAddr())
}

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

/*
TODO: Add more fun and interesting commands into
the Command Interpreter.
*/
func (bot *Bot) CmdInterpreter(username string, usermessage string) {
	message := strings.ToLower(usermessage)
	if strings.HasPrefix(message, "!hi") {
		bot.Message("Hi there!")
	} else if strings.HasPrefix(message, "!tredo") {
		bot.Message("\"can't have a good night without some cock\" -Tredo 2013")
	} else if strings.HasPrefix(message, "http://") {
		bot.Message("^ " + webTitle(usermessage))
	} else if strings.HasPrefix(message, "https://") {
		bot.Message("^ " + webTitle(usermessage))
	}
}

func (bot *Bot) Message(message string) {
	fmt.Printf("Bot: " + message + "\n")
	fmt.Fprintf(bot.conn, "PRIVMSG "+ bot.channel +" :"+message+"\r\n")
}

func (bot *Bot) AutoMessage() {
	for {
		time.Sleep(10 * time.Minute)
		bot.Message(bot.autoMSG1)
	}
}

func (bot *Bot) ConsoleInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if text != "" {
			bot.Message(text)
		}
	}
}

func main() {
	//INIT
	fmt.Printf("Twitch IRC Bot made in Go!\n")
	ircbot := NewBot()
	ircbot.Connect()
	messagesCount := 0
	pass1, err := ioutil.ReadFile("twitch_pass")
	pass := strings.Replace(string(pass1), "\n", "", 0)
	if err != nil {
		panic(err)
	}
	go ircbot.AutoMessage()
	fmt.Fprintf(ircbot.conn, "USER %s 8 * :%s\r\n", ircbot.nick, ircbot.nick)
	fmt.Fprintf(ircbot.conn, "PASS %s\r\n", pass)
	fmt.Fprintf(ircbot.conn, "NICK %s\r\n", ircbot.nick)
	fmt.Fprintf(ircbot.conn, "JOIN %s\r\n", ircbot.channel)
	fmt.Printf("Inserted information to server...\n")
	fmt.Printf("If you don't see the stream chat it probably means the Twitch oAuth password is wrong\n")
	defer ircbot.conn.Close()
	reader := bufio.NewReader(ircbot.conn)
	tp := textproto.NewReader(reader)
	go ircbot.ConsoleInput()
	for {
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors
		}
		if strings.Contains(line, "PING") {
			pongdata := strings.Split(line, "PING ")
			fmt.Fprintf(ircbot.conn, "PONG %s\r\n", pongdata[1])
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG "+ircbot.channel) {
			messagesCount++
			if messagesCount == ircbot.autoMSG2Count {
				ircbot.Message(ircbot.autoMSG2)
			}
			userdata := strings.Split(line, ".tmi.twitch.tv PRIVMSG "+ircbot.channel)
			username := strings.Split(userdata[0], "@")
			usermessage := strings.Replace(userdata[1], " :", "", 1)
			fmt.Printf(username[1] + ": " + usermessage + "\n") //TODO: Put this in a command interpretter
			go ircbot.CmdInterpreter(username[1], usermessage)

		} else if strings.Contains(line, ".tmi.twitch.tv JOIN "+ircbot.channel) {
			userjoindata := strings.Split(line, ".tmi.twitch.tv JOIN "+ircbot.channel)
			userjoined := strings.Split(userjoindata[0], "@")
			fmt.Printf(userjoined[1] + " has joined!\n") //TODO: Database for all people joined and to have it check to see if they have joined.  If not, greet them
		} else if strings.Contains(line, ".tmi.twitch.tv PART "+ircbot.channel) {
			userjoindata := strings.Split(line, ".tmi.twitch.tv PART "+ircbot.channel)
			userjoined := strings.Split(userjoindata[0], "@")
			fmt.Printf(userjoined[1] + " has left!\n")
		} else if strings.Contains(line, ":jtv MODE "+ircbot.channel+" +o ") {
			usermod := strings.Split(line, ":jtv MODE "+ircbot.channel+" +o ")
			fmt.Printf(usermod[1] + " is a moderator!\n")
		} else if strings.Contains(line, ":jtv MODE "+ircbot.channel+" -o ") {
			usermod := strings.Split(line, ":jtv MODE "+ircbot.channel+" -o ")
			fmt.Printf(usermod[1] + " isn't a moderator anymore!\n")
		}
	}

}
