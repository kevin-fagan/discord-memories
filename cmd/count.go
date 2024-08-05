package cmd

import (
	"github.com/KevinFagan/discord-memories/config"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Count(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3, args []string) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "count",
		"prefix":  args[0],
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")
}
