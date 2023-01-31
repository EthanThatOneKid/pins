package filter

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/ethanthatonekid/pins/pins"
	"github.com/google/cel-go/cel"
)

// MessageData is the data for a message in a CEL expression.
type MessageData struct {
	ChannelParentID discord.ChannelID
	ChannelID       discord.ChannelID
	ChannelName     string
	Timestamp       time.Time
	AuthorID        discord.UserID
	AuthorName      string
	Text            string
}

// New creates a new Discord client.
func New(expr pins.FilterExpr) (*Filter, error) {
	if expr == "" {
		return &Filter{}, nil
	}

	env, err := cel.NewEnv(
		cel.Variable("channel_parent_id", cel.IntType),
		cel.Variable("channel_id", cel.IntType),
		cel.Variable("channel_name", cel.StringType),
		cel.Variable("timestamp", cel.TimestampType),
		cel.Variable("author_id", cel.IntType),
		cel.Variable("author_name", cel.StringType),
		cel.Variable("text", cel.StringType),
	)
	if err != nil {
		return nil, err
	}

	ast, issues := env.Compile(string(expr))
	if issues.Err() != nil {
		return nil, issues.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	return &Filter{prg}, nil
}

type Filter struct {
	prg cel.Program
}

func (m Filter) Filter(data MessageData) (bool, error) {
	if m.prg == nil {
		return true, nil
	}

	out, _, err := m.prg.Eval(map[string]interface{}{
		"channel_parent_id": data.ChannelParentID,
		"channel_id":        data.ChannelID,
		"channel_name":      data.ChannelName,
		"timestamp":         data.Timestamp,
		"author_id":         data.AuthorID,
		"author_name":       data.AuthorName,
		"text":              data.Text,
	})
	if err != nil {
		return false, err
	}

	return out.Value().(bool), nil
}
