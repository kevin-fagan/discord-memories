Discord Memories is a bot that lets you upload and recall cherished moments with your friends, whether they are pictures, videos, or gifs.


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