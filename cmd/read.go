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

func Read(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3, option string) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "read",
		"prefix":  option,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	if !c.OptionExists(option) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s does not exist.", option))
		logrus.WithFields(logs).Info("prefix does not exist")
		return
	}

	object, name, err := storage.GetRandomObjectUnderPrefix(sv, c.Storage.Bucket, c.Options[option].Path)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retieve object from s3")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %s", err))
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
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to send discord message")
	}

	logs["file"] = name
	logs["size"] = *object.ContentLength
	logrus.WithFields(logs).Info("file retrieved")
}
