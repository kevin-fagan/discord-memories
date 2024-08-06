package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/KevinFagan/discord-memories/config"
	"github.com/KevinFagan/discord-memories/storage"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Upload(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3, option string) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "upload",
		"prefix":  option,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	if !c.OptionExists(option) {
		logrus.WithFields(logs).Info("prefix does not exist")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s does not exist.", option))
		return
	}

	for _, attachment := range m.Attachments {
		logs["file"] = attachment.Filename
		logs["size"] = attachment.Size

		if attachment.Size > c.Storage.MaxFileSize {
			logrus.WithFields(logs).Info("uploaded rejected due to file size being too large")
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("File size cannot be greater than %d bytes.", c.Storage.MaxFileSize))
			return
		}

		if !c.SupportsExtension(attachment.Filename) {
			logrus.WithFields(logs).Info("uploaded rejected due to file type")
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s files are not allowed.", filepath.Ext(attachment.Filename)))
			return
		}

		err := storage.UploadObject(sv, c.Storage.Bucket, c.Options[option].Path, *attachment)
		if err != nil {
			logs["error"] = err
			logrus.WithFields(logs).Error("unable to upload file")
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error while uploading %s: %s", attachment.Filename, err))
			return
		}

		logrus.WithFields(logs).Info("file uploaded")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s uploaded.", attachment.Filename))
	}
}
