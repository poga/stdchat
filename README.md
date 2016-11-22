# stdchat

facebook messenger to unix pipe

### Install

```
go get -u github.com/poga/stdchat/...
```

### Usage

First, prepare `chat.yaml`. For example:

```yaml
app_secret: "YOUR_APP_SECRET"
token: "YOUR_BOT_TOKEN"
verify: "VERIFY_TOKEN_FOR_REGISTER_WEBHOOK"
tls_cert: "cert.pem" # webhook requires https. use letsEncrypt to generate your own cert
tls_key: "key.pem" # webhook requires https. use letsEncrypt to generate your own cert
```

Then write a simple bot `bot.rb`:

```ruby
require('json')

loop do
        input = gets

        msg = JSON.parse(input) // expect a input like {"event":{"id":id,"time":time},"opts":{"sender":{"id":"id"},"recipient":{"id":"id"},"timestamp":ts},"message":{"mid":"mid","text":"hi","seq":53},"profile":{...}}
        sender = msg["opts"]["sender"]["id"]
        out = {recipient: sender, message: "hi"}.to_json
        $stdout.puts out
        $stdout.flush
end
```

Finally:

```
chatin | ruby bot.rb | chatout
```

**note**: remember to turn off buffering of `your_program`'s stdout(or remember to flush it).

### Privacy Policy

Facebook messenger requires a privacy police for your app. `chatin` provides a basic template for you to customize: `/privacy`.

### License

The MIT License
