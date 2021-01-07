package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	botlib "github.com/Owicca/logbot-go/bot"
	bufferlib "github.com/Owicca/logbot-go/buffer"
)

/*
var Reserved = map[string]string
var ErrorReplies = map[string]string
var CommandResponses = map[string]string
*/
/*
correctly create buffers and close them
open files in each buffer
*/

const (
	host = "sp44.local"
	port = 6667
)

var (
	nick         = "nick"
	user         = "user"
	real_name    = "real_name"
	hostname     = "hostname"
	new_channels = []string{"#test_chan", "#another_chan"}
)

var (
	channels = map[string]*bufferlib.Buffer{}
)

func main() {
	verbose := flag.Bool("v", false, "Print verbose output")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Could not dial %s (%s)", addr, err)
	}

	bot := botlib.NewBot(&conn)
	stopSig := SetupGracefulStop(bot)

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if *verbose {
				log.Printf("Something happened while reading line (%s)\n", err)
			}
		}

		if strings.HasPrefix(line, "PING") {
			key := strings.Split(line, " ")[1][1:]
			bot.Pong(key)

			if *verbose {
				log.Println("ponged to ", key)
			}

			continue
		}

		if !bot.IsLoggedIn() {
			bot.Nick(nick)
			if *verbose {
				log.Println("Nick ", nick)
			}

			bot.User(user, hostname, host, real_name)
			if *verbose {
				log.Println("User ", user)
			}

			bot.LoggedIn(true)
		}

		if bot.IsLoggedIn() || strings.Contains(line, "376 "+nick) {
			for _, chn := range new_channels {
				channel, ok := channels[chn]

				if !ok || !channel.IsConnected() {
					channels[chn] = bufferlib.NewBuffer(chn)
					bot.Join([]string{chn})
					if *verbose {
						log.Printf("Joined %s\n", chn)
					}
				}
			}
		}

		if strings.HasPrefix(line, "PRIVMSG") {
			parts := strings.Split(line, "PRIVMSG")
			log.Fatal(parts)
		}

		if strings.HasPrefix(line, "ERROR") {
			if *verbose {
				log.Fatalf("Quit because (%s)", line)
			}

			stopSig <- syscall.SIGINT
		}

		if *verbose {
			log.Println(line)
		}
	}
}

func SetupGracefulStop(bot *botlib.Bot) chan os.Signal {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt)

	go func() {
		sig := <-gracefulStop

		bot.Quit("Buh'bye")
		log.Println("Quiting with 'Buh'bye")
		log.Printf("Stopping because %+v\n", sig)

		for _, chn := range channels {
			chn.Close()
		}

		time.Sleep(1 * time.Second)

		os.Exit(0)
	}()

	return gracefulStop
}
