package main

import (
	"bufio"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type Bot struct {
	server        string
	port          string
	nick          string
	channel       string
	pread, pwrite chan string
	conn          net.Conn
}

func NewBot() *Bot {
	return &Bot{server: "irc.twitch.tv",
		port:    "6667",
		nick:    "quanticbot", //Change to your Twitch username
		channel: "#vaultpls", //Change to your channel
		conn:    nil}
}
func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s)\n", bot.server, bot.conn.RemoteAddr())
	return bot.conn, nil
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
func CmdInterpreter(conn net.Conn, channel string, username string, usermessage string) {
	message := strings.ToLower(usermessage)
	if strings.HasPrefix(message, "!hi") {
		Message(conn, channel, "Hi there!")
	} else if strings.HasPrefix(message, "!tredo") {
		Message(conn, channel, "\"can't have a good night without some cock\" -Tredo 2013")
	} else if strings.HasPrefix(message, "http://") {
		Message(conn, channel, "^ " + webTitle(usermessage))
	} else if strings.HasPrefix(message, "https://") {
		Message(conn, channel, "^ " + webTitle(usermessage))
	}
}

func Message(conn net.Conn, channel string, message string) {
	fmt.Printf("Bot: " + message + "\n")
	fmt.Fprintf(conn, "PRIVMSG "+channel+" :"+message+"\r\n")
}

func AutoMessage(conn net.Conn, channel string) {
	for {
		time.Sleep(10 * time.Minute)
		Message(conn, channel, "Hi this is a automessage.  Please change me, senpai.")
	}
}

func ConsoleInput(conn net.Conn, channel string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if text != "" {
			Message(conn, channel, text)
		}
	}
}

func main() {
	//INIT
	ircbot := NewBot()
	conn, _ := ircbot.Connect()
	messagesCount := 0
	pass1, err := ioutil.ReadFile("twitch_pass")
	pass := strings.Replace(string(pass1), "\n", "", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pass))
	go AutoMessage(conn, ircbot.channel)
	fmt.Printf("Initialized auto messager...\n")
	fmt.Fprintf(conn, "USER %s 8 * :%s\r\n", ircbot.nick, ircbot.nick)
	fmt.Fprintf(conn, "PASS %s\r\n", pass)
	fmt.Fprintf(conn, "NICK %s\r\n", ircbot.nick)
	fmt.Fprintf(conn, "JOIN %s\r\n", ircbot.channel)
	fmt.Println("Inserted credentials...\n")
	defer conn.Close()
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	go ConsoleInput(conn, ircbot.channel)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors
		}
		if strings.Contains(line, "PING") {
			pongdata := strings.Split(line, "PING ")
			fmt.Fprintf(conn, "PONG %s\r\n", pongdata[1])
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG "+ircbot.channel) {
			messagesCount = messagesCount + 1
			if messagesCount == 50 {
				Message(conn, ircbot.channel, "Welcome to Vaults stream!  Please take a break from your usual life and watch.  You'd be surprised how much retardation goes on in this stream.")
			}
			userdata := strings.Split(line, ".tmi.twitch.tv PRIVMSG "+ircbot.channel)
			username := strings.Split(userdata[0], "@")
			usermessage := strings.Replace(userdata[1], " :", "", 1)
			fmt.Printf(username[1] + ": " + usermessage + "\n") //TODO: Put this in a command interpretter
			go CmdInterpreter(conn, ircbot.channel, username[1], usermessage)

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
