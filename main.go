package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

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
	cmdArg     = 1
	flagArg    = 2
	configFile = "memories.json"
)

func main() {
	logrus.Info("starting Discord Memories")
	logrus.Info("loading custom configurations")
	config, err := config.LoadConfig(configFile)
	if err != nil {
		logrus.Fatalf("error reading custom configuration: %s\n", err)
	}

	logrus.Info("creating S3 session")
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
	err = s3helper.Sync(service, config, config.Storage.Bucket)
	if err != nil {
		logrus.Fatalf("error syncing folders: %s\n", err)
	}

	// Create a new Discord session using the provided bot token
	logrus.Info("creating Discord session")
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
	logrus.Info("creating websocket connection to Discord")
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

	logrus.WithFields(logrus.Fields{
		"author":  m.Author,
		"command": args[cmdArg],
	}).Info("message received")

	// Invoking the Help command
	if len(args) == 2 && args[cmdArg] == "help" {
		s.ChannelMessageSend(m.ChannelID, c.Help())
	}

	// Getting random content from a "folder" defined by the first argument
	if _, argExists := c.Commands[args[cmdArg]]; len(args) == 2 && argExists {
		if !c.BotAllowed(m.GuildID, m.ChannelID) {
			s.ChannelMessageSend(m.ChannelID, "This channel or server is not allowed to use this bot.")
			return
		}

		if !c.CommandAllowed(args[cmdArg]) {
			s.ChannelMessageSend(m.ChannelID, "This argument is either not allowed.")
			return
		}

		object, name, err := s3helper.GetRandomObjectUnderPrefix(service, c.Storage.Bucket, c.Commands[args[cmdArg]].Path)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			logrus.WithFields(logrus.Fields{
				"error":   err,
				"author":  m.Author,
				"command": args[cmdArg],
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
				"error":   err,
				"author":  m.Author,
				"command": args[cmdArg],
			}).Error("error occured while sending a message")
			return
		}

		// File succesfully sent
		logrus.WithFields(logrus.Fields{
			"author":  m.Author,
			"command": args[cmdArg],
			"size":    *object.ContentLength,
			"file":    name,
		}).Info("message sent")
	}

	// Uploading content into a "folder" defined by the first argument
	if len(args) == 3 && args[flagArg] == "--upload" {
		for _, attachment := range m.Attachments {
			if attachment.Size > c.Storage.MaxFileSize {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("File size cannot be greater than %d bytes.", c.Storage.MaxFileSize))
				return
			}

			if !c.SupportsExtension(attachment.Filename) {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s files are now allowed.", filepath.Ext(attachment.Filename)))
				return
			}

			err := s3helper.UploadObject(service, c.Storage.Bucket, c.Commands[args[cmdArg]].Path, *attachment)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("An error has occured while uploading %s.", attachment.Filename))
				logrus.WithFields(logrus.Fields{
					"error":   err,
					"author":  m.Author,
					"command": args[cmdArg],
					"size":    attachment.Size,
					"file":    filepath.Join(c.Commands[args[cmdArg]].Path, attachment.Filename),
				}).Error("error while uploading a file")
				return
			}

			// File was succesfully uploaded
			s.ChannelMessageSend(m.ChannelID, "File succesfully uploaded!")
			logrus.WithFields(logrus.Fields{
				"author":  m.Author,
				"command": args[cmdArg],
				"size":    attachment.Size,
				"file":    filepath.Join(c.Commands[args[cmdArg]].Path, attachment.Filename),
			}).Info("file uploaded")
		}
	}
}
