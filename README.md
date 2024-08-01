# Discord Memories

Discord Memories is a bot that lets you upload and recall cherished moments with your friends. Whether they are pictures, videos, gifs, or whatever file types you choose to upload! This bot is not hosted publicly, as I do not want to mange media for the public. However, feel free to fork the repo and deploy your own Discord Memories bot! 

### Retrieving Content
![Screencastfrom07-31-202410_31_54PMonline-video-cutter com-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/cee508c1-e7c3-4c31-a2cc-ece5bbb3ae31)

### Uploading Content
![Screencastfrom07-31-202410_38_30PMonline-video-cutter com-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/98b52bad-5ae3-4819-8b18-bb99810a1639)

### Config Driven
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
            ".mp4"
        ]
    },
    "commands": {
        "kevin": {
            "path": "kevin/",
            "enabled": true,
            "description": "Content related to Kevin"
        },
    },
    "permissions": {
        "servers": {
            "172589280089210880": {
                "name": "Slayers of the Bright Realm",
                "description": "Optional description",
                "enabled": true
            }
        },
        "channels": {
            "1255566989963952148": {
                "name:": "bot-testing",
                "description": "Optional description",
                "enabled": true
            },
        }
    }
}
```

### Secrets
On startup, the bot will read the configuration file along with environment variables. These environment variables are sourced from the `.env` file. The following secrets are required:

```env
DISCORD_TOKEN=xxx
S3_ACCESS_KEY=xxx
S3_SECRET_KEY=xxx
```
