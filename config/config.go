package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/joho/godotenv"
)

type Config struct {
	Tokens      Tokens
	Storage     Storage             `json:"storage"`
	Arguments   map[string]Argument `json:"arguments"`
	Permissions Permissions         `json:"permissions"`
}

type Tokens struct {
	DiscordToken string
	S3AccessKey  string
	S3SecretKey  string
}

type Storage struct {
	Endpoint    string   `json:"endpoint"`
	Region      string   `json:"region"`
	Bucket      string   `json:"bucket"`
	MaxFileSize int      `json:"maxFileSize"`
	Extensions  []string `json:"extensions"`
}

type Argument struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type Permissions struct {
	Servers  map[string]Permission `json:"servers"`
	Channels map[string]Permission `json:"channels"`
}

type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// LoadConfig will load the configuration file along with required environment variables
func LoadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}

	// Loading environment variables
	err = godotenv.Load()
	if err != nil {
		return Config{}, nil
	}
	config.Tokens.DiscordToken = os.Getenv("DISCORD_TOKEN")
	config.Tokens.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	config.Tokens.S3SecretKey = os.Getenv("S3_SECRET_KEY")

	return config, nil
}

// ArgumentAllowed checks if an argument is enabled or not
func (c Config) ArgumentAllowed(arg string) bool {
	argument, ok := c.Arguments[arg]
	if !ok {
		return false
	}
	return argument.Enabled
}

// BotAllowed checks if the bot is allowed in either the Server or channel
// Channel permissions override server permissions
func (c Config) BotAllowed(serverID, channelID string) bool {
	channel, cok := c.Permissions.Channels[channelID]
	server, sok := c.Permissions.Servers[serverID]

	if !cok && !sok {
		return false
	}
	if cok {
		return channel.Enabled
	}

	return server.Enabled
}

// SupportsExtension checks if an extension is allowed or not
func (c Config) SupportsExtension(file string) bool {
	supported := false
	for _, extension := range c.Storage.Extensions {
		if filepath.Ext(file) == extension {
			supported = true
		}
	}
	return supported
}

// Help returns a formatted help string based off the provided configs
func (c Config) Help() string {
	var sb strings.Builder

	sb.WriteString("```\n")
	sb.WriteString("The Discord Memories bot allows you to upload and recall memories made with your friends.\nCommands, permissions, file types, and file sizes are all determined by the Memories\nconfiguration file.\n\n")
	sb.WriteString("Usage:\n")
	sb.WriteString("\t !memories [command]\n\n")
	sb.WriteString("Commands:\n")

	w := tabwriter.NewWriter(&sb, 0, 0, 4, ' ', 0)
	for k, v := range c.Arguments {
		fmt.Fprintf(w, "\t%s\t%s\n", k, v.Description)
	}
	w.Flush()

	sb.WriteString("\n\n")
	sb.WriteString("Flags:\n")
	sb.WriteString("\t--upload\tAllows you to upload a file for a given command\n\n")
	sb.WriteString("```\n")

	return sb.String()
}
