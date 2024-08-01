package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/KevinFagan/discord-memories/config"
	s3helper "github.com/KevinFagan/discord-memories/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	firstArg   = 1
	secondArg  = 2
	configFile = "memories.json"
)

func main() {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		logrus.Fatalf("error reading config: %s\n", err)
	}

	// Establishing S3 bucket Session
	sess, err := session.NewSession(&aws.Config{
		Region:      &config.Storage.Region,
		Endpoint:    &config.Storage.Endpoint,
		Credentials: credentials.NewStaticCredentials(config.Tokens.S3AccessKey, config.Tokens.S3SecretKey, ""),
	})
	if err != nil {
		logrus.Fatalf("error creating session: %s\n", err)
	}
	service := s3.New(sess)
	err = s3helper.Sync(service, config, config.Storage.Bucket)
	if err != nil {
		logrus.Fatalf("error syncing folders: %s\n", err)
	}

	// Create a new Discord session using the provided bot token
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
	err = dg.Open()
	if err != nil {
		logrus.Fatalf("error opening connection: %s\n", err)
	}

	// Wait here until CTRL-C or other term signal is received
	logrus.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, service *s3.S3) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Ignore all messages that do not start with the command prefix
	if !strings.HasPrefix(m.Content, "!memories") {
		return
	}
	// Ensuring we have the correct amount of argument
	args := strings.Split(strings.TrimSpace(m.Content), " ")
	if len(args) != 2 && len(args) != 3 {
		return
	}
	// Invoking the Help command
	if len(args) == 2 && args[firstArg] == "help" {
		s.ChannelMessageSend(m.ChannelID, c.Help())
	}
	// Getting random content from a "folder" defined by the first argument
	if _, argExists := c.Arguments[args[firstArg]]; len(args) == 2 && argExists {
		// Checking Sever/Channel permissions
		if !c.BotAllowed(m.GuildID, m.ChannelID) {
			s.ChannelMessageSend(m.ChannelID, "This channel or server is not allowed to use this bot.")
			return
		}
		// Checking argument permissions
		if !c.ArgumentAllowed(args[firstArg]) {
			s.ChannelMessageSend(m.ChannelID, "This argument is either not allowed.")
			return
		}

		object, name, err := s3helper.GetRandomObjectUnderPrefix(service, c.Storage.Bucket, c.Arguments[args[firstArg]].Path)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"author": m.Author,
			}).Error("error while retrieving an object")
			return
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   name,
					Reader: object.Body,
				},
			},
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"author": m.Author,
			}).Error("error occured while sending a message")
		}
	}
	// Uploading content into a "folder" defined by the first argument
	if len(args) == 3 && args[secondArg] == "upload" {
		for _, attachment := range m.Attachments {
			if attachment.Size > c.Storage.MaxFileSize {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("File size cannot be greater than %d bytes.", c.Storage.MaxFileSize))
				return
			}

			if !c.SupportsExtension(attachment.Filename) {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s files are now allowed.", filepath.Ext(attachment.Filename)))
				return
			}

			err := s3helper.UploadObject(service, c.Storage.Bucket, c.Arguments[args[firstArg]].Path, *attachment)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("An error has occured while uploading %s.", attachment.Filename))
				logrus.WithFields(logrus.Fields{
					"error":  err,
					"author": m.Author,
					"size":   attachment.Size,
					"file":   filepath.Join(c.Arguments[args[firstArg]].Path, attachment.Filename),
				}).Error("error while uploading a file")
				return
			}

			// File was succesfully uploaded
			s.ChannelMessageSend(m.ChannelID, "File succesfully uploaded!")
			logrus.WithFields(logrus.Fields{
				"author": m.Author,
				"size":   attachment.Size,
				"file":   filepath.Join(c.Arguments[args[firstArg]].Path, attachment.Filename),
			}).Info("file successfully uploaded")
		}
	}
}
