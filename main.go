package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	discordToken = os.Getenv("DISCORD_BOT_TOKEN")
	imageDir     = os.Getenv("FILE_PATH")
)

func main() {
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Printf("error creating Discord session: %s", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Printf("error opening connection: %s", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!jake") {
		fmt.Println("Command received")
		imagePath := filepath.Join(imageDir, selectPhoto(imageDir))

		file, err := os.Open(imagePath)
		if err != nil {
			fmt.Printf("error opening file: %s", err)
		}
		defer file.Close()

		message := &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   filepath.Base(imagePath),
					Reader: file,
				},
			},
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, message)
		if err != nil {
			fmt.Printf("error sending message: %s", err)
		}
	}
}

func selectPhoto(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("error reading directory: %s", err)
	}
	num := rand.Intn(len(files))
	return files[num].Name()
}
