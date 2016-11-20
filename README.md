# botpipe

facebook messenger to unix pipe

### Install

```
go get -u github.com/poga/botlisten/...
```

### Usage

First, prepare `bot.toml`. For example:

```toml
app_secret = <YOUR_APP_SECRET>
token = <YOUR_BOT_TOKEN>
verify = <VERIFY_TOKEN_FOR_REGISTER_WEBHOOK>
tls_cert = "cert.pem" # webhook requires https. use letsEncrypt to generate your own cert
tls_key = "key.pem" # webhook requires https. use letsEncrypt to generate your own cert
```

Then:

```
botlisten | your_program | botsay
```

**note**: remember to turn off buffering of your_program's stdout(or remember to flush it).

### License

The MIT License
