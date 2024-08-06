# Discord Memories

Discord Memories is a bot that lets you upload and recall cherished moments with your friends. Whether they are pictures, videos, gifs, or whatever file types you choose! This bot is not hosted publicly, as I do not want to mange media for the public. However, feel free to fork the repo and deploy the bot yourself, for you and your friends.

## Motivation

My friends and I frequently share photos and videos on Discord, ranging from personal moments and memes to inside jokes and gaming clips. We enjoy reliving these funny moments, but it can be challenging to find them. We often forget where we saved them or accidentally delete them. That's why I created this bot—it makes it easy to pull up random memories we've created over the years and share them with everyone in our Discord, making for a fun and nostalgic experience.

## Demo

![Screencastfrom08-06-2024120528PM-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/37d945b3-7a84-427d-ab75-534bb9a0147e)

## Config Driven
The bot is powered by a JSON configuration file called `memories.json` located in the root directory. The configruation file lets you manage the types of files allowed for upload, their maximum file size, as well as custom commands and permissions. Currently, S3 is the only form of storage supported by the Discord Memories bot:

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

### Secrets
On startup, the bot will load the configuration file and read environment variables. You can either source these from a .env file or set the environment variables manually using another method. Below is the required secrets:

```env
DISCORD_TOKEN=xxx
S3_ACCESS_KEY=xxx
S3_SECRET_KEY=xxx
```
