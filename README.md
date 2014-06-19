<h1>Twitch IRC Bot</h1>
Created in Go!


Original: https://stackoverflow.com/questions/13342128/simple-golang-irc-bot-keeps-timing-out

<h2>Features:</h2>
**Auto-messenger** - Messages every (x) minutes and every (x) lines

**Interpretation of Commands** - Self explanatory "!hi" = "Hi there!"

**Notifications** - (x) joined the chatroom! (x) got moderator status! (only in text)

**Website Title fetcher** - (WIP) Someone posted a link and the bot goes and fetches the title of the page so you can know what it is before you even go there!

**Console Input** - Don't feel like going in your browser/IRC client? You can type your input into the program and it comes out if the bot were saying it!

<h3>Why not use a IRC Lib instead of making your own?</h3>
It just seemed a little dirty using a IRC Library at the moment.

<h4>How to use</h4>
'''
git clone https://github.com/Vaultpls/Twitch-IRC-Bot
or Download ZIP -> On the right sidebar

Create twitch_pass and insert the oauth from
http://twitchapps.com/tmi/

Then either "go run bot.go" or "go build bot.go"
'''

Please add/remove/modify as you please!  :D
