Discord Memories is a bot that lets you upload and recall cherished moments with your friends, whether they are pictures, videos, gifs, or whatever files you choose to upload! This bot is not hosted publicly, as I do not want to mange media for the public. However, feel free to fork the repo and deploy your own Discord Memories bot! 

### Retrieving Random Content
![Screencastfrom07-31-202410_31_54PMonline-video-cutter com-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/cee508c1-e7c3-4c31-a2cc-ece5bbb3ae31)

### Uploading Content

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
    "arguments": {
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
