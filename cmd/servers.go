package cmd

import (
	"fmt"
	"strings"

	"github.com/KevinFagan/discord-memories/config"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Servers(s *discordgo.Session, m *discordgo.MessageCreate, c config.Config) {
	logs := logrus.Fields{
		"author":  m.Author.Username,
		"command": "servers",
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command received")

	var sb strings.Builder

	sb.WriteString("Server allowlist:\n")
	for k := range c.Permissions.Servers {
		st, _ := s.Guild(k)
		sb.WriteString(fmt.Sprintf("- %s\n", st.Name))
	}

	s.ChannelMessageSend(m.ChannelID, sb.String())
}
