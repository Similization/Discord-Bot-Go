package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const BotPrefix = "!gobot"

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var TOKEN string = os.Getenv("BOT_TOKEN")

	discord, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(hello)
	discord.AddHandler(flipCoin)
	discord.AddHandler(rollDice)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func hello(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	args := strings.Split(msg.Content, " ")
	if args[0] != BotPrefix {
		return
	}

	if strings.Join(args[1:], " ") == "hello" {
		sess.ChannelMessageSend(msg.ChannelID, "Hi, "+msg.Author.Username)
	}
}

func flipCoin(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	args := strings.Split(msg.Content, " ")
	if args[0] != BotPrefix {
		return
	}

	if strings.Join(args[1:], " ") == "flip coin" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		res := r.Intn(2)
		if res == 0 {
			sess.ChannelMessageSend(msg.ChannelID, msg.Author.Username+" flipped tails")
		} else {
			sess.ChannelMessageSend(msg.ChannelID, msg.Author.Username+" flipped heads")
		}
	}
}

func rollDice(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	args := strings.Split(msg.Content, " ")
	if args[0] != BotPrefix {
		return
	}

	if strings.Join(args[1:], " ") == "roll dice" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		res := r.Intn(6) + 1
		sess.ChannelMessageSend(msg.ChannelID, msg.Author.Username+" rolled "+fmt.Sprint(res))
	}
}
