# Simple notifier

# Development

```
fresh
```

# Usage local

```
go get && go build && ./simple-notifer
curl -i -X POST -F "service=telegram" -F "message=Hello" -F "chat_id=-296931564" -F "token=your_token" localhost:9191/send-message
```