package inputs

import "github.com/EMCECS/influx/telegraf/plugins"

type baseInput int

func (b baseInput) Type() plugins.Type {
	return plugins.Input
}
