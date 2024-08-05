package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/KevinFagan/discord-memories/config"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config, sv *s3.S3) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "help",
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 0, 4, ' ', 0)
	fmt.Fprint(w, "```\n")
	fmt.Fprint(w, "The Discord Memories bot allows you to upload and recall memories made with your friends.\nCommands, permissions, file types, and file sizes are all determined by the Memories\nconfiguration file.\n\n")
	fmt.Fprint(w, "Usage:\n")
	fmt.Fprint(w, "\t!memories [command] [option]\n\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "\thelp\tPrints information about the Discord Memories bot\n")
	fmt.Fprint(w, "\tcount\tCounts the number of files under an option\n")
	fmt.Fprint(w, "\tread\tRetrieves a random file under an option\n")
	fmt.Fprint(w, "\tupload\tUploads one or more files under an option\n")
	fmt.Fprint(w, "\tchannels\tLists channels that have permissions to use this bot\n")
	fmt.Fprint(w, "\tservers\tLists servers that have permissions to use this bot\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Options:\n")
	for k, v := range c.Commands {
		fmt.Fprintf(w, "\t%s\t%s\n", k, v.Description)
	}
	fmt.Fprint(w, "```\n")
	w.Flush()

	s.ChannelMessageSend(m.ChannelID, sb.String())
}
