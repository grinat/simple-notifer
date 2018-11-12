### Simple notifier
Universal adapter for send message to instant messengers.

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
Build:

```
go build && ./simple-notifer
```

Send response for see help page with services and methods:

```
http://localhost:9191
```