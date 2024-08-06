package cmd

import (
	"fmt"

	"github.com/KevinFagan/discord-memories/config"
	"github.com/KevinFagan/discord-memories/storage"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Count(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3, option string) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "count",
		"prefix":  option,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	count, err := storage.Count(sv, c.Storage.Bucket, option)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to count the number of files under an option")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %s", err))
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s contains %d files.", option, count))
}
