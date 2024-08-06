package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/KevinFagan/discord-memories/cmd"
	"github.com/KevinFagan/discord-memories/config"
	"github.com/KevinFagan/discord-memories/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	configFile = "memories.json"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	logrus.Info("loading configurations...")
	config, err := config.LoadConfig(configFile)
	if err != nil {
		logrus.Fatalf("error reading custom configuration: %s\n", err)
	}

	logrus.Info("creating S3 session...")
	sess, err := session.NewSession(&aws.Config{
		Region:      &config.Storage.Region,
		Endpoint:    &config.Storage.Endpoint,
		Credentials: credentials.NewStaticCredentials(config.Tokens.S3AccessKey, config.Tokens.S3SecretKey, ""),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	})
	if err != nil {
		logrus.Fatalf("error creating session: %s\n", err)
	}

	service := s3.New(sess)
	err = storage.Sync(service, config, config.Storage.Bucket)
	if err != nil {
		logrus.Fatalf("error syncing folders: %s\n", err)
	}

	// Create a new Discord session using the provided bot token
	logrus.Info("creating Discord session...")
	dg, err := discordgo.New("Bot " + config.Tokens.DiscordToken)
	if err != nil {
		logrus.Fatalf("error creating Discord session: %s\n", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, config, service)
	})
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening
	logrus.Info("opening websocket connection to Discord...")
	err = dg.Open()
	if err != nil {
		logrus.Fatalf("error opening connection: %s\n", err)
	}

	// Wait here until CTRL-C or other term signal is received
	logrus.Info("Discord Memories is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!memories") {
		return
	}

	if !c.BotAllowed(m.GuildID, m.ChannelID) {
		s.ChannelMessageSend(m.ChannelID, "This Discord Server or Channel does not have permissions to run this bot.")
		return
	}

	args := strings.Split(strings.TrimSpace(m.Content), " ")
	switch command := args[1]; command {
	case "help":
		cmd.Help(s, m, c, sv)
	case "count":
		cmd.Count(s, m, c, sv, args[2:])
	case "upload":
		cmd.Upload(s, m, c, sv, args[2:])
	case "read":
		cmd.Read(s, m, c, sv, args[2:])
	case "servers":
		cmd.Servers(s, m, c)
	case "channels":
		cmd.Channels(s, m, c)
	}
}
