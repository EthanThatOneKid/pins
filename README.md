# pins

List pins from Discord channel, thread, category, or server.

## Example usage

- Go to <https://discord.com/developers/applications>
- Create a new application
- Go to the "Bot" tab
- Click "Add Bot"
- Enable "Message Content Intent"
- Copy the token
- Paste the token into file `.env` like this:

```bash
DISCORD_TOKEN="YOUR_TOKEN_HERE"
```

- Go to the "OAuth2" tab
- Select "bot" under "Scopes"
- Select "Read Message History" under "Bot Permissions"
- Copy the URL and open it in your browser
- Add the bot to your server
- Go to the "General Information" tab
- Copy the "Application ID"
- Run the program with your guild ID like this:

```bash
go run . --guild "710225099923521558" --expr "author_name.size() == 0 || author_name.startsWith(\"EthanThatOneKid\")" --verbose
```

```bash
go run . --guild "710225099923521558" --expr "channel_parent_id == 1055675648296890388" --verbose
```

### Reference

```go
// MessageData is the data for a message in a CEL expression.
type MessageData struct {
	ChannelParentID string
	ChannelID       string
	ChannelName     string
	Timestamp       time.Time
	AuthorID        string
	AuthorName      string
	Text            string
}
```