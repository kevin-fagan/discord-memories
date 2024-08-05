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

func Upload(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3, args []string) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "upload",
		"prefix":  args[0],
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	if !c.PrefixExists(args[0]) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s does not exist.", args[0]))
		logrus.WithFields(logs).Info("prefix does not exist")
		return
	}

	for _, attachment := range m.Attachments {
		logs["file"] = attachment.Filename
		logs["size"] = attachment.Size

		if attachment.Size > c.Storage.MaxFileSize {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("File size cannot be greater than %d bytes.", c.Storage.MaxFileSize))
			logrus.WithFields(logs).Info("uploaded rejected due to file size being too large")
			return
		}

		if !c.SupportsExtension(attachment.Filename) {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s files are not allowed.", filepath.Ext(attachment.Filename)))
			logrus.WithFields(logs).Info("uploaded rejected due to file type")
			return
		}

		err := storage.UploadObject(sv, c.Storage.Bucket, c.Commands[args[0]].Path, *attachment)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("An error has occured while uploading %s.", attachment.Filename))
			logs["error"] = err
			logrus.WithFields(logs).Error("unable to upload file")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "File succesfully uploaded!")
		logrus.WithFields(logs).Info("file uploaded")
	}
}
