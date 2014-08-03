package main

import (
	"bufio"
	"math/rand"
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
	quotes		map[string]string
	mods		map[string]bool
	userLastMsg	map[string]int64
	lastmsg		int64
	maxMsgTime	int64
	userMaxLastMsg	int64
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
		quotes:		make(map[string]string),
		mods:		make(map[string]bool),
		lastmsg:	0,
		maxMsgTime:	5,
		userLastMsg:make(map[string]int64),
		userMaxLastMsg:2,
	}
}

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
	dst, err := os.Create("quotes.txt")
	defer dst.Close()
	if err != nil {
		fmt.Println("Can't write to QuoteDB!")
		return
	}
	for split1, split2 := range bot.quotes {
		fmt.Fprintf(dst, split1 + "|" + split2 + "\n")
	}
}

func (bot *Bot) readQuoteDB() {
	quotes, err := ioutil.ReadFile("quotes.txt")
	if err != nil {
		fmt.Println("Unable to read QuoteDB")
		return
	}
	split1 := strings.Split(string(quotes), "\n")
	for _, splitted1 := range split1 {
		split2 := strings.Split(splitted1, "|")
		bot.quotes[split2[0]] = split2[1]
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

func (bot *Bot) isMod(username string) bool {
	temp := strings.Replace(bot.channel, "#", "", 1)
	if bot.mods[username] == true || temp == username{
		return true
	}
	return false
}

/*
TODO: Add more fun and interesting commands into
the Command Interpreter.
*/
func (bot *Bot) CmdInterpreter(username string, usermessage string) {
	message := strings.ToLower(usermessage)
	tempstr := strings.Split(message, " ")
	
	for _, str := range tempstr {
		if strings.HasPrefix(str, "https://") || strings.HasPrefix(str, "http://") {
			bot.Message("^ " + webTitle(str))
		} else if isWebsite(str) {
			bot.Message("^ " + webTitle("http://" + str))
		}
	}
	
	if strings.HasPrefix(message, "!hi") {
		bot.Message("Hi there!")
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
	}
}

func isWebsite(website string) bool {
	domains := []string {".com", ".net", ".org", ".info",}
	for _, domain := range domains {
		if strings.Contains(website, domain) {
			return true
		}
	}
	return false
}

func (bot *Bot) Message(message string) {
	if bot.lastmsg + bot.maxMsgTime <= time.Now().Unix() {
		fmt.Printf("Bot: " + message + "\n")
		fmt.Fprintf(bot.conn, "PRIVMSG "+ bot.channel +" :"+message+"\r\n")
		bot.lastmsg = time.Now().Unix()
	} else {
		fmt.Println("Attempted to spam message")
	}
}

func (bot *Bot) timeout (username string, reason string) {
	fmt.Fprintf(bot.conn, "PRIVMSG "+ bot.channel +" :/timeout " + username + "\r\n")
	bot.Message(username + " was timed out(" + reason +")!")
}

func (bot *Bot) ban (username string, reason string) {
	fmt.Fprintf(bot.conn, "PRIVMSG "+ bot.channel +" :/ban " + username + "\r\n")
	bot.Message(username + " was banned(" + reason +")!")
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
		if text == "/quit" {
			bot.conn.Close()
			os.Exit(0)
		}
		if text != "" {
			bot.Message(text)
		}
	}
}

func main() {
	//INIT
	fmt.Printf("Twitch IRC Bot made in Go!\n")
	ircbot := NewBot()
	go ircbot.ConsoleInput()
	ircbot.Connect()
	messagesCount := 0
	pass1, err := ioutil.ReadFile("twitch_pass.txt")
	pass := strings.Replace(string(pass1), "\n", "", 0)
	if err != nil {
		panic(err)
	}
	ircbot.readQuoteDB()
	fmt.Fprintf(ircbot.conn, "USER %s 8 * :%s\r\n", ircbot.nick, ircbot.nick)
	fmt.Fprintf(ircbot.conn, "PASS %s\r\n", pass)
	fmt.Fprintf(ircbot.conn, "NICK %s\r\n", ircbot.nick)
	fmt.Fprintf(ircbot.conn, "JOIN %s\r\n", ircbot.channel)
	fmt.Printf("Inserted information to server...\n")
	fmt.Printf("If you don't see the stream chat it probably means the Twitch oAuth password is wrong\n")
	defer ircbot.conn.Close()
	go ircbot.AutoMessage()
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
			fmt.Printf(username[1] + ": " + usermessage + "\n")
			if ircbot.userLastMsg[username[1]] + ircbot.userMaxLastMsg >= time.Now().Unix() {
				ircbot.timeout(username[1], "spam")
			}
			ircbot.userLastMsg[username[1]] = time.Now().Unix()
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
			ircbot.mods[usermod[1]] = true
			fmt.Printf(usermod[1] + " is a moderator!\n")
		} else if strings.Contains(line, ":jtv MODE "+ircbot.channel+" -o ") {
			usermod := strings.Split(line, ":jtv MODE "+ircbot.channel+" -o ")
			ircbot.mods[usermod[1]] = false
			fmt.Printf(usermod[1] + " isn't a moderator anymore!\n")
		}
	}

}
