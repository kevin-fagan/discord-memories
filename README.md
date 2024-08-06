# Discord Memories

Discord Memories is a bot that lets you upload and recall cherished moments with your friends. Whether they are pictures, videos, gifs, or whatever file types you choose! This bot is not hosted publicly, as I do not want to mange media for the public. However, feel free to fork the repo and deploy the bot yourself, for you and your friends.

![Screencastfrom08-06-2024120528PM-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/37d945b3-7a84-427d-ab75-534bb9a0147e)

## Motivation

My friends and I frequently share photos and videos on Discord, ranging from personal moments and memes to inside jokes and gaming clips. We enjoy reliving these funny moments, but it can be challenging to find them. We often forget where we saved them or accidentally delete them. That's why I created this bot—it makes it easy to pull up random memories we've created over the years and share them with everyone in our Discord, making for a fun and nostalgic experience.

## Setup

Setting up the bot is quite simple. You'll need to provide a configuration file and a few secrets. The bot uses a JSON configuration file named `memories.json`, located in the root directory. This file allows you to manage the types of files permitted for upload, their maximum size, and custom options and permissions. Currently, S3 is the only supported storage option for the Discord Memories bot.

### Config 

```json
{
    "storage": {
        "endpoint": "nyc3.digitaloceanspaces.com",
        "region": "nyc3",
        "bucket": "discord-memories",
        "maxFileSize": 25000000,
        "extensions": [
            ".jpg",
            ".jpeg",
            ".png",
            ".gif",
            ".mp4",
            ".webp",
            ".mov"
        ]
    },
    "options": {
        "loki": {
            "path": "loki/",
            "enabled": true,
            "description": "Files related to Loki"
        },
        "lucy": {
            "path": "lucy/",
            "enabled": true,
            "description": "Files related to Lucy"
        }
    },
    "permissions": {
        "servers": {
            "172589280089210880": {
                "enabled": true
            }
        },
        "channels": {
            "1255566989963952148": {
                "enabled": true
            },
        }
    }
}
```

#### Options
The `options` field in `memories.json` allows you to define "buckets" where files will be stored.

```json
"options": {
    "loki": {
        "path": "loki/",
        "enabled": true,
        "description": "Files related to Loki"
    },
}
```

These options will also show up when a user invokes the `!memories help` command.

```
The Discord Memories bot allows you to upload and recall memories made with your friends.
Commands, permissions, file types, and file sizes are all determined by the Memories
configuration file.

Usage:
    !memories [command] [option]

Commands:
    help        Prints information about the Discord Memories bot
    count       Counts the number of files under an option
    read        Retrieves a random file under an option
    upload      Uploads one or more files under an option
    channels    Lists channels that have permissions to use this bot
    servers     Lists servers that have permissions to use this bot

Options:
    loki                    Files related to loki
```

#### Permissions
The `permissions` field in `memories.json` lets you customize access settings. You can whitelist specific servers and channels where the bot can be used. Channel permissions take precedence over server permissions, which is helpful if you want to disable the bot for an entire server except for one channel.

```json
"permissions": {
    "servers": {
        "172589280089210880": {
            "enabled": true
        }
    },
    "channels": {
        "1255566989963952148": {
            "enabled": true
        },
    }
}
```

### Secrets
The final step is to configure the secrets. When the bot starts up, it will load the configuration file and read the required environment variables. You can source these envrionment variables from a `.env ` file or set them manually using another method. Below are the required secrets:

```env
DISCORD_TOKEN=xxx
S3_ACCESS_KEY=xxx
S3_SECRET_KEY=xxx
```
