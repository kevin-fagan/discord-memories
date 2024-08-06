package cmd

import (
	"fmt"
	"strings"

	"github.com/KevinFagan/discord-memories/config"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Channels(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "channels",
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	var sb strings.Builder
	sb.WriteString("Channel allowlist:\n")
	for k := range c.Permissions.Channels {
		st, _ := s.Channel(k)
		sb.WriteString(fmt.Sprintf("- %s\n", st.Name))
	}

	s.ChannelMessageSend(m.ChannelID, sb.String())
}
