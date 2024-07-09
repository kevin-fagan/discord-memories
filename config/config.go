package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Channels []Channel `json:"channels"`
	Commands []Command `json:"commands"`
}

type Channel struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type Command struct {
	Command string `json:"command"`
	Images  string `json:"images"`
	Enabled bool   `json:"enabled"`
}

func ReadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c Config) ChannelAllowed(ID string) bool {
	for _, channel := range c.Channels {
		if channel.ID == ID {
			return channel.Enabled
		}
	}
	return false
}

func (c Config) CommandAllowed(command string) (Command, bool) {
	for _, cmd := range c.Commands {
		if cmd.Command == command {
			return cmd, cmd.Enabled
		}
	}
	return Command{}, false
}
