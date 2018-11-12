### Simple notifier
Universal adapter for send message to instant messengers. Now supported:

- Telegram

### Install to heroku
For repos, clone and exec:
```
heroku create you_app_name --buildpack heroku/go
```

### Development
Install fresh:
```
go get github.com/pilu/fresh
```

Run app:
```
fresh
```

### Usage local

```
go build && ./simple-notifer
curl -i -X POST -F "service=telegram" -F "message=Hello" -F "chat_id=-296931564" -F "token=your_token" localhost:9191/send-message
```