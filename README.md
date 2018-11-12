### Simple notifier
Universal adapter for send message to instant messengers.

### Install to heroku
For repos, clone and exec:
```
# create app
heroku create you_app_name --buildpack heroku/go

# enable metadata support https://devcenter.heroku.com/articles/dyno-metadata
heroku labs:enable runtime-dyno-metadata -a you_app_name
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