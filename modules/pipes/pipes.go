package pipes

import (
	"github.com/bheru27/glitzz/config"
	"github.com/bheru27/glitzz/core"
	"strings"
)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &pipes{
		Base: core.NewBase("pipes", sender, conf),
	}
	rv.AddCommand("upper", rv.upper)
	rv.AddCommand("lower", rv.lower)
	rv.AddCommand("echo", rv.echo)
	return rv, nil
}

type pipes struct {
	core.Base
}

func (p *pipes) upper(arguments core.CommandArguments) ([]string, error) {
	return []string{strings.ToUpper(strings.Join(arguments.Arguments, " "))}, nil
}

func (p *pipes) lower(arguments core.CommandArguments) ([]string, error) {
	return []string{strings.ToLower(strings.Join(arguments.Arguments, " "))}, nil
}

func (p *pipes) echo(arguments core.CommandArguments) ([]string, error) {
	return []string{strings.Join(arguments.Arguments, " ")}, nil
}
