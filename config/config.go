package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// Help returns a formatted help string
func (c Config) Help() string {
	var sb strings.Builder

	sb.WriteString("Help Command:\n\n")
	sb.WriteString("!memories <command> [upload]\n\n")
	for k, v := range c.Arguments {
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", k, v.Description))
	}

	return sb.String()
}
