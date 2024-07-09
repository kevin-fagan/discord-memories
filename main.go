package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/KevinFagan/project-friendship/config"
	"github.com/bwmarrin/discordgo"
)

var (
	discordToken = os.Getenv("DISCORD_BOT_TOKEN")
)

const (
	configFile = "config/config.json"
)

func main() {
	// Load configuration
	config, err := config.ReadConfig(configFile)
	if err != nil {
		fmt.Printf("error reading config: %s\n", err)
		return
	}

	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Printf("error creating Discord session: %s\n", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, config)
	})
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Printf("error opening connection: %s\n", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	reqCmd := strings.Split(m.Content, " ")[0]
	cmd, cmdAllowed := c.CommandAllowed(reqCmd)

	// Checking if the channel and the command is allowed
	if !c.ChannelAllowed(m.ChannelID) {
		s.ChannelMessageSend(m.ChannelID, "This channel is not allowed to use this bot.")
		return
	}
	if !cmdAllowed {
		s.ChannelMessageSend(m.ChannelID, "This command is either not allowed or does not exist.")
		return
	}

	// Selecting and reading random photo based of the provided command
	photo, err := selectPhoto(cmd.Images)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error selecting photo: %s", err))
	}
	file, err := os.Open(filepath.Join(cmd.Images, photo))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error reading photo: %s", err))
	}
	defer file.Close()

	// Sending the selected photo to the channel
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "Random Photo",
				Reader: file,
			},
		},
	})

	if err != nil {
		fmt.Printf("error sending message: %s", err)
	}
}

func selectPhoto(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(len(files))
	return files[num].Name(), nil
}
