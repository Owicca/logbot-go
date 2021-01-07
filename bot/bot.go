package bot

import (
	"fmt"
	"net"
	"strings"
)

/*
commands:
- PASS
- NICK
- USER
- JOIN
- PONG
- PRIVMSG
- QUIT
*/

type Bot struct {
	conn     *net.Conn
	endline  string
	loggedIn bool
}

func NewBot(conn *net.Conn) *Bot {
	bot := &Bot{}
	bot.conn = conn
	bot.endline = "\n\r"
	bot.loggedIn = false

	return bot
}

func (b *Bot) IsLoggedIn() bool {
	return b.loggedIn
}

func (b *Bot) LoggedIn(isLoggedIn bool) {
	b.loggedIn = isLoggedIn
}

func (b *Bot) Pong(key string) (bool, error) {
	res, err := fmt.Fprintf(*b.conn, "PONG :%s%s",
		key, b.endline)
	return res > 0, err
}

func (b *Bot) Pass(pass string) (bool, error) {
	res, err := fmt.Fprintf(*b.conn, "PASS %s%s",
		pass, b.endline)
	return res > 0, err
}

func (b *Bot) Nick(nick string) (bool, error) {
	res, err := fmt.Fprintf(*b.conn, "NICK %s%s",
		nick, b.endline)
	return res > 0, err
}

func (b *Bot) User(username string, hostname string,
	servername string, realname string) (bool, error) {
	res, err := fmt.Fprintf(*b.conn, "USER %s %s %s :%s%s",
		username, hostname, servername, realname, b.endline)
	return res > 0, err
}

func (b *Bot) Join(channels []string) (bool, error) {
	res, err := fmt.Fprintf(*b.conn, "JOIN %s%s",
		strings.Join(channels, ","), b.endline)
	return res > 0, err
}

func (b *Bot) PrivMsg(message string) bool {
	return true
}

func (b *Bot) Quit(message string) (bool, error) {
	res, err := fmt.Fprintf(*b.conn, "QUIT :%s", message)
	return res > 0, err
}
